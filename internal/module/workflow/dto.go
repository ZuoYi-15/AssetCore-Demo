package workflow

type NodeRequest struct {
	NodeName     string `json:"node_name" binding:"required"`
	ApproverRole string `json:"approver_role" binding:"required"`
	SortOrder    int    `json:"sort_order"`
}

type SaveDefinitionRequest struct {
	FlowType string        `json:"flow_type" binding:"required"`
	Name     string        `json:"name" binding:"required"`
	Status   string        `json:"status"`
	Nodes    []NodeRequest `json:"nodes" binding:"required"`
}

type StartRequest struct {
	FlowType string                 `json:"flow_type" binding:"required"`
	AssetID  uint64                 `json:"asset_id"`
	Title    string                 `json:"title" binding:"required"`
	Payload  map[string]interface{} `json:"payload"`
}

type ApproveRequest struct {
	Action  string `json:"action" binding:"required"`
	Comment string `json:"comment"`
}

type Actor struct {
	UserID   uint64
	Username string
	Roles    []string
}
