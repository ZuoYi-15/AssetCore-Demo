package asset

import "time"

type Asset struct {
	ID                      uint64     `gorm:"primaryKey" json:"id"`
	IdentityID              string     `gorm:"size:128;uniqueIndex" json:"identity_id"`
	AssetHashID             string     `gorm:"size:128;index" json:"asset_hash_id"`
	RFIDUID                 string     `gorm:"column:rfid_uid;size:128;index" json:"rfid_uid"`
	AssetName               string     `gorm:"size:128;not null" json:"asset_name"`
	AssetType               string     `gorm:"size:64;index" json:"asset_type"`
	Vendor                  string     `gorm:"size:128" json:"vendor"`
	Model                   string     `gorm:"size:128" json:"model"`
	SerialNumber            string     `gorm:"size:128;index" json:"serial_number"`
	MACAddress              string     `gorm:"size:64;index" json:"mac_address"`
	IPAddress               string     `gorm:"size:64;index" json:"ip_address"`
	Hostname                string     `gorm:"size:128" json:"hostname"`
	OwnerDepartment         string     `gorm:"size:128" json:"owner_department"`
	OwnerUser               string     `gorm:"size:128" json:"owner_user"`
	Location                string     `gorm:"size:255" json:"location"`
	Building                string     `gorm:"size:128;index" json:"building"`
	Floor                   string     `gorm:"size:64;index" json:"floor"`
	Room                    string     `gorm:"size:128;index" json:"room"`
	Source                  string     `gorm:"size:64" json:"source"`
	TrustLevel              string     `gorm:"size:32" json:"trust_level"`
	Status                  string     `gorm:"size:32;index" json:"status"`
	InitialValue            float64    `json:"initial_value"`
	DepMethod               string     `gorm:"size:32" json:"depreciation_method"`
	DepMonths               int        `json:"depreciation_months"`
	SalvageRate             float64    `json:"salvage_rate"`
	InServiceDate           *time.Time `json:"in_service_date"`
	DeactivationDate        *time.Time `json:"deactivation_date"`
	DepreciatedMonths       int        `json:"depreciated_months"`
	AccumulatedDepreciation float64    `json:"accumulated_depreciation"`
	ImpairmentProvision     float64    `json:"impairment_provision"`
	CurrentNetValue         float64    `json:"current_net_value"`
	DepreciationStopped     bool       `json:"depreciation_stopped"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
	DeletedAt               *time.Time `gorm:"index" json:"-"`
}

func (Asset) TableName() string { return "asset" }

type ChangeLog struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	AssetID   uint64    `gorm:"index;not null" json:"asset_id"`
	Field     string    `gorm:"size:64" json:"field"`
	OldValue  string    `gorm:"type:text" json:"old_value"`
	NewValue  string    `gorm:"type:text" json:"new_value"`
	Operator  string    `gorm:"size:128" json:"operator"`
	CreatedAt time.Time `json:"created_at"`
}

func (ChangeLog) TableName() string { return "asset_change_log" }

type Insurance struct {
	ID            uint64    `gorm:"primaryKey" json:"id"`
	AssetID       uint64    `gorm:"index;not null" json:"asset_id"`
	PolicyNo      string    `gorm:"column:insurance_policy_no;size:128;index" json:"insurance_policy_no"`
	AnnualPremium float64   `json:"annual_premium"`
	InsuredAmount float64   `json:"insured_amount"`
	PeriodStart   time.Time `json:"period_start"`
	PeriodEnd     time.Time `json:"period_end"`
	Operator      string    `gorm:"size:128" json:"operator"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Insurance) TableName() string { return "asset_insurance" }

type ImpairmentRecord struct {
	ID                uint64    `gorm:"primaryKey" json:"id"`
	AssetID           uint64    `gorm:"index;not null" json:"asset_id"`
	Reason            string    `gorm:"type:text" json:"reason"`
	EvidenceFileHash  string    `gorm:"size:128" json:"evidence_file_hash"`
	RecoverableAmount float64   `json:"recoverable_amount"`
	ImpairmentAmount  float64   `json:"impairment_amount"`
	Reviewer          string    `gorm:"size:128" json:"reviewer"`
	CreatedAt         time.Time `json:"created_at"`
}

func (ImpairmentRecord) TableName() string { return "asset_impairment_record" }

const (
	StatusDiscovered      = "discovered"
	StatusRegistered      = "registered"
	StatusVerified        = "verified"
	StatusAbnormal        = "abnormal"
	StatusRetired         = "retired"
	StatusLegacyPending   = "legacy_pending"
	StatusInWarehouse     = "in_warehouse"
	StatusInService       = "in_service"
	StatusUnderRepair     = "under_repair"
	StatusTransferring    = "transferring"
	StatusPendingDisposal = "pending_disposal"
	StatusDisposed        = "disposed"
)

const DepMethodStraightLine = "straight_line"
