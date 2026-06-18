package asset

import "time"

type CreateRequest struct {
	AssetName               string     `json:"asset_name" binding:"required"`
	AssetHashID             string     `json:"asset_hash_id"`
	RFIDUID                 string     `json:"rfid_uid"`
	AssetType               string     `json:"asset_type"`
	Vendor                  string     `json:"vendor"`
	Model                   string     `json:"model"`
	SerialNumber            string     `json:"serial_number"`
	MACAddress              string     `json:"mac_address"`
	IPAddress               string     `json:"ip_address"`
	Hostname                string     `json:"hostname"`
	OwnerDepartment         string     `json:"owner_department"`
	OwnerUser               string     `json:"owner_user"`
	Location                string     `json:"location"`
	Building                string     `json:"building"`
	Floor                   string     `json:"floor"`
	Room                    string     `json:"room"`
	Source                  string     `json:"source"`
	Status                  string     `json:"status"`
	InitialValue            float64    `json:"initial_value"`
	DepMonths               int        `json:"depreciation_months"`
	AccumulatedDepreciation float64    `json:"accumulated_depreciation"`
	ImpairmentProvision     float64    `json:"impairment_provision"`
	InServiceDate           *time.Time `json:"in_service_date"`
}

type UpdateRequest struct {
	AssetName               string     `json:"asset_name"`
	AssetHashID             string     `json:"asset_hash_id"`
	RFIDUID                 string     `json:"rfid_uid"`
	AssetType               string     `json:"asset_type"`
	Vendor                  string     `json:"vendor"`
	Model                   string     `json:"model"`
	SerialNumber            string     `json:"serial_number"`
	MACAddress              string     `json:"mac_address"`
	IPAddress               string     `json:"ip_address"`
	Hostname                string     `json:"hostname"`
	OwnerDepartment         string     `json:"owner_department"`
	OwnerUser               string     `json:"owner_user"`
	Location                string     `json:"location"`
	Building                string     `json:"building"`
	Floor                   string     `json:"floor"`
	Room                    string     `json:"room"`
	Source                  string     `json:"source"`
	TrustLevel              string     `json:"trust_level"`
	Status                  string     `json:"status"`
	InitialValue            float64    `json:"initial_value"`
	DepMonths               int        `json:"depreciation_months"`
	AccumulatedDepreciation float64    `json:"accumulated_depreciation"`
	ImpairmentProvision     float64    `json:"impairment_provision"`
	InServiceDate           *time.Time `json:"in_service_date"`
	DeactivationDate        *time.Time `json:"deactivation_date"`
}

type StatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type Query struct {
	Keyword   string
	Status    string
	AssetType string
	Building  string
	Floor     string
	Room      string
}

type InsuranceRequest struct {
	PolicyNo      string    `json:"insurance_policy_no" binding:"required"`
	AnnualPremium float64   `json:"annual_premium"`
	InsuredAmount float64   `json:"insured_amount"`
	PeriodStart   time.Time `json:"period_start" binding:"required"`
	PeriodEnd     time.Time `json:"period_end" binding:"required"`
	Operator      string    `json:"operator"`
}

type ImpairmentRequest struct {
	Reason            string  `json:"reason" binding:"required"`
	EvidenceFileHash  string  `json:"evidence_file_hash" binding:"required"`
	RecoverableAmount float64 `json:"recoverable_amount" binding:"required"`
	Reviewer          string  `json:"reviewer"`
}
