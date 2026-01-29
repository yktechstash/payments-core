package dal

import (
	"context"

	"github.com/payments-core/internal/db/dto"
	"github.com/payments-core/internal/domain"
	"gorm.io/gorm"
)

type AccountsLedger interface {
	WithTx(tx *gorm.DB) AccountsLedger
	Create(ctx context.Context, id domain.AccountID, initial domain.Money) error
	Get(ctx context.Context, id domain.AccountID) (dto.Account, error)
	GetForUpdate(ctx context.Context, id domain.AccountID) (dto.Account, error)
	UpdateBalance(ctx context.Context, id domain.AccountID, newBalance domain.Money) error
}

type TransactionsLedger interface {
	WithTx(tx *gorm.DB) TransactionsLedger
	Insert(ctx context.Context, t dto.Transaction) error
}
