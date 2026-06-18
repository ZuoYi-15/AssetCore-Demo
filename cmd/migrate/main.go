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
	"asset-core/internal/module/workflow"
)

func main() {
	cfg := config.Load()
	db := mysql.New(cfg.MySQL)
	if err := db.AutoMigrate(
		&asset.Asset{},
		&asset.ChangeLog{},
		&asset.Insurance{},
		&asset.ImpairmentRecord{},
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
		&auth.UserPermission{},
		&auth.RolePermission{},
		&workflow.Definition{},
		&workflow.Node{},
		&workflow.Instance{},
		&workflow.Task{},
	); err != nil {
		panic(err)
	}
}
