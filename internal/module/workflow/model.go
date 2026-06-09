package workflow

import "time"

const (
	FlowPurchase = "purchase"
	FlowTransfer = "transfer"
	FlowRetire   = "retire"

	StatusActive   = "active"
	StatusInactive = "inactive"

	InstancePending  = "pending"
	InstanceApproved = "approved"
	InstanceRejected = "rejected"

	TaskPending  = "pending"
	TaskApproved = "approved"
	TaskRejected = "rejected"
)

type Definition struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	FlowType  string    `gorm:"size:32;uniqueIndex;not null" json:"flow_type"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	Status    string    `gorm:"size:32;index;not null" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Nodes     []Node    `gorm:"foreignKey:DefinitionID" json:"nodes,omitempty"`
}

func (Definition) TableName() string { return "workflow_definition" }

type Node struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	DefinitionID uint64    `gorm:"index;not null" json:"definition_id"`
	NodeName     string    `gorm:"size:128;not null" json:"node_name"`
	ApproverRole string    `gorm:"size:64;not null" json:"approver_role"`
	SortOrder    int       `gorm:"index;not null" json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Node) TableName() string { return "workflow_node" }

type Instance struct {
	ID            uint64     `gorm:"primaryKey" json:"id"`
	DefinitionID  uint64     `gorm:"index;not null" json:"definition_id"`
	FlowType      string     `gorm:"size:32;index;not null" json:"flow_type"`
	AssetID       uint64     `gorm:"index" json:"asset_id"`
	Title         string     `gorm:"size:255;not null" json:"title"`
	Status        string     `gorm:"size:32;index;not null" json:"status"`
	CurrentTaskID uint64     `gorm:"index" json:"current_task_id"`
	ApplicantID   uint64     `gorm:"index;not null" json:"applicant_id"`
	ApplicantName string     `gorm:"size:128" json:"applicant_name"`
	Payload       string     `gorm:"type:text" json:"payload"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	Tasks         []Task     `gorm:"foreignKey:InstanceID" json:"tasks,omitempty"`
}

func (Instance) TableName() string { return "workflow_instance" }

type Task struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	InstanceID   uint64     `gorm:"index;not null" json:"instance_id"`
	NodeID       uint64     `gorm:"index;not null" json:"node_id"`
	NodeName     string     `gorm:"size:128;not null" json:"node_name"`
	ApproverRole string     `gorm:"size:64;index;not null" json:"approver_role"`
	SortOrder    int        `gorm:"index;not null" json:"sort_order"`
	Status       string     `gorm:"size:32;index;not null" json:"status"`
	ApproverID   uint64     `gorm:"index" json:"approver_id"`
	ApproverName string     `gorm:"size:128" json:"approver_name"`
	Comment      string     `gorm:"size:512" json:"comment"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	Instance     *Instance  `gorm:"foreignKey:InstanceID" json:"instance,omitempty"`
}

func (Task) TableName() string { return "workflow_task" }
