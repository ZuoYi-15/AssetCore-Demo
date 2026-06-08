package asset

type CreateRequest struct {
	AssetName       string `json:"asset_name" binding:"required"`
	AssetType       string `json:"asset_type"`
	Vendor          string `json:"vendor"`
	Model           string `json:"model"`
	SerialNumber    string `json:"serial_number"`
	MACAddress      string `json:"mac_address"`
	IPAddress       string `json:"ip_address"`
	Hostname        string `json:"hostname"`
	OwnerDepartment string `json:"owner_department"`
	OwnerUser       string `json:"owner_user"`
	Location        string `json:"location"`
	Source          string `json:"source"`
}

type UpdateRequest struct {
	AssetName       string `json:"asset_name"`
	AssetType       string `json:"asset_type"`
	Vendor          string `json:"vendor"`
	Model           string `json:"model"`
	SerialNumber    string `json:"serial_number"`
	MACAddress      string `json:"mac_address"`
	IPAddress       string `json:"ip_address"`
	Hostname        string `json:"hostname"`
	OwnerDepartment string `json:"owner_department"`
	OwnerUser       string `json:"owner_user"`
	Location        string `json:"location"`
	Source          string `json:"source"`
	TrustLevel      string `json:"trust_level"`
	Status          string `json:"status"`
}

type StatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type Query struct {
	Keyword   string
	Status    string
	AssetType string
}
