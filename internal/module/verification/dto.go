package verification

import "time"

type CreateRequest struct {
	AssetID uint64 `json:"asset_id" binding:"required"`
}

type Query struct {
	Keyword string
	Result  string
	Status  string
}

type ResultResponse struct {
	Task      Task       `json:"task"`
	Conflicts []Conflict `json:"conflicts"`
}

type TaskRecord struct {
	ID              uint64    `json:"id"`
	TaskNo          string    `json:"task_no"`
	AssetID         uint64    `json:"asset_id"`
	AssetName       string    `json:"asset_name"`
	AssetType       string    `json:"asset_type"`
	IdentityID      string    `json:"identity_id"`
	SerialNumber    string    `json:"serial_number"`
	OwnerDepartment string    `json:"owner_department"`
	Status          string    `json:"status"`
	Score           int       `json:"score"`
	Result          string    `json:"result"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
