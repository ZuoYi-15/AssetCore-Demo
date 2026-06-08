package audit

import "time"

type Log struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	TraceID   string    `gorm:"size:128;index" json:"trace_id"`
	Operator  string    `gorm:"size:128" json:"operator"`
	Action    string    `gorm:"size:128;index" json:"action"`
	Resource  string    `gorm:"size:128;index" json:"resource"`
	Detail    string    `gorm:"type:text" json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

func (Log) TableName() string { return "audit_log" }
