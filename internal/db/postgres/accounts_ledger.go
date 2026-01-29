package postgres

import (
	"context"
	"errors"

	"github.com/payments-core/internal/db/dal"
	"github.com/payments-core/internal/db/dto"
	"github.com/payments-core/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Accounts struct{
	db *gorm.DB
	tx *gorm.DB
}

func NewAccountsLedger(db *gorm.DB) *Accounts { return &Accounts{db: db} }

func (r *Accounts) WithTx(tx *gorm.DB) dal.AccountsLedger { return &Accounts{db: r.db, tx: tx} }

func (r *Accounts) Create(ctx context.Context, id domain.AccountID, initial domain.Money) error {
	a := dto.Account{AccountID: int64(id), Balance: initial.String()}
	if err := r.db.WithContext(ctx).Create(&a).Error; err != nil {
		if isUniqueViolation(err) {
			return domain.ErrAccountExists
		}
		return err
	}
	return nil
}

func (r *Accounts) Get(ctx context.Context, id domain.AccountID) (dto.Account, error) {
	var a dto.Account
	if err := r.db.WithContext(ctx).Where("account_id = ?", int64(id)).Take(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.Account{}, domain.ErrAccountNotFound
		}
		return dto.Account{}, err
	}

	return a, nil
}

func (r *Accounts) GetForUpdate(ctx context.Context, id domain.AccountID) (dto.Account, error) {
	var a dto.Account
	d := r.tx
	if d == nil { d = r.db }
	if err := d.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Where("account_id = ?", int64(id)).Take(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.Account{}, domain.ErrAccountNotFound
		}
		return dto.Account{}, err
	}

	return a, nil
}

func (r *Accounts) UpdateBalance(ctx context.Context, id domain.AccountID, newBalance domain.Money) error {
	d := r.tx
	if d == nil { d = r.db }
	res := d.WithContext(ctx).Model(&dto.Account{}).Where("account_id = ?", int64(id)).Update("balance", newBalance.String())
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
