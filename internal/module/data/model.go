package data

import "time"

type ImportTask struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	TaskNo       string    `gorm:"size:64;uniqueIndex;not null" json:"task_no"`
	FileName     string    `gorm:"size:255" json:"file_name"`
	FileURL      string    `gorm:"size:512" json:"file_url"`
	Status       string    `gorm:"size:32;index" json:"status"`
	TotalCount   int       `json:"total_count"`
	SuccessCount int       `json:"success_count"`
	FailedCount  int       `json:"failed_count"`
	OperatorID   string    `gorm:"size:128" json:"operator_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (ImportTask) TableName() string { return "data_import_task" }

type ImportError struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	TaskID       uint64    `gorm:"index;not null" json:"task_id"`
	RowNumber    int       `json:"row_number"`
	ErrorField   string    `gorm:"size:64" json:"error_field"`
	ErrorMessage string    `gorm:"size:512" json:"error_message"`
	RawData      string    `gorm:"type:text" json:"raw_data"`
	CreatedAt    time.Time `json:"created_at"`
}

func (ImportError) TableName() string { return "data_import_error" }

const (
	ImportPending   = "pending"
	ImportCompleted = "completed"
	ImportFailed    = "failed"
)
