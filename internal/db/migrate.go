package db

import (
	"github.com/payments-core/internal/config"
	"github.com/payments-core/internal/db/dto"
	"gorm.io/gorm"
)

func Migrate(cfg config.Config, db *gorm.DB) error {
	return db.AutoMigrate(&dto.Account{}, &dto.Transaction{})
}
