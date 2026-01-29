package repo

import (
	"context"

	"github.com/payments-core/internal/domain"
	"gorm.io/gorm"
)

type Account struct {
	AccountID domain.AccountID
	Balance   domain.Money
}

type AccountRepository interface {
	Create(ctx context.Context, db *gorm.DB, id domain.AccountID, initial domain.Money) error
	Get(ctx context.Context, db *gorm.DB, id domain.AccountID) (Account, error)
	GetForUpdate(ctx context.Context, tx *gorm.DB, id domain.AccountID) (Account, error)
	UpdateBalance(ctx context.Context, tx *gorm.DB, id domain.AccountID, newBalance domain.Money) error
}

type TransactionRepository interface {
	Insert(ctx context.Context, tx *gorm.DB, t Transaction) error
}

type Transaction struct {
	ID                   string
	SourceAccountID      domain.AccountID
	DestinationAccountID domain.AccountID
	Amount               domain.Money
}
