package sqlite

import (
	"context"

	"github.com/payments-core/internal/repo"
	"gorm.io/gorm"
)

type TransactionRepo struct{}

func NewTransactionRepo() *TransactionRepo { return &TransactionRepo{} }

func (r *TransactionRepo) Insert(ctx context.Context, tx *gorm.DB, t repo.Transaction) error {
	return tx.WithContext(ctx).Exec(
		`INSERT INTO transactions(id, source_account_id, destination_account_id, amount) VALUES (?, ?, ?, ?)`,
		t.ID, int64(t.SourceAccountID), int64(t.DestinationAccountID), t.Amount.String(),
	).Error
}
