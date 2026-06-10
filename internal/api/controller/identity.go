package controller

import (
	"net/http"

	"asset-core/internal/module/identity"
	apperrors "asset-core/internal/pkg/errors"
	"asset-core/internal/pkg/pagination"
	"asset-core/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type IdentityController struct {
	service *identity.Service
}

func NewIdentityController(service *identity.Service) *IdentityController {
	return &IdentityController{service: service}
}

func (ctl *IdentityController) Generate(c *gin.Context) {
	var req identity.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, err.Error())
		return
	}
	item, err := ctl.service.Generate(req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Created(c, item)
}

func (ctl *IdentityController) List(c *gin.Context) {
	page := pagination.FromQuery(c)
	items, total, err := ctl.service.List(identity.Query{
		Keyword: c.Query("keyword"),
		Status:  c.Query("status"),
	}, page.Offset(), page.PageSize)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, pagination.Result{Items: items, Page: page.Page, PageSize: page.PageSize, Total: total})
}

func (ctl *IdentityController) Get(c *gin.Context) {
	item, err := ctl.service.Get(c.Param("identity_id"))
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, item)
}

func (ctl *IdentityController) Bind(c *gin.Context) {
	var req identity.BindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, err.Error())
		return
	}
	item, err := ctl.service.Bind(c.Param("identity_id"), req.AssetID)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, item)
}

func (ctl *IdentityController) Unbind(c *gin.Context) {
	item, err := ctl.service.Unbind(c.Param("identity_id"))
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, item)
}

func (ctl *IdentityController) Features(c *gin.Context) {
	items, err := ctl.service.Features(c.Param("identity_id"))
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, items)
}
