package controller

import (
	"errors"
	"net/http"

	"asset-core/internal/api/middleware"
	"asset-core/internal/module/auth"
	apperrors "asset-core/internal/pkg/errors"
	"asset-core/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	service *auth.Service
}

func NewAuthController(service *auth.Service) *AuthController {
	return &AuthController{service: service}
}

func (ctl *AuthController) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, err.Error())
		return
	}
	result, err := ctl.service.Login(req)
	if err != nil {
		handleAuthError(c, err)
		return
	}
	response.OK(c, result)
}

func (ctl *AuthController) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, err.Error())
		return
	}
	user, err := ctl.service.Register(req)
	if err != nil {
		handleAuthError(c, err)
		return
	}
	response.Created(c, user)
}

func (ctl *AuthController) Me(c *gin.Context) {
	value, ok := c.Get(middleware.ClaimsKey)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, apperrors.CodeUnauthorized, "missing authorization token")
		return
	}
	claims, ok := value.(*auth.Claims)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, apperrors.CodeUnauthorized, "invalid authorization token")
		return
	}
	profile, err := ctl.service.Profile(claims.UserID)
	if err != nil {
		handleAuthError(c, err)
		return
	}
	response.OK(c, profile)
}

func handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, auth.ErrInvalidCredential):
		response.Fail(c, http.StatusUnauthorized, apperrors.CodeUnauthorized, "invalid username or password")
	case errors.Is(err, auth.ErrInactiveUser):
		response.Fail(c, http.StatusForbidden, apperrors.CodeForbidden, "user is disabled")
	case errors.Is(err, auth.ErrUnsupportedRole):
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, "unsupported role")
	case errors.Is(err, auth.ErrWeakPassword):
		response.Fail(c, http.StatusBadRequest, apperrors.CodeInvalidParameter, "password must be at least 8 characters")
	case errors.Is(err, auth.ErrUserExists):
		response.Fail(c, http.StatusConflict, apperrors.CodeAssetConflict, "username already exists")
	case errors.Is(err, gorm.ErrRecordNotFound):
		response.Fail(c, http.StatusNotFound, apperrors.CodeNotFound, "resource not found")
	default:
		response.Fail(c, http.StatusInternalServerError, apperrors.CodeInternal, err.Error())
	}
}
