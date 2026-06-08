package controller

import (
	"net/http"

	"asset-core/internal/module/data"
	apperrors "asset-core/internal/pkg/errors"
	"asset-core/internal/pkg/pagination"
	"asset-core/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type DataController struct {
	service *data.Service
}

func NewDataController(service *data.Service) *DataController {
	return &DataController{service: service}
}

func (ctl *DataController) CreateImportTask(c *gin.Context) {
	var req data.CreateImportTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, err.Error())
		return
	}
	item, err := ctl.service.CreateImportTask(req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Created(c, item)
}

func (ctl *DataController) ListImportTasks(c *gin.Context) {
	page := pagination.FromQuery(c)
	items, total, err := ctl.service.ListImportTasks(page.Offset(), page.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, pagination.Result{Items: items, Page: page.Page, PageSize: page.PageSize, Total: total})
}

func (ctl *DataController) GetImportTask(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	item, err := ctl.service.GetImportTask(id)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, item)
}

func (ctl *DataController) ImportErrors(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	items, err := ctl.service.ImportErrors(id)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, items)
}

func (ctl *DataController) ExportAssets(c *gin.Context) {
	response.OK(c, gin.H{"status": "pending", "message": "asset export task is reserved for MVP extension"})
}
