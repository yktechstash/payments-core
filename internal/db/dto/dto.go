package dto

import "time"

type Account struct {
	AccountID int64     `gorm:"column:account_id;primaryKey"`
	Balance   string    `gorm:"column:balance;type:numeric(36,18)"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Account) TableName() string { return "accounts" }

type Transaction struct {
	ID                   string    `gorm:"column:id;primaryKey"`
	SourceAccountID      int64     `gorm:"column:source_account_id;index:idx_transactions_source_created_at,priority:1"`
	DestinationAccountID int64     `gorm:"column:destination_account_id;index:idx_transactions_dest_created_at,priority:1"`
	Amount               string    `gorm:"column:amount;type:numeric(36,18)"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime;index:idx_transactions_source_created_at,priority:2;index:idx_transactions_dest_created_at,priority:2"`
}

func (Transaction) TableName() string { return "transactions" }
