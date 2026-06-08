package mysql

import (
	"time"

	"asset-core/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(cfg config.MySQLConfig) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}
