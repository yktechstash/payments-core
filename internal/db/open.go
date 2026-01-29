package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/payments-core/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(cfg config.Config) (*gorm.DB, error) {
	var gdb *gorm.DB
	var err error
	switch cfg.DBDriver {
	case "postgres":
		gdb, err = gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	case "sqlite":
		gdb, err = gorm.Open(sqlite.Open(cfg.DBDSN), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported DB_DRIVER: %q", cfg.DBDriver)
	}
	if err != nil {
		return nil, err
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}
	applyPoolDefaults(sqlDB)
	if err := ping(sqlDB); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return gdb, nil
}

func applyPoolDefaults(db *sql.DB) {
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)
}

func ping(db *sql.DB) error {
	return db.Ping()
}
