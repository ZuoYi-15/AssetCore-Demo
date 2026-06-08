package controller

import (
	"net/http"

	"asset-core/internal/module/verification"
	apperrors "asset-core/internal/pkg/errors"
	"asset-core/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type VerificationController struct {
	service *verification.Service
}

func NewVerificationController(service *verification.Service) *VerificationController {
	return &VerificationController{service: service}
}

func (ctl *VerificationController) Create(c *gin.Context) {
	var req verification.CreateRequest
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

func (ctl *VerificationController) Get(c *gin.Context) {
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

func (ctl *VerificationController) VerifyAsset(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	item, err := ctl.service.Create(verification.CreateRequest{AssetID: id})
	if err != nil {
		handleError(c, err)
		return
	}
	response.Created(c, item)
}

func (ctl *VerificationController) LatestByAsset(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	item, err := ctl.service.LatestByAsset(id)
	if err != nil {
		handleError(c, err)
		return
	}
	response.OK(c, item)
}
