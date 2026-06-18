package asset

import (
	"context"
	"fmt"
	"time"

	"asset-core/internal/infrastructure/kafka"
	"asset-core/internal/module/identity"
)

type SearchIndexer interface {
	Enabled() bool
	IndexAsset(ctx context.Context, item Asset) error
	DeleteAsset(ctx context.Context, id uint64) error
	SearchAssets(ctx context.Context, q Query, offset, limit int) (*SearchResult, error)
}

type SearchResult struct {
	Items []Asset
	Total int64
}

type Service struct {
	repo            *Repository
	identityService *identity.Service
	producer        kafka.Producer
	search          SearchIndexer
}

func NewService(repo *Repository, identityService *identity.Service, producer kafka.Producer, search SearchIndexer) *Service {
	return &Service{repo: repo, identityService: identityService, producer: producer, search: search}
}

func (s *Service) Create(req CreateRequest) (*Asset, error) {
	id, err := s.identityService.Generate(identity.GenerateRequest{
		SerialNumber: req.SerialNumber,
		Vendor:       req.Vendor,
		Model:        req.Model,
		MACAddress:   req.MACAddress,
		IPAddress:    req.IPAddress,
		Source:       req.Source,
	})
	if err != nil {
		return nil, err
	}

	item := &Asset{
		IdentityID:              id.IdentityID,
		AssetHashID:             req.AssetHashID,
		RFIDUID:                 req.RFIDUID,
		AssetName:               req.AssetName,
		AssetType:               req.AssetType,
		Vendor:                  req.Vendor,
		Model:                   req.Model,
		SerialNumber:            req.SerialNumber,
		MACAddress:              req.MACAddress,
		IPAddress:               req.IPAddress,
		Hostname:                req.Hostname,
		OwnerDepartment:         req.OwnerDepartment,
		OwnerUser:               req.OwnerUser,
		Location:                req.Location,
		Building:                req.Building,
		Floor:                   req.Floor,
		Room:                    req.Room,
		Source:                  req.Source,
		TrustLevel:              id.IdentityLevel,
		Status:                  req.Status,
		InitialValue:            req.InitialValue,
		DepMonths:               req.DepMonths,
		AccumulatedDepreciation: req.AccumulatedDepreciation,
		ImpairmentProvision:     req.ImpairmentProvision,
		InServiceDate:           req.InServiceDate,
	}
	prepareAsset(item)
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	_, _ = s.identityService.Bind(id.IdentityID, item.ID)
	s.indexAsset(*item)
	_ = s.producer.Publish("asset.created", kafka.NewEvent("asset.created", item))
	return item, nil
}

func (s *Service) List(q Query, offset, limit int) ([]Asset, int64, error) {
	if s.search != nil && s.search.Enabled() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		result, err := s.search.SearchAssets(ctx, q, offset, limit)
		if err == nil && result != nil {
			return result.Items, result.Total, nil
		}
	}
	return s.repo.List(q, offset, limit)
}

func (s *Service) Get(id uint64) (*Asset, error) {
	return s.repo.FindByID(id)
}

func (s *Service) GenerateIdentity(id uint64) (*Asset, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item.IdentityID != "" {
		return item, nil
	}
	source := item.Source
	if source == "" {
		source = "asset-ledger"
	}
	generated, err := s.identityService.Generate(identity.GenerateRequest{
		TenantID:     "default",
		SerialNumber: item.SerialNumber,
		Vendor:       item.Vendor,
		Model:        item.Model,
		MACAddress:   item.MACAddress,
		IPAddress:    item.IPAddress,
		Source:       source,
	})
	if err != nil {
		return nil, err
	}
	if generated.AssetID != 0 && generated.AssetID != item.ID {
		return nil, identity.ErrIdentityAlreadyBound
	}
	boundIdentity, err := s.identityService.Bind(generated.IdentityID, item.ID)
	if err != nil {
		return nil, err
	}
	before := *item
	item.IdentityID = boundIdentity.IdentityID
	item.TrustLevel = boundIdentity.IdentityLevel
	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	s.writeChangeLogs(before, *item)
	s.indexAsset(*item)
	_ = s.producer.Publish("asset.identity.bound", kafka.NewEvent("asset.identity.bound", item))
	return item, nil
}

func (s *Service) RefreshIdentity(id uint64) (*Asset, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item.IdentityID == "" {
		return item, nil
	}
	source := item.Source
	if source == "" {
		source = "asset-ledger"
	}
	refreshed, err := s.identityService.RefreshBound(item.IdentityID, item.ID, identity.GenerateRequest{
		TenantID:     "default",
		SerialNumber: item.SerialNumber,
		Vendor:       item.Vendor,
		Model:        item.Model,
		MACAddress:   item.MACAddress,
		IPAddress:    item.IPAddress,
		Source:       source,
	})
	if err != nil {
		return nil, err
	}
	before := *item
	item.IdentityID = refreshed.IdentityID
	item.TrustLevel = refreshed.IdentityLevel
	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	s.writeChangeLogs(before, *item)
	s.indexAsset(*item)
	_ = s.producer.Publish("asset.identity.refreshed", kafka.NewEvent("asset.identity.refreshed", item))
	return item, nil
}

func (s *Service) Update(id uint64, req UpdateRequest) (*Asset, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	before := *item
	applyUpdate(item, req)
	prepareAsset(item)
	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	s.writeChangeLogs(before, *item)
	s.indexAsset(*item)
	_ = s.producer.Publish("asset.updated", kafka.NewEvent("asset.updated", item))
	return item, nil
}

func (s *Service) Delete(id uint64) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.deleteAssetIndex(id)
	_ = s.producer.Publish("asset.deleted", kafka.NewEvent("asset.deleted", map[string]uint64{"id": id}))
	return nil
}

func (s *Service) UpdateStatus(id uint64, status string) (*Asset, error) {
	return s.Update(id, UpdateRequest{Status: status})
}

func (s *Service) AddInsurance(id uint64, req InsuranceRequest) (*Insurance, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	insurance := &Insurance{
		AssetID:       item.ID,
		PolicyNo:      req.PolicyNo,
		AnnualPremium: req.AnnualPremium,
		InsuredAmount: req.InsuredAmount,
		PeriodStart:   req.PeriodStart,
		PeriodEnd:     req.PeriodEnd,
		Operator:      req.Operator,
	}
	if err := s.repo.CreateInsurance(insurance); err != nil {
		return nil, err
	}
	_ = s.repo.CreateChangeLog(&ChangeLog{
		AssetID:  item.ID,
		Field:    "insurance_policy_no",
		OldValue: "",
		NewValue: req.PolicyNo,
		Operator: operatorOrSystem(req.Operator),
	})
	_ = s.producer.Publish("asset.insurance.updated", kafka.NewEvent("asset.insurance.updated", insurance))
	return insurance, nil
}

func (s *Service) ListInsurance(assetID uint64) ([]Insurance, error) {
	return s.repo.ListInsurance(assetID)
}

func (s *Service) RecordImpairment(id uint64, req ImpairmentRequest) (*ImpairmentRecord, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	prepareAsset(item)
	current := item.CurrentNetValue
	amount := current - req.RecoverableAmount
	if amount < 0 {
		amount = 0
	}
	before := *item
	item.ImpairmentProvision += amount
	prepareAsset(item)
	record := &ImpairmentRecord{
		AssetID:           item.ID,
		Reason:            req.Reason,
		EvidenceFileHash:  req.EvidenceFileHash,
		RecoverableAmount: req.RecoverableAmount,
		ImpairmentAmount:  amount,
		Reviewer:          req.Reviewer,
	}
	if err := s.repo.CreateImpairmentRecord(record); err != nil {
		return nil, err
	}
	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	s.writeChangeLogs(before, *item)
	_ = s.repo.CreateChangeLog(&ChangeLog{
		AssetID:  item.ID,
		Field:    "asset_impairment",
		OldValue: fmt.Sprintf("%.2f", before.ImpairmentProvision),
		NewValue: fmt.Sprintf("%.2f", item.ImpairmentProvision),
		Operator: operatorOrSystem(req.Reviewer),
	})
	s.indexAsset(*item)
	_ = s.producer.Publish("asset.impairment.recorded", kafka.NewEvent("asset.impairment.recorded", record))
	return record, nil
}

func (s *Service) ListImpairments(assetID uint64) ([]ImpairmentRecord, error) {
	return s.repo.ListImpairmentRecords(assetID)
}

func (s *Service) ChangeLogs(assetID uint64) ([]ChangeLog, error) {
	return s.repo.ListChangeLogs(assetID)
}

func (s *Service) RecordChangeLog(assetID uint64, field, oldValue, newValue, operator string) error {
	return s.repo.CreateChangeLog(&ChangeLog{
		AssetID:  assetID,
		Field:    field,
		OldValue: oldValue,
		NewValue: newValue,
		Operator: operator,
	})
}

func (s *Service) HasChangeLog(assetID uint64, field, newValue string) (bool, error) {
	return s.repo.HasChangeLog(assetID, field, newValue)
}

func (s *Service) indexAsset(item Asset) {
	if s.search == nil || !s.search.Enabled() {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = s.search.IndexAsset(ctx, item)
}

func (s *Service) deleteAssetIndex(id uint64) {
	if s.search == nil || !s.search.Enabled() {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = s.search.DeleteAsset(ctx, id)
}

func applyUpdate(item *Asset, req UpdateRequest) {
	if req.AssetHashID != "" {
		item.AssetHashID = req.AssetHashID
	}
	if req.RFIDUID != "" {
		item.RFIDUID = req.RFIDUID
	}
	if req.AssetName != "" {
		item.AssetName = req.AssetName
	}
	if req.AssetType != "" {
		item.AssetType = req.AssetType
	}
	if req.Vendor != "" {
		item.Vendor = req.Vendor
	}
	if req.Model != "" {
		item.Model = req.Model
	}
	if req.SerialNumber != "" {
		item.SerialNumber = req.SerialNumber
	}
	if req.MACAddress != "" {
		item.MACAddress = req.MACAddress
	}
	if req.IPAddress != "" {
		item.IPAddress = req.IPAddress
	}
	if req.Hostname != "" {
		item.Hostname = req.Hostname
	}
	if req.OwnerDepartment != "" {
		item.OwnerDepartment = req.OwnerDepartment
	}
	if req.OwnerUser != "" {
		item.OwnerUser = req.OwnerUser
	}
	if req.Location != "" {
		item.Location = req.Location
	}
	if req.Building != "" {
		item.Building = req.Building
	}
	if req.Floor != "" {
		item.Floor = req.Floor
	}
	if req.Room != "" {
		item.Room = req.Room
	}
	if req.Source != "" {
		item.Source = req.Source
	}
	if req.TrustLevel != "" {
		item.TrustLevel = req.TrustLevel
	}
	if req.Status != "" {
		item.Status = req.Status
	}
	if req.InitialValue > 0 {
		item.InitialValue = req.InitialValue
	}
	if req.DepMonths > 0 {
		item.DepMonths = req.DepMonths
	}
	if req.AccumulatedDepreciation > 0 {
		item.AccumulatedDepreciation = req.AccumulatedDepreciation
	}
	if req.ImpairmentProvision > 0 {
		item.ImpairmentProvision = req.ImpairmentProvision
	}
	if req.InServiceDate != nil {
		item.InServiceDate = req.InServiceDate
	}
	if req.DeactivationDate != nil {
		item.DeactivationDate = req.DeactivationDate
	}
}

func (s *Service) writeChangeLogs(before, after Asset) {
	changes := map[string][2]string{
		"asset_hash_id":            {before.AssetHashID, after.AssetHashID},
		"rfid_uid":                 {before.RFIDUID, after.RFIDUID},
		"asset_name":               {before.AssetName, after.AssetName},
		"asset_type":               {before.AssetType, after.AssetType},
		"vendor":                   {before.Vendor, after.Vendor},
		"model":                    {before.Model, after.Model},
		"serial_number":            {before.SerialNumber, after.SerialNumber},
		"mac_address":              {before.MACAddress, after.MACAddress},
		"ip_address":               {before.IPAddress, after.IPAddress},
		"hostname":                 {before.Hostname, after.Hostname},
		"owner_department":         {before.OwnerDepartment, after.OwnerDepartment},
		"owner_user":               {before.OwnerUser, after.OwnerUser},
		"location":                 {before.Location, after.Location},
		"building":                 {before.Building, after.Building},
		"floor":                    {before.Floor, after.Floor},
		"room":                     {before.Room, after.Room},
		"source":                   {before.Source, after.Source},
		"trust_level":              {before.TrustLevel, after.TrustLevel},
		"status":                   {before.Status, after.Status},
		"initial_value":            {formatAmount(before.InitialValue), formatAmount(after.InitialValue)},
		"depreciation_months":      {fmt.Sprint(before.DepMonths), fmt.Sprint(after.DepMonths)},
		"accumulated_depreciation": {formatAmount(before.AccumulatedDepreciation), formatAmount(after.AccumulatedDepreciation)},
		"impairment_provision":     {formatAmount(before.ImpairmentProvision), formatAmount(after.ImpairmentProvision)},
		"current_net_value":        {formatAmount(before.CurrentNetValue), formatAmount(after.CurrentNetValue)},
	}
	for field, values := range changes {
		if values[0] == values[1] {
			continue
		}
		_ = s.repo.CreateChangeLog(&ChangeLog{
			AssetID:  after.ID,
			Field:    field,
			OldValue: values[0],
			NewValue: values[1],
			Operator: "system",
		})
	}
}

func prepareAsset(item *Asset) {
	if item.Status == "" {
		item.Status = StatusRegistered
	}
	if item.DepMethod == "" {
		item.DepMethod = DepMethodStraightLine
	}
	item.SalvageRate = 0
	now := time.Now()
	if item.Status == StatusInService && item.InServiceDate == nil {
		item.InServiceDate = &now
	}
	if isReducedStatus(item.Status) && item.DeactivationDate == nil {
		item.DeactivationDate = &now
	}
	if item.DepMonths > 0 && item.InServiceDate != nil {
		item.DepreciatedMonths = elapsedDepreciationMonths(*item.InServiceDate, item.DeactivationDate, now)
		if item.DepreciatedMonths > item.DepMonths {
			item.DepreciatedMonths = item.DepMonths
		}
		expected := expectedAccumulatedDepreciation(item)
		if expected > item.AccumulatedDepreciation {
			item.AccumulatedDepreciation = expected
		}
	}
	net := item.InitialValue - item.AccumulatedDepreciation - item.ImpairmentProvision
	if net < 0 {
		net = 0
	}
	item.CurrentNetValue = net
	item.DepreciationStopped = item.DeactivationDate != nil && startOfMonth(now).After(startOfMonth(*item.DeactivationDate))
}

func expectedAccumulatedDepreciation(item *Asset) float64 {
	if item.DepMonths <= 0 || item.DepreciatedMonths <= 0 || item.InitialValue <= 0 {
		return item.AccumulatedDepreciation
	}
	if item.ImpairmentProvision <= 0 {
		return item.InitialValue / float64(item.DepMonths) * float64(item.DepreciatedMonths)
	}
	remainingMonths := item.DepMonths - item.DepreciatedMonths
	if remainingMonths <= 0 {
		return item.InitialValue - item.ImpairmentProvision
	}
	preImpairmentValue := item.InitialValue - item.ImpairmentProvision - item.AccumulatedDepreciation
	if preImpairmentValue <= 0 {
		return item.AccumulatedDepreciation
	}
	return item.AccumulatedDepreciation + preImpairmentValue/float64(remainingMonths)
}

func elapsedDepreciationMonths(inService time.Time, deactivated *time.Time, now time.Time) int {
	start := startOfMonth(inService).AddDate(0, 1, 0)
	end := startOfMonth(now)
	if deactivated != nil {
		deactivationMonth := startOfMonth(*deactivated)
		if deactivationMonth.Before(end) || deactivationMonth.Equal(end) {
			end = deactivationMonth
		}
	}
	if end.Before(start) {
		return 0
	}
	return (end.Year()-start.Year())*12 + int(end.Month()-start.Month()) + 1
}

func startOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func isReducedStatus(status string) bool {
	switch status {
	case StatusRetired, StatusPendingDisposal, StatusDisposed:
		return true
	default:
		return false
	}
}

func formatAmount(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func operatorOrSystem(operator string) string {
	if operator == "" {
		return "system"
	}
	return operator
}
