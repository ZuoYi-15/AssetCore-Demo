package controller

import (
	"net/http"

	"asset-core/internal/module/asset"
	apperrors "asset-core/internal/pkg/errors"
	"asset-core/internal/pkg/pagination"
	"asset-core/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type AssetController struct {
	service *asset.Service
}

func NewAssetController(service *asset.Service) *AssetController {
	return &AssetController{service: service}
}

func (ctl *AssetController) Create(c *gin.Context) {
	var req asset.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, err.Error())
		return
	}
	item, err := ctl.service.Create(req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Created(c, item)
}

func (ctl *AssetController) List(c *gin.Context) {
	page := pagination.FromQuery(c)
	items, total, err := ctl.service.List(asset.Query{
		Keyword:   c.Query("keyword"),
		Status:    c.Query("status"),
		AssetType: c.Query("asset_type"),
	}, page.Offset(), page.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, pagination.Result{Items: items, Page: page.Page, PageSize: page.PageSize, Total: total})
}

func (ctl *AssetController) Get(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	item, err := ctl.service.Get(id)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, item)
}

func (ctl *AssetController) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req asset.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, err.Error())
		return
	}
	item, err := ctl.service.Update(id, req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, item)
}

func (ctl *AssetController) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := ctl.service.Delete(id); err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

func (ctl *AssetController) ChangeStatus(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req asset.StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, err.Error())
		return
	}
	item, err := ctl.service.UpdateStatus(id, req.Status)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, item)
}

func (ctl *AssetController) ChangeLogs(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	items, err := ctl.service.ChangeLogs(id)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, items)
}
