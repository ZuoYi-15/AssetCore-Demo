package data

type CreateImportTaskRequest struct {
	FileName   string `json:"file_name" binding:"required"`
	FileURL    string `json:"file_url"`
	OperatorID string `json:"operator_id"`
}

type ImportAssetsResult struct {
	Task *ImportTask `json:"task"`
}
