package db

import (
	"fmt"

	"github.com/payments-core/internal/config"
	"gorm.io/gorm"
)

func Migrate(cfg config.Config, db *gorm.DB) error {
	var stmts []string
	switch cfg.DBDriver {
	case "postgres":
		stmts = postgresSchema()
	case "sqlite":
		stmts = sqliteSchema()
	default:
		return fmt.Errorf("unsupported DB_DRIVER: %q", cfg.DBDriver)
	}

	for _, s := range stmts {
		if err := db.Exec(s).Error; err != nil {
			return fmt.Errorf("exec schema: %w", err)
		}
	}
	return nil
}

func postgresSchema() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS accounts (
			account_id BIGINT PRIMARY KEY,
			balance NUMERIC(36, 18) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);`,
		`CREATE TABLE IF NOT EXISTS transactions (
			id UUID PRIMARY KEY,
			source_account_id BIGINT NOT NULL REFERENCES accounts(account_id),
			destination_account_id BIGINT NOT NULL REFERENCES accounts(account_id),
			amount NUMERIC(36, 18) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_source_created_at ON transactions(source_account_id, created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_dest_created_at ON transactions(destination_account_id, created_at);`,
	}
}

func sqliteSchema() []string {
	return []string{
		`PRAGMA journal_mode=WAL;`,
		`CREATE TABLE IF NOT EXISTS accounts (
			account_id INTEGER PRIMARY KEY,
			balance TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT (datetime('now'))
		);`,
		`CREATE TABLE IF NOT EXISTS transactions (
			id TEXT PRIMARY KEY,
			source_account_id INTEGER NOT NULL,
			destination_account_id INTEGER NOT NULL,
			amount TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			FOREIGN KEY(source_account_id) REFERENCES accounts(account_id),
			FOREIGN KEY(destination_account_id) REFERENCES accounts(account_id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_source_created_at ON transactions(source_account_id, created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_dest_created_at ON transactions(destination_account_id, created_at);`,
	}
}
