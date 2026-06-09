package middleware

import (
	"net/http"
	"strings"

	"asset-core/internal/module/auth"
	apperrors "asset-core/internal/pkg/errors"
	"asset-core/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

const ClaimsKey = "auth_claims"

func AuthRequired(service *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Fail(c, http.StatusUnauthorized, apperrors.CodeUnauthorized, "missing authorization token")
			c.Abort()
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		if token == "" || token == header {
			response.Fail(c, http.StatusUnauthorized, apperrors.CodeUnauthorized, "invalid authorization token")
			c.Abort()
			return
		}
		claims, err := service.Parse(token)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, apperrors.CodeUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}
		c.Set(ClaimsKey, claims)
		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, ok := c.Get(ClaimsKey)
		if !ok {
			response.Fail(c, http.StatusUnauthorized, apperrors.CodeUnauthorized, "missing authorization token")
			c.Abort()
			return
		}
		claims, ok := value.(*auth.Claims)
		if !ok || (!hasRole(claims.Roles, auth.RoleSuperAdmin) && !hasPermission(claims.Permissions, permission)) {
			response.Fail(c, http.StatusForbidden, apperrors.CodeForbidden, "permission denied")
			c.Abort()
			return
		}
		c.Next()
	}
}

func hasRole(roles []string, required string) bool {
	for _, role := range roles {
		if role == required {
			return true
		}
	}
	return false
}

func hasPermission(permissions []string, required string) bool {
	for _, permission := range permissions {
		if permission == required {
			return true
		}
	}
	return false
}
