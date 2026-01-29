package postgres

import (
	"context"

	"github.com/payments-core/internal/db/dal"
	repo "github.com/payments-core/internal/db/dto"
	"gorm.io/gorm"
)

type TransactionsLedger struct{
	db *gorm.DB
	tx *gorm.DB
}

func NewTransactionsLedger(db *gorm.DB) *TransactionsLedger { return &TransactionsLedger{db: db} }

func (r *TransactionsLedger) WithTx(tx *gorm.DB) dal.TransactionsLedger { return &TransactionsLedger{db: r.db, tx: tx} }

func (r *TransactionsLedger) Insert(ctx context.Context, t repo.Transaction) error {
	d := r.tx
	if d == nil { d = r.db }
	dto := repo.Transaction{
		ID:                   t.ID,
		SourceAccountID:      int64(t.SourceAccountID),
		DestinationAccountID: int64(t.DestinationAccountID),
		Amount:               t.Amount,
	}
	return d.WithContext(ctx).Create(&dto).Error
}
