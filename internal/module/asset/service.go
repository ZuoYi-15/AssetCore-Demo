package asset

import (
	"context"
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
		IdentityID:      id.IdentityID,
		AssetName:       req.AssetName,
		AssetType:       req.AssetType,
		Vendor:          req.Vendor,
		Model:           req.Model,
		SerialNumber:    req.SerialNumber,
		MACAddress:      req.MACAddress,
		IPAddress:       req.IPAddress,
		Hostname:        req.Hostname,
		OwnerDepartment: req.OwnerDepartment,
		OwnerUser:       req.OwnerUser,
		Location:        req.Location,
		Source:          req.Source,
		TrustLevel:      id.IdentityLevel,
		Status:          StatusRegistered,
	}
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

func (s *Service) Update(id uint64, req UpdateRequest) (*Asset, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	before := *item
	applyUpdate(item, req)
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
	if req.Source != "" {
		item.Source = req.Source
	}
	if req.TrustLevel != "" {
		item.TrustLevel = req.TrustLevel
	}
	if req.Status != "" {
		item.Status = req.Status
	}
}

func (s *Service) writeChangeLogs(before, after Asset) {
	changes := map[string][2]string{
		"asset_name":       {before.AssetName, after.AssetName},
		"asset_type":       {before.AssetType, after.AssetType},
		"vendor":           {before.Vendor, after.Vendor},
		"model":            {before.Model, after.Model},
		"serial_number":    {before.SerialNumber, after.SerialNumber},
		"mac_address":      {before.MACAddress, after.MACAddress},
		"ip_address":       {before.IPAddress, after.IPAddress},
		"hostname":         {before.Hostname, after.Hostname},
		"owner_department": {before.OwnerDepartment, after.OwnerDepartment},
		"owner_user":       {before.OwnerUser, after.OwnerUser},
		"location":         {before.Location, after.Location},
		"source":           {before.Source, after.Source},
		"trust_level":      {before.TrustLevel, after.TrustLevel},
		"status":           {before.Status, after.Status},
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
