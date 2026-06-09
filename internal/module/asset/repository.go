package asset

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(a *Asset) error {
	return r.db.Create(a).Error
}

func (r *Repository) Update(a *Asset) error {
	return r.db.Save(a).Error
}

func (r *Repository) FindByID(id uint64) (*Asset, error) {
	var a Asset
	if err := r.db.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) FindByIdentityID(identityID string) (*Asset, error) {
	var a Asset
	if err := r.db.Where("identity_id = ?", identityID).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) List(q Query, offset, limit int) ([]Asset, int64, error) {
	db := r.db.Model(&Asset{})
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		db = db.Where("asset_name LIKE ? OR serial_number LIKE ? OR mac_address LIKE ? OR ip_address LIKE ? OR hostname LIKE ? OR identity_id LIKE ?", like, like, like, like, like, like)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}
	if q.AssetType != "" {
		db = db.Where("asset_type = ?", q.AssetType)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []Asset
	if err := db.Order("id DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) Delete(id uint64) error {
	return r.db.Delete(&Asset{}, id).Error
}

func (r *Repository) CreateChangeLog(log *ChangeLog) error {
	return r.db.Create(log).Error
}

func (r *Repository) HasChangeLog(assetID uint64, field, newValue string) (bool, error) {
	var count int64
	err := r.db.Model(&ChangeLog{}).
		Where("asset_id = ? AND field = ? AND new_value = ?", assetID, field, newValue).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) ListChangeLogs(assetID uint64) ([]ChangeLog, error) {
	var logs []ChangeLog
	err := r.db.Where("asset_id = ?", assetID).Order("id DESC").Find(&logs).Error
	return logs, err
}
