package controller

import (
	"errors"
	"net/http"
	"strconv"

	"asset-core/internal/module/identity"
	apperrors "asset-core/internal/pkg/errors"
	"asset-core/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func parseID(c *gin.Context, name string) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, "invalid id")
		return 0, false
	}
	return id, true
}

func handleError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, apperrors.CodeNotFound, "resource not found")
		return
	}
	if errors.Is(err, identity.ErrIdentityAlreadyBound) || errors.Is(err, identity.ErrAssetAlreadyBound) {
		response.Fail(c, http.StatusConflict, apperrors.CodeIdentityConflict, err.Error())
		return
	}
	response.Fail(c, http.StatusInternalServerError, apperrors.CodeInternal, err.Error())
}
