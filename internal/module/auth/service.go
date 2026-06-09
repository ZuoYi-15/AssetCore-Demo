package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"asset-core/internal/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredential = errors.New("invalid username or password")
	ErrInactiveUser      = errors.New("user is disabled")
	ErrUnsupportedRole   = errors.New("unsupported role")
	ErrUserExists        = errors.New("username already exists")
	ErrWeakPassword      = errors.New("password must be at least 8 characters")
)

type Service struct {
	repo *Repository
	cfg  config.JWTConfig
}

func NewService(repo *Repository, cfg config.JWTConfig) *Service {
	return &Service{repo: repo, cfg: cfg}
}

func (s *Service) Bootstrap() error {
	if err := s.repo.AutoMigrate(); err != nil {
		return err
	}
	if err := s.seedRolesAndPermissions(); err != nil {
		return err
	}
	count, err := s.repo.CountUsers()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	user, err := s.Register(RegisterRequest{
		Username:    "superadmin",
		Password:    "Admin@123456",
		DisplayName: "Super Administrator",
		RoleCode:    RoleSuperAdmin,
	})
	if err != nil {
		return err
	}
	_ = user
	return nil
}

func (s *Service) Register(req RegisterRequest) (*UserProfile, error) {
	req.RoleCode = strings.TrimSpace(req.RoleCode)
	if req.RoleCode != RoleSuperAdmin && req.RoleCode != RoleAdmin && req.RoleCode != RoleUser {
		return nil, ErrUnsupportedRole
	}
	username := strings.TrimSpace(req.Username)
	if username == "" || len(req.Password) < 8 {
		return nil, ErrWeakPassword
	}
	if _, err := s.repo.FindUserByUsername(username); err == nil {
		return nil, ErrUserExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &User{
		Username:     username,
		PasswordHash: string(passwordHash),
		DisplayName:  strings.TrimSpace(req.DisplayName),
		Status:       "active",
	}
	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	if err := s.repo.AssignRole(user.ID, req.RoleCode); err != nil {
		return nil, err
	}
	return s.Profile(user.ID)
}

func (s *Service) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.FindUserByUsername(strings.TrimSpace(req.Username))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidCredential
	}
	if err != nil {
		return nil, err
	}
	if user.Status != "active" {
		return nil, ErrInactiveUser
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredential
	}
	profile, err := s.Profile(user.ID)
	if err != nil {
		return nil, err
	}
	token, err := SignToken(Claims{
		UserID:      profile.ID,
		Username:    profile.Username,
		Roles:       profile.Roles,
		Permissions: profile.Permissions,
		ExpiresAt:   time.Now().Add(24 * time.Hour).Unix(),
		Issuer:      s.cfg.Issuer,
	}, s.cfg.Secret)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{Token: token, User: *profile}, nil
}

func (s *Service) Profile(userID uint64) (*UserProfile, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	roles, err := s.repo.UserRoles(user.ID)
	if err != nil {
		return nil, err
	}
	permissions, err := s.repo.UserPermissions(user.ID)
	if err != nil {
		return nil, err
	}
	return &UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

func (s *Service) Parse(token string) (*Claims, error) {
	claims, err := ParseToken(token, s.cfg.Secret)
	if err != nil {
		return nil, err
	}
	if claims.Issuer != s.cfg.Issuer || claims.ExpiresAt < time.Now().Unix() {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (s *Service) seedRolesAndPermissions() error {
	roles := []Role{
		{Code: RoleSuperAdmin, Name: "超级管理员", Description: "拥有系统全部权限"},
		{Code: RoleAdmin, Name: "管理员", Description: "可管理资产台账"},
		{Code: RoleUser, Name: "普通用户", Description: "仅可查看资产台账"},
	}
	for _, role := range roles {
		if err := s.repo.UpsertRole(role); err != nil {
			return err
		}
	}

	permissions := []Permission{
		{Code: PermissionAssetRead, Name: "查看资产", Resource: "asset", Action: "read"},
		{Code: PermissionAssetCreate, Name: "新增资产", Resource: "asset", Action: "create"},
		{Code: PermissionAssetUpdate, Name: "编辑资产", Resource: "asset", Action: "update"},
		{Code: PermissionAssetDelete, Name: "删除资产", Resource: "asset", Action: "delete"},
		{Code: PermissionUserCreate, Name: "注册账号", Resource: "user", Action: "create"},
	}
	for _, permission := range permissions {
		if err := s.repo.UpsertPermission(permission); err != nil {
			return err
		}
	}

	rolePermissions := map[string][]string{
		RoleSuperAdmin: {PermissionAssetRead, PermissionAssetCreate, PermissionAssetUpdate, PermissionAssetDelete, PermissionUserCreate},
		RoleAdmin:      {PermissionAssetRead, PermissionAssetCreate, PermissionAssetUpdate, PermissionAssetDelete},
		RoleUser:       {PermissionAssetRead},
	}
	for role, permissionCodes := range rolePermissions {
		for _, permissionCode := range permissionCodes {
			if err := s.repo.AssignPermission(role, permissionCode); err != nil {
				return fmt.Errorf("assign %s to %s: %w", permissionCode, role, err)
			}
		}
	}
	return nil
}
