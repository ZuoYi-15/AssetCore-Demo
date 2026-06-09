package auth

import "time"

const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleUser       = "user"

	PermissionAssetRead   = "asset:read"
	PermissionAssetCreate = "asset:create"
	PermissionAssetUpdate = "asset:update"
	PermissionAssetDelete = "asset:delete"
	PermissionUserCreate  = "user:create"

	PermissionWorkflowConfig  = "workflow:config"
	PermissionWorkflowStart   = "workflow:start"
	PermissionWorkflowApprove = "workflow:approve"
)

type User struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	DisplayName  string    `gorm:"size:128" json:"display_name"`
	Email        string    `gorm:"size:128;index" json:"email"`
	Status       string    `gorm:"size:32;index;not null;default:active" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string { return "user_account" }

type Role struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"size:128;not null" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Role) TableName() string { return "auth_role" }

type Permission struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"size:128;uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"size:128;not null" json:"name"`
	Resource    string    `gorm:"size:64;index" json:"resource"`
	Action      string    `gorm:"size:64" json:"action"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Permission) TableName() string { return "auth_permission" }

type UserRole struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"uniqueIndex:idx_user_role;not null" json:"user_id"`
	RoleID    uint64    `gorm:"uniqueIndex:idx_user_role;not null" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserRole) TableName() string { return "auth_user_role" }

type UserPermission struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	UserID       uint64    `gorm:"uniqueIndex:idx_user_permission;not null" json:"user_id"`
	PermissionID uint64    `gorm:"uniqueIndex:idx_user_permission;not null" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (UserPermission) TableName() string { return "auth_user_permission" }

type RolePermission struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	RoleID       uint64    `gorm:"uniqueIndex:idx_role_permission;not null" json:"role_id"`
	PermissionID uint64    `gorm:"uniqueIndex:idx_role_permission;not null" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (RolePermission) TableName() string { return "auth_role_permission" }
