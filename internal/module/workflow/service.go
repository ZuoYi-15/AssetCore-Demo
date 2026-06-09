package workflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"asset-core/internal/module/asset"
	"asset-core/internal/module/auth"

	"gorm.io/gorm"
)

var (
	ErrDefinitionInactive = errors.New("workflow definition is inactive")
	ErrDefinitionEmpty    = errors.New("workflow definition must contain at least one node")
	ErrDefinitionActive   = errors.New("only inactive workflow definitions can be deleted")
	ErrInvalidAction      = errors.New("invalid workflow action")
	ErrTaskNotPending     = errors.New("workflow task is not pending")
	ErrApproverRole       = errors.New("current user does not match approver role")
)

type Service struct {
	repo         *Repository
	assetService AssetStatusUpdater
}

type AssetStatusUpdater interface {
	Update(id uint64, req asset.UpdateRequest) (*asset.Asset, error)
	UpdateStatus(id uint64, status string) (*asset.Asset, error)
	RecordChangeLog(assetID uint64, field, oldValue, newValue, operator string) error
	HasChangeLog(assetID uint64, field, newValue string) (bool, error)
}

func NewService(repo *Repository, assetService AssetStatusUpdater) *Service {
	return &Service{repo: repo, assetService: assetService}
}

func (s *Service) Bootstrap() error {
	if err := s.repo.AutoMigrate(); err != nil {
		return err
	}
	count, err := s.repo.CountDefinitions()
	if err != nil {
		return err
	}
	if count == 0 {
		defaults := []SaveDefinitionRequest{
			{FlowType: FlowPurchase, Name: "采购审批", Status: StatusActive, Nodes: []NodeRequest{{NodeName: "管理员审批", ApproverRole: auth.RoleAdmin, SortOrder: 1}}},
			{FlowType: FlowTransfer, Name: "调拨审批", Status: StatusActive, Nodes: []NodeRequest{{NodeName: "管理员审批", ApproverRole: auth.RoleAdmin, SortOrder: 1}}},
			{FlowType: FlowRetire, Name: "报废审批", Status: StatusActive, Nodes: []NodeRequest{{NodeName: "管理员审批", ApproverRole: auth.RoleAdmin, SortOrder: 1}}},
		}
		for _, item := range defaults {
			if _, err := s.SaveDefinition(item); err != nil {
				return err
			}
		}
	}
	if err := s.syncApprovedRetireAssets(); err != nil {
		return err
	}
	return s.syncApprovedTransferLogs()
}

func (s *Service) SaveDefinition(req SaveDefinitionRequest) (*Definition, error) {
	req.FlowType = strings.TrimSpace(req.FlowType)
	if req.Status == "" {
		req.Status = StatusActive
	}
	nodes := normalizeNodes(req.Nodes)
	if len(nodes) == 0 {
		return nil, ErrDefinitionEmpty
	}

	var saved *Definition
	err := s.repo.Transaction(func(tx *Repository) error {
		definition, err := tx.FindDefinitionByType(req.FlowType)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			definition = &Definition{FlowType: req.FlowType}
			err = nil
		}
		if err != nil {
			return err
		}
		definition.Name = strings.TrimSpace(req.Name)
		definition.Status = req.Status
		if definition.ID == 0 {
			if err := tx.CreateDefinition(definition); err != nil {
				return err
			}
		} else if err := tx.UpdateDefinition(definition); err != nil {
			return err
		}
		if err := tx.DeleteNodes(definition.ID); err != nil {
			return err
		}
		items := make([]Node, 0, len(nodes))
		for _, node := range nodes {
			items = append(items, Node{
				DefinitionID: definition.ID,
				NodeName:     node.NodeName,
				ApproverRole: node.ApproverRole,
				SortOrder:    node.SortOrder,
			})
		}
		if err := tx.CreateNodes(items); err != nil {
			return err
		}
		saved = definition
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.repo.FindDefinitionByType(saved.FlowType)
}

func (s *Service) ListDefinitions() ([]Definition, error) {
	return s.repo.ListDefinitions()
}

func (s *Service) DeleteDefinition(id uint64) error {
	return s.repo.Transaction(func(tx *Repository) error {
		definition, err := tx.FindDefinitionByID(id)
		if err != nil {
			return err
		}
		if definition.Status != StatusInactive {
			return ErrDefinitionActive
		}
		if err := tx.DeleteNodes(definition.ID); err != nil {
			return err
		}
		return tx.DeleteDefinition(definition.ID)
	})
}

func (s *Service) Start(req StartRequest, actor Actor) (*Instance, error) {
	definition, err := s.repo.FindDefinitionByType(req.FlowType)
	if err != nil {
		return nil, err
	}
	if definition.Status != StatusActive {
		return nil, ErrDefinitionInactive
	}
	if len(definition.Nodes) == 0 {
		return nil, ErrDefinitionEmpty
	}
	payload, _ := json.Marshal(req.Payload)
	var created *Instance
	err = s.repo.Transaction(func(tx *Repository) error {
		instance := &Instance{
			DefinitionID:  definition.ID,
			FlowType:      definition.FlowType,
			AssetID:       req.AssetID,
			Title:         strings.TrimSpace(req.Title),
			Status:        InstancePending,
			ApplicantID:   actor.UserID,
			ApplicantName: actor.Username,
			Payload:       string(payload),
		}
		if err := tx.CreateInstance(instance); err != nil {
			return err
		}
		first := definition.Nodes[0]
		task := &Task{
			InstanceID:   instance.ID,
			NodeID:       first.ID,
			NodeName:     first.NodeName,
			ApproverRole: first.ApproverRole,
			SortOrder:    first.SortOrder,
			Status:       TaskPending,
		}
		if err := tx.CreateTask(task); err != nil {
			return err
		}
		instance.CurrentTaskID = task.ID
		if err := tx.UpdateInstance(instance); err != nil {
			return err
		}
		created = instance
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.repo.FindInstance(created.ID)
}

func (s *Service) Approve(taskID uint64, req ApproveRequest, actor Actor) (*Instance, error) {
	req.Action = strings.TrimSpace(req.Action)
	if req.Action != "approve" && req.Action != "reject" {
		return nil, ErrInvalidAction
	}
	var instanceID uint64
	var shouldRetireAsset bool
	var retireAssetID uint64
	var approvedTransfer *Instance
	err := s.repo.Transaction(func(tx *Repository) error {
		task, err := tx.FindTask(taskID)
		if err != nil {
			return err
		}
		if task.Status != TaskPending {
			return ErrTaskNotPending
		}
		if !hasRole(actor.Roles, task.ApproverRole) && !hasRole(actor.Roles, auth.RoleSuperAdmin) {
			return ErrApproverRole
		}
		now := time.Now()
		task.ApproverID = actor.UserID
		task.ApproverName = actor.Username
		task.Comment = req.Comment
		task.CompletedAt = &now
		if req.Action == "reject" {
			task.Status = TaskRejected
			if err := tx.UpdateTask(task); err != nil {
				return err
			}
			task.Instance.Status = InstanceRejected
			task.Instance.CurrentTaskID = 0
			task.Instance.CompletedAt = &now
			instanceID = task.Instance.ID
			return tx.UpdateInstance(task.Instance)
		}

		task.Status = TaskApproved
		if err := tx.UpdateTask(task); err != nil {
			return err
		}
		definition, err := tx.FindDefinitionByType(task.Instance.FlowType)
		if err != nil {
			return err
		}
		next := nextNode(definition.Nodes, task.SortOrder)
		if next == nil {
			task.Instance.Status = InstanceApproved
			task.Instance.CurrentTaskID = 0
			task.Instance.CompletedAt = &now
			instanceID = task.Instance.ID
			if task.Instance.FlowType == FlowRetire && task.Instance.AssetID > 0 {
				shouldRetireAsset = true
				retireAssetID = task.Instance.AssetID
			}
			if task.Instance.FlowType == FlowTransfer && task.Instance.AssetID > 0 {
				copy := *task.Instance
				approvedTransfer = &copy
			}
			return tx.UpdateInstance(task.Instance)
		}
		nextTask := &Task{
			InstanceID:   task.Instance.ID,
			NodeID:       next.ID,
			NodeName:     next.NodeName,
			ApproverRole: next.ApproverRole,
			SortOrder:    next.SortOrder,
			Status:       TaskPending,
		}
		if err := tx.CreateTask(nextTask); err != nil {
			return err
		}
		task.Instance.CurrentTaskID = nextTask.ID
		instanceID = task.Instance.ID
		return tx.UpdateInstance(task.Instance)
	})
	if err != nil {
		return nil, err
	}
	if shouldRetireAsset && s.assetService != nil {
		if _, err := s.assetService.UpdateStatus(retireAssetID, asset.StatusRetired); err != nil {
			return nil, err
		}
	}
	if approvedTransfer != nil {
		if err := s.applyApprovedTransfer(*approvedTransfer); err != nil {
			return nil, err
		}
	}
	return s.repo.FindInstance(instanceID)
}

func (s *Service) ListInstances(assetID uint64, offset, limit int) ([]Instance, int64, error) {
	return s.repo.ListInstances(assetID, offset, limit)
}

func (s *Service) ListTasks(status string, offset, limit int) ([]Task, int64, error) {
	return s.repo.ListTasks(status, offset, limit)
}

func (s *Service) syncApprovedRetireAssets() error {
	if s.assetService == nil {
		return nil
	}
	items, err := s.repo.ListApprovedRetireInstances()
	if err != nil {
		return err
	}
	for _, item := range items {
		if _, err := s.assetService.UpdateStatus(item.AssetID, asset.StatusRetired); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) syncApprovedTransferLogs() error {
	if s.assetService == nil {
		return nil
	}
	items, err := s.repo.ListApprovedTransferInstances()
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := s.applyApprovedTransfer(item); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) applyApprovedTransfer(item Instance) error {
	if s.assetService == nil || item.AssetID == 0 {
		return nil
	}
	newValue := workflowChangeValue(item.ID)
	exists, err := s.assetService.HasChangeLog(item.AssetID, "workflow_transfer", newValue)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	payload := parseTransferPayload(item.Payload)
	update := asset.UpdateRequest{
		OwnerDepartment: payload.TargetDepartment,
		OwnerUser:       payload.TargetOwner,
		Location:        payload.TargetLocation,
	}
	if update.OwnerDepartment != "" || update.OwnerUser != "" || update.Location != "" {
		if _, err := s.assetService.Update(item.AssetID, update); err != nil {
			return err
		}
	}
	return s.assetService.RecordChangeLog(item.AssetID, "workflow_transfer", item.Title, newValue, "workflow")
}

type transferPayload struct {
	TargetDepartment string
	TargetOwner      string
	TargetLocation   string
}

func parseTransferPayload(raw string) transferPayload {
	var data map[string]interface{}
	_ = json.Unmarshal([]byte(raw), &data)
	return transferPayload{
		TargetDepartment: stringValue(data["target_department"]),
		TargetOwner:      stringValue(data["target_owner"]),
		TargetLocation:   stringValue(data["target_location"]),
	}
}

func stringValue(value interface{}) string {
	if value == nil {
		return ""
	}
	if s, ok := value.(string); ok {
		return strings.TrimSpace(s)
	}
	return strings.TrimSpace(fmt.Sprint(value))
}

func workflowChangeValue(instanceID uint64) string {
	return "approved:" + strconv.FormatUint(instanceID, 10)
}

func normalizeNodes(nodes []NodeRequest) []NodeRequest {
	items := make([]NodeRequest, 0, len(nodes))
	for i, node := range nodes {
		node.NodeName = strings.TrimSpace(node.NodeName)
		node.ApproverRole = strings.TrimSpace(node.ApproverRole)
		if node.NodeName == "" || node.ApproverRole == "" {
			continue
		}
		if node.SortOrder <= 0 {
			node.SortOrder = i + 1
		}
		items = append(items, node)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].SortOrder < items[j].SortOrder
	})
	for i := range items {
		items[i].SortOrder = i + 1
	}
	return items
}

func nextNode(nodes []Node, currentOrder int) *Node {
	for i := range nodes {
		if nodes[i].SortOrder > currentOrder {
			return &nodes[i]
		}
	}
	return nil
}

func hasRole(roles []string, required string) bool {
	for _, role := range roles {
		if role == required {
			return true
		}
	}
	return false
}
