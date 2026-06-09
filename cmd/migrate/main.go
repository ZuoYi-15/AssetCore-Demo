package main

import (
	"asset-core/internal/config"
	"asset-core/internal/infrastructure/mysql"
	"asset-core/internal/module/asset"
	"asset-core/internal/module/audit"
	"asset-core/internal/module/auth"
	"asset-core/internal/module/data"
	"asset-core/internal/module/identity"
	"asset-core/internal/module/verification"
)

func main() {
	cfg := config.Load()
	db := mysql.New(cfg.MySQL)
	if err := db.AutoMigrate(
		&asset.Asset{},
		&asset.ChangeLog{},
		&identity.Identity{},
		&identity.Feature{},
		&verification.Task{},
		&verification.Conflict{},
		&data.ImportTask{},
		&data.ImportError{},
		&audit.Log{},
		&auth.User{},
		&auth.Role{},
		&auth.Permission{},
		&auth.UserRole{},
		&auth.RolePermission{},
	); err != nil {
		panic(err)
	}
}
