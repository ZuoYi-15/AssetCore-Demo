package identity

import (
	"strings"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(identity *Identity, features []Feature) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(identity).Error; err != nil {
			return err
		}
		if len(features) > 0 {
			if err := tx.Create(&features).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) FindByIdentityID(identityID string) (*Identity, error) {
	var item Identity
	if err := r.db.Where("identity_id = ?", identityID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) FindByFingerprint(hash string) (*Identity, error) {
	var item Identity
	if err := r.db.Where("fingerprint_hash = ?", hash).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) FindByAssetID(assetID uint64) (*Identity, error) {
	var item Identity
	if err := r.db.Where("asset_id = ?", assetID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) List(q Query, offset, limit int) ([]IdentityRecord, int64, error) {
	db := r.db.Table("asset_identity").
		Select(`asset_identity.id,
			asset_identity.identity_id,
			asset_identity.fingerprint_hash,
			asset_identity.identity_level,
			asset_identity.asset_id,
			asset.asset_name,
			asset.asset_type,
			asset.serial_number,
			asset.owner_department,
			asset.location,
			asset_identity.status,
			asset_identity.created_at,
			asset_identity.updated_at`).
		Joins("LEFT JOIN asset ON asset.id = asset_identity.asset_id")
	if strings.TrimSpace(q.Keyword) != "" {
		like := "%" + strings.TrimSpace(q.Keyword) + "%"
		db = db.Where("asset_identity.identity_id LIKE ? OR asset_identity.fingerprint_hash LIKE ? OR asset.asset_name LIKE ? OR asset.serial_number LIKE ? OR asset.mac_address LIKE ? OR asset.ip_address LIKE ?", like, like, like, like, like, like)
	}
	if strings.TrimSpace(q.Status) != "" {
		db = db.Where("asset_identity.status = ?", strings.TrimSpace(q.Status))
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []IdentityRecord
	if err := db.Order("asset_identity.id DESC").Offset(offset).Limit(limit).Scan(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) Update(item *Identity) error {
	return r.db.Save(item).Error
}

func (r *Repository) UpdateWithFeatures(item *Identity, oldIdentityID string, features []Feature) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(item).Error; err != nil {
			return err
		}
		if err := tx.Where("identity_id = ?", oldIdentityID).Delete(&Feature{}).Error; err != nil {
			return err
		}
		if len(features) > 0 {
			if err := tx.Create(&features).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) ListFeatures(identityID string) ([]Feature, error) {
	var items []Feature
	err := r.db.Where("identity_id = ?", identityID).Order("id ASC").Find(&items).Error
	return items, err
}
