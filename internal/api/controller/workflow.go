package controller

import (
	"errors"
	"net/http"
	"strconv"

	"asset-core/internal/api/middleware"
	"asset-core/internal/module/auth"
	"asset-core/internal/module/workflow"
	apperrors "asset-core/internal/pkg/errors"
	"asset-core/internal/pkg/pagination"
	"asset-core/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type WorkflowController struct {
	service *workflow.Service
}

func NewWorkflowController(service *workflow.Service) *WorkflowController {
	return &WorkflowController{service: service}
}

func (ctl *WorkflowController) ListDefinitions(c *gin.Context) {
	items, err := ctl.service.ListDefinitions()
	if err != nil {
		handleWorkflowError(c, err)
		return
	}
	response.OK(c, items)
}

func (ctl *WorkflowController) SaveDefinition(c *gin.Context) {
	var req workflow.SaveDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, err.Error())
		return
	}
	item, err := ctl.service.SaveDefinition(req)
	if err != nil {
		handleWorkflowError(c, err)
		return
	}
	response.OK(c, item)
}

func (ctl *WorkflowController) DeleteDefinition(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := ctl.service.DeleteDefinition(id); err != nil {
		handleWorkflowError(c, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

func (ctl *WorkflowController) Start(c *gin.Context) {
	var req workflow.StartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, err.Error())
		return
	}
	item, err := ctl.service.Start(req, actorFromContext(c))
	if err != nil {
		handleWorkflowError(c, err)
		return
	}
	response.Created(c, item)
}

func (ctl *WorkflowController) ListInstances(c *gin.Context) {
	page := pagination.FromQuery(c)
	var assetID uint64
	if value := c.Query("asset_id"); value != "" {
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, "invalid asset_id")
			return
		}
		assetID = parsed
	}
	items, total, err := ctl.service.ListInstances(assetID, page.Offset(), page.PageSize)
	if err != nil {
		handleWorkflowError(c, err)
		return
	}
	response.OK(c, pagination.Result{Items: items, Page: page.Page, PageSize: page.PageSize, Total: total})
}

func (ctl *WorkflowController) ListTasks(c *gin.Context) {
	page := pagination.FromQuery(c)
	items, total, err := ctl.service.ListTasks(c.Query("status"), page.Offset(), page.PageSize)
	if err != nil {
		handleWorkflowError(c, err)
		return
	}
	response.OK(c, pagination.Result{Items: items, Page: page.Page, PageSize: page.PageSize, Total: total})
}

func (ctl *WorkflowController) Approve(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req workflow.ApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, err.Error())
		return
	}
	item, err := ctl.service.Approve(id, req, actorFromContext(c))
	if err != nil {
		handleWorkflowError(c, err)
		return
	}
	response.OK(c, item)
}

func actorFromContext(c *gin.Context) workflow.Actor {
	value, _ := c.Get(middleware.ClaimsKey)
	claims, _ := value.(*auth.Claims)
	if claims == nil {
		return workflow.Actor{}
	}
	return workflow.Actor{
		UserID:   claims.UserID,
		Username: claims.Username,
		Roles:    claims.Roles,
	}
}

func handleWorkflowError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, workflow.ErrDefinitionEmpty), errors.Is(err, workflow.ErrInvalidAction):
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, err.Error())
	case errors.Is(err, workflow.ErrDefinitionInactive), errors.Is(err, workflow.ErrTaskNotPending), errors.Is(err, workflow.ErrDefinitionActive):
		response.Fail(c, http.StatusConflict, apperrors.CodeAssetConflict, err.Error())
	case errors.Is(err, workflow.ErrApproverRole):
		response.Fail(c, http.StatusForbidden, apperrors.CodeForbidden, err.Error())
	default:
		handleError(c, err)
	}
}
