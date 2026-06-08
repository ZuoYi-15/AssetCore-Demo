package identity

import "gorm.io/gorm"

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

func (r *Repository) Update(item *Identity) error {
	return r.db.Save(item).Error
}

func (r *Repository) ListFeatures(identityID string) ([]Feature, error) {
	var items []Feature
	err := r.db.Where("identity_id = ?", identityID).Order("id ASC").Find(&items).Error
	return items, err
}
