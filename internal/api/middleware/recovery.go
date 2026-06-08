package middleware

import (
	"net/http"

	apperrors "asset-core/internal/pkg/errors"
	"asset-core/internal/pkg/logger"
	"asset-core/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func Recovery(log *logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Error("panic recovered", logger.Any("recovered", recovered))
		response.Fail(c, http.StatusInternalServerError, apperrors.CodeInternal, "internal error")
	})
}
