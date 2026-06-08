package verification

import "time"

type Task struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	TaskNo    string    `gorm:"size:64;uniqueIndex;not null" json:"task_no"`
	AssetID   uint64    `gorm:"index;not null" json:"asset_id"`
	Status    string    `gorm:"size:32;index" json:"status"`
	Score     int       `json:"score"`
	Result    string    `gorm:"size:32;index" json:"result"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Task) TableName() string { return "asset_verification_task" }

type Conflict struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	TaskID    uint64    `gorm:"index;not null" json:"task_id"`
	AssetID   uint64    `gorm:"index;not null" json:"asset_id"`
	Field     string    `gorm:"size:64" json:"field"`
	Expected  string    `gorm:"type:text" json:"expected"`
	Actual    string    `gorm:"type:text" json:"actual"`
	Severity  string    `gorm:"size:32" json:"severity"`
	CreatedAt time.Time `json:"created_at"`
}

func (Conflict) TableName() string { return "asset_verification_conflict" }

const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	ResultPassed    = "passed"
	ResultWarning   = "warning"
	ResultFailed    = "failed"
	ResultUnknown   = "unknown"
)
