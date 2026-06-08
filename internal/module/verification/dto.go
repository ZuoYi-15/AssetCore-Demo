package verification

type CreateRequest struct {
	AssetID uint64 `json:"asset_id" binding:"required"`
}

type ResultResponse struct {
	Task      Task       `json:"task"`
	Conflicts []Conflict `json:"conflicts"`
}
