package identity

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
