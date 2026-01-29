package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/payments-core/internal/domain"
	"github.com/payments-core/internal/repo"
	"gorm.io/gorm"
)

type AccountRepo struct{}

func NewAccountRepo() *AccountRepo { return &AccountRepo{} }

func (r *AccountRepo) Create(ctx context.Context, db *gorm.DB, id domain.AccountID, initial domain.Money) error {
	if err := db.WithContext(ctx).Exec(
		`INSERT INTO accounts(account_id, balance) VALUES (?, ?)`,
		int64(id), initial.String(),
	).Error; err != nil {
		if isConstraint(err) {
			return domain.ErrAccountExists
		}
		return err
	}
	return nil
}

func (r *AccountRepo) Get(ctx context.Context, db *gorm.DB, id domain.AccountID) (repo.Account, error) {
	var balStr string
	row := db.WithContext(ctx).Raw(
		`SELECT balance FROM accounts WHERE account_id=?`,
		int64(id),
	).Row()
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
	return r.Get(ctx, tx, id)
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

func isConstraint(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "constraint") || strings.Contains(msg, "unique")
}
