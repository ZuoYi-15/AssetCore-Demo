package identity

import "time"

type GenerateRequest struct {
	TenantID     string `json:"tenant_id"`
	SerialNumber string `json:"serial_number"`
	Vendor       string `json:"vendor"`
	Model        string `json:"model"`
	MACAddress   string `json:"mac_address"`
	IPAddress    string `json:"ip_address"`
	Source       string `json:"source"`
}

type BindRequest struct {
	AssetID uint64 `json:"asset_id" binding:"required"`
}

type Query struct {
	Keyword string
	Status  string
}

type IdentityRecord struct {
	ID              uint64    `json:"id"`
	IdentityID      string    `json:"identity_id"`
	FingerprintHash string    `json:"fingerprint_hash"`
	IdentityLevel   string    `json:"identity_level"`
	AssetID         uint64    `json:"asset_id"`
	AssetName       string    `json:"asset_name"`
	AssetType       string    `json:"asset_type"`
	SerialNumber    string    `json:"serial_number"`
	OwnerDepartment string    `json:"owner_department"`
	Location        string    `json:"location"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
