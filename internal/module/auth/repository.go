package auth

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&User{}, &Role{}, &Permission{}, &UserRole{}, &UserPermission{}, &RolePermission{})
}

func (r *Repository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *Repository) UpdateUser(user *User) error {
	return r.db.Save(user).Error
}

func (r *Repository) ListUsers() ([]User, error) {
	var users []User
	err := r.db.Order("id ASC").Find(&users).Error
	return users, err
}

func (r *Repository) FindUserByUsername(username string) (*User, error) {
	var user User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindUserByID(id uint64) (*User, error) {
	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CountUsers() (int64, error) {
	var count int64
	err := r.db.Model(&User{}).Count(&count).Error
	return count, err
}

func (r *Repository) UpsertRole(role Role) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "description", "updated_at"}),
	}).Create(&role).Error
}

func (r *Repository) UpsertPermission(permission Permission) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "resource", "action", "description", "updated_at"}),
	}).Create(&permission).Error
}

func (r *Repository) FindRoleByCode(code string) (*Role, error) {
	var role Role
	if err := r.db.Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) FindPermissionByCode(code string) (*Permission, error) {
	var permission Permission
	if err := r.db.Where("code = ?", code).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *Repository) ListPermissions() ([]Permission, error) {
	var items []Permission
	err := r.db.Order("id ASC").Find(&items).Error
	return items, err
}

func (r *Repository) AssignRole(userID uint64, roleCode string) error {
	role, err := r.FindRoleByCode(roleCode)
	if err != nil {
		return err
	}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&UserRole{
		UserID: userID,
		RoleID: role.ID,
	}).Error
}

func (r *Repository) ReplaceRole(userID uint64, roleCode string) error {
	role, err := r.FindRoleByCode(roleCode)
	if err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&UserRole{}).Error; err != nil {
			return err
		}
		return tx.Create(&UserRole{UserID: userID, RoleID: role.ID}).Error
	})
}

func (r *Repository) ReplacePermissions(userID uint64, permissionCodes []string) error {
	permissions := make([]Permission, 0, len(permissionCodes))
	for _, code := range permissionCodes {
		permission, err := r.FindPermissionByCode(code)
		if err != nil {
			return err
		}
		permissions = append(permissions, *permission)
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&UserPermission{}).Error; err != nil {
			return err
		}
		for _, permission := range permissions {
			if err := tx.Create(&UserPermission{UserID: userID, PermissionID: permission.ID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) AssignPermission(roleCode, permissionCode string) error {
	role, err := r.FindRoleByCode(roleCode)
	if err != nil {
		return err
	}
	permission, err := r.FindPermissionByCode(permissionCode)
	if err != nil {
		return err
	}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&RolePermission{
		RoleID:       role.ID,
		PermissionID: permission.ID,
	}).Error
}

func (r *Repository) UserRoles(userID uint64) ([]string, error) {
	var roles []string
	err := r.db.Table("auth_role").
		Select("auth_role.code").
		Joins("JOIN auth_user_role ON auth_user_role.role_id = auth_role.id").
		Where("auth_user_role.user_id = ?", userID).
		Order("auth_role.id ASC").
		Scan(&roles).Error
	return roles, err
}

func (r *Repository) UserPermissions(userID uint64) ([]string, error) {
	var permissions []string
	rolePermissions := r.db.Table("auth_permission").
		Select("DISTINCT auth_permission.code").
		Joins("JOIN auth_role_permission ON auth_role_permission.permission_id = auth_permission.id").
		Joins("JOIN auth_user_role ON auth_user_role.role_id = auth_role_permission.role_id").
		Where("auth_user_role.user_id = ?", userID).
		Scan(&permissions).Error
	if rolePermissions != nil {
		return nil, rolePermissions
	}
	var direct []string
	err := r.db.Table("auth_permission").
		Select("DISTINCT auth_permission.code").
		Joins("JOIN auth_user_permission ON auth_user_permission.permission_id = auth_permission.id").
		Where("auth_user_permission.user_id = ?", userID).
		Scan(&direct).Error
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	merged := make([]string, 0, len(permissions)+len(direct))
	for _, permission := range append(permissions, direct...) {
		if !seen[permission] {
			seen[permission] = true
			merged = append(merged, permission)
		}
	}
	permissions = merged
	return permissions, err
}

func (r *Repository) UserDirectPermissions(userID uint64) ([]string, error) {
	var permissions []string
	err := r.db.Table("auth_permission").
		Select("DISTINCT auth_permission.code").
		Joins("JOIN auth_user_permission ON auth_user_permission.permission_id = auth_permission.id").
		Where("auth_user_permission.user_id = ?", userID).
		Order("auth_permission.code ASC").
		Scan(&permissions).Error
	return permissions, err
}

func IsDuplicate(err error) bool {
	return err != nil && !errors.Is(err, gorm.ErrRecordNotFound)
}
