package identity

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"asset-core/internal/infrastructure/kafka"

	"gorm.io/gorm"
)

var (
	ErrIdentityAlreadyBound = errors.New("identity already bound to another asset")
	ErrAssetAlreadyBound    = errors.New("asset already bound to another identity")
)

type Service struct {
	repo     *Repository
	producer kafka.Producer
}

func NewService(repo *Repository, producer kafka.Producer) *Service {
	return &Service{repo: repo, producer: producer}
}

func (s *Service) Generate(req GenerateRequest) (*Identity, error) {
	fingerprint, level := fingerprint(req)
	if existing, err := s.repo.FindByFingerprint(fingerprint); err == nil {
		return existing, nil
	}

	identityID := "did:asset:" + fingerprint[:32]
	item := &Identity{
		IdentityID:      identityID,
		FingerprintHash: fingerprint,
		IdentityLevel:   level,
		Status:          StatusActive,
	}

	features := buildFeatures(identityID, req)
	if err := s.repo.Create(item, features); err != nil {
		return nil, err
	}
	_ = s.producer.Publish("asset.identity.generated", kafka.NewEvent("identity.generated", item))
	return item, nil
}

func (s *Service) Get(identityID string) (*Identity, error) {
	return s.repo.FindByIdentityID(identityID)
}

func (s *Service) List(q Query, offset, limit int) ([]IdentityRecord, int64, error) {
	return s.repo.List(q, offset, limit)
}

func (s *Service) Bind(identityID string, assetID uint64) (*Identity, error) {
	item, err := s.repo.FindByIdentityID(identityID)
	if err != nil {
		return nil, err
	}
	if item.AssetID != 0 && item.AssetID != assetID {
		return nil, ErrIdentityAlreadyBound
	}
	existing, err := s.repo.FindByAssetID(assetID)
	if err == nil && existing.IdentityID != identityID {
		return nil, ErrAssetAlreadyBound
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	item.AssetID = assetID
	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) RefreshBound(identityID string, assetID uint64, req GenerateRequest) (*Identity, error) {
	item, err := s.repo.FindByIdentityID(identityID)
	if err != nil {
		return nil, err
	}
	if item.AssetID != 0 && item.AssetID != assetID {
		return nil, ErrIdentityAlreadyBound
	}

	fingerprint, level := fingerprint(req)
	nextIdentityID := "did:asset:" + fingerprint[:32]
	if existing, err := s.repo.FindByFingerprint(fingerprint); err == nil && existing.ID != item.ID {
		if existing.AssetID != 0 && existing.AssetID != assetID {
			return nil, ErrIdentityAlreadyBound
		}
		return nil, ErrAssetAlreadyBound
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	oldIdentityID := item.IdentityID
	item.IdentityID = nextIdentityID
	item.FingerprintHash = fingerprint
	item.IdentityLevel = level
	item.AssetID = assetID
	item.Status = StatusActive
	features := buildFeatures(nextIdentityID, req)
	if err := s.repo.UpdateWithFeatures(item, oldIdentityID, features); err != nil {
		return nil, err
	}
	_ = s.producer.Publish("asset.identity.refreshed", kafka.NewEvent("identity.refreshed", item))
	return item, nil
}

func (s *Service) Unbind(identityID string) (*Identity, error) {
	item, err := s.repo.FindByIdentityID(identityID)
	if err != nil {
		return nil, err
	}
	item.AssetID = 0
	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) Features(identityID string) ([]Feature, error) {
	return s.repo.ListFeatures(identityID)
}

func fingerprint(req GenerateRequest) (string, string) {
	parts := []string{req.TenantID, req.SerialNumber, req.Vendor, req.Model, req.MACAddress}
	level := LevelStrong
	if req.SerialNumber == "" || req.Vendor == "" || req.Model == "" || req.MACAddress == "" {
		parts = []string{req.TenantID, req.SerialNumber, req.Vendor, req.MACAddress}
		level = LevelMedium
	}
	if req.SerialNumber == "" && req.MACAddress == "" {
		parts = []string{req.TenantID, req.IPAddress, req.Vendor, req.Model}
		level = LevelWeak
	}
	joined := strings.ToLower(strings.Join(parts, "|"))
	sum := sha256.Sum256([]byte(joined))
	return hex.EncodeToString(sum[:]), level
}

func buildFeatures(identityID string, req GenerateRequest) []Feature {
	raw := map[string]string{
		"serial_number": req.SerialNumber,
		"vendor":        req.Vendor,
		"model":         req.Model,
		"mac_address":   req.MACAddress,
		"ip_address":    req.IPAddress,
	}
	items := make([]Feature, 0, len(raw))
	for k, v := range raw {
		if v == "" {
			continue
		}
		sum := sha256.Sum256([]byte(strings.ToLower(v)))
		items = append(items, Feature{
			IdentityID:       identityID,
			FeatureKey:       k,
			FeatureValueHash: hex.EncodeToString(sum[:]),
			Confidence:       80,
			Source:           req.Source,
		})
	}
	return items
}
