package asset

import "time"

type Asset struct {
	ID              uint64     `gorm:"primaryKey" json:"id"`
	IdentityID      string     `gorm:"size:128;uniqueIndex" json:"identity_id"`
	AssetName       string     `gorm:"size:128;not null" json:"asset_name"`
	AssetType       string     `gorm:"size:64;index" json:"asset_type"`
	Vendor          string     `gorm:"size:128" json:"vendor"`
	Model           string     `gorm:"size:128" json:"model"`
	SerialNumber    string     `gorm:"size:128;index" json:"serial_number"`
	MACAddress      string     `gorm:"size:64;index" json:"mac_address"`
	IPAddress       string     `gorm:"size:64;index" json:"ip_address"`
	Hostname        string     `gorm:"size:128" json:"hostname"`
	OwnerDepartment string     `gorm:"size:128" json:"owner_department"`
	OwnerUser       string     `gorm:"size:128" json:"owner_user"`
	Location        string     `gorm:"size:255" json:"location"`
	Source          string     `gorm:"size:64" json:"source"`
	TrustLevel      string     `gorm:"size:32" json:"trust_level"`
	Status          string     `gorm:"size:32;index" json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `gorm:"index" json:"-"`
}

func (Asset) TableName() string { return "asset" }

type ChangeLog struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	AssetID   uint64    `gorm:"index;not null" json:"asset_id"`
	Field     string    `gorm:"size:64" json:"field"`
	OldValue  string    `gorm:"type:text" json:"old_value"`
	NewValue  string    `gorm:"type:text" json:"new_value"`
	Operator  string    `gorm:"size:128" json:"operator"`
	CreatedAt time.Time `json:"created_at"`
}

func (ChangeLog) TableName() string { return "asset_change_log" }

const (
	StatusDiscovered = "discovered"
	StatusRegistered = "registered"
	StatusVerified   = "verified"
	StatusAbnormal   = "abnormal"
	StatusRetired    = "retired"
)
