package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/payments-core/internal/domain"
	"github.com/payments-core/internal/repo"
	"gorm.io/gorm"
)

type AccountRepo struct{}

func NewAccountRepo() *AccountRepo { return &AccountRepo{} }

func (r *AccountRepo) Create(ctx context.Context, db *gorm.DB, id domain.AccountID, initial domain.Money) error {
	q := `INSERT INTO accounts(account_id, balance) VALUES (?, ?)`
	if err := db.WithContext(ctx).Exec(q, int64(id), initial.String()).Error; err != nil {
		if isUniqueViolation(err) {
			return domain.ErrAccountExists
		}
		return err
	}
	return nil
}

func (r *AccountRepo) Get(ctx context.Context, db *gorm.DB, id domain.AccountID) (repo.Account, error) {
	var balStr string
	row := db.WithContext(ctx).Raw(`SELECT balance FROM accounts WHERE account_id=?`, int64(id)).Row()
	err := row.Scan(&balStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repo.Account{}, domain.ErrAccountNotFound
		}
		return repo.Account{}, err
	}
	bal, err := domain.ParseMoney(balStr)
	if err != nil {
		return repo.Account{}, fmt.Errorf("parse balance: %w", err)
	}
	return repo.Account{AccountID: id, Balance: bal}, nil
}

func (r *AccountRepo) GetForUpdate(ctx context.Context, tx *gorm.DB, id domain.AccountID) (repo.Account, error) {
	var balStr string
	row := tx.WithContext(ctx).Raw(`SELECT balance FROM accounts WHERE account_id=? FOR UPDATE`, int64(id)).Row()
	err := row.Scan(&balStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repo.Account{}, domain.ErrAccountNotFound
		}
		return repo.Account{}, err
	}
	bal, err := domain.ParseMoney(balStr)
	if err != nil {
		return repo.Account{}, fmt.Errorf("parse balance: %w", err)
	}
	return repo.Account{AccountID: id, Balance: bal}, nil
}

func (r *AccountRepo) UpdateBalance(ctx context.Context, tx *gorm.DB, id domain.AccountID, newBalance domain.Money) error {
	res := tx.WithContext(ctx).Exec(
		`UPDATE accounts SET balance=? WHERE account_id=?`,
		newBalance.String(), int64(id),
	)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrAccountNotFound
	}
	return nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return contains(msg, "duplicate key value") || contains(msg, "unique constraint")
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool {
		return (len(sub) == 0) || (stringIndex(s, sub) >= 0)
	})()
}

func stringIndex(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
