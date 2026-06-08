package identity

import "time"

type Identity struct {
	ID              uint64    `gorm:"primaryKey" json:"id"`
	IdentityID      string    `gorm:"size:128;uniqueIndex;not null" json:"identity_id"`
	FingerprintHash string    `gorm:"size:128;index;not null" json:"fingerprint_hash"`
	IdentityLevel   string    `gorm:"size:32" json:"identity_level"`
	AssetID         uint64    `gorm:"index" json:"asset_id"`
	Status          string    `gorm:"size:32" json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (Identity) TableName() string { return "asset_identity" }

type Feature struct {
	ID               uint64    `gorm:"primaryKey" json:"id"`
	IdentityID       string    `gorm:"size:128;index;not null" json:"identity_id"`
	FeatureKey       string    `gorm:"size:64;not null" json:"feature_key"`
	FeatureValueHash string    `gorm:"size:128;not null" json:"feature_value_hash"`
	Confidence       int       `json:"confidence"`
	Source           string    `gorm:"size:64" json:"source"`
	CreatedAt        time.Time `json:"created_at"`
}

func (Feature) TableName() string { return "asset_identity_feature" }

const (
	LevelStrong  = "strong"
	LevelMedium  = "medium"
	LevelWeak    = "weak"
	StatusActive = "active"
)
