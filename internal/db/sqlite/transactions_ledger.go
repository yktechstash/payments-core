package sqlite

import (
	"context"

	"github.com/payments-core/internal/db/dal"
	"github.com/payments-core/internal/db/dto"
	"gorm.io/gorm"
)

type TransactionsLedger struct{
	db *gorm.DB
	tx *gorm.DB
}

func NewTransactionsLedger(db *gorm.DB) *TransactionsLedger { return &TransactionsLedger{db: db} }

func (r *TransactionsLedger) WithTx(tx *gorm.DB) dal.TransactionsLedger { return &TransactionsLedger{db: r.db, tx: tx} }

func (r *TransactionsLedger) Insert(ctx context.Context, t dto.Transaction) error {
	d := r.tx
	if d == nil { d = r.db }
	return d.WithContext(ctx).Create(&t).Error
}
