package app

import (
	"context"
	"fmt"

	"github.com/payments-core/internal/domain"
	"github.com/payments-core/internal/repo"
	"gorm.io/gorm"
)

type TransferService struct {
	DB           *gorm.DB
	Accounts     repo.AccountRepository
	Transactions repo.TransactionRepository
	DBDriver     string
}

func NewTransferService(db *gorm.DB, driver string, ar repo.AccountRepository, tr repo.TransactionRepository) *TransferService {
	return &TransferService{
		DB:           db,
		DBDriver:     driver,
		Accounts:     ar,
		Transactions: tr,
	}
}

func (s *TransferService) CreateAccount(ctx context.Context, id domain.AccountID, initial domain.Money) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}
	if initial.Decimal.IsNegative() {
		return domain.ErrInvalidInput
	}
	return s.Accounts.Create(ctx, s.DB, id, initial)
}

func (s *TransferService) GetAccount(ctx context.Context, id domain.AccountID) (repo.Account, error) {
	if id <= 0 {
		return repo.Account{}, domain.ErrInvalidInput
	}
	return s.Accounts.Get(ctx, s.DB, id)
}

func (s *TransferService) Transfer(ctx context.Context, source, dest domain.AccountID, amount domain.Money, txID string) error {
	if source <= 0 || dest <= 0 || source == dest || amount.IsNegativeOrZero() {
		return domain.ErrInvalidInput
	}

	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		aID, bID := source, dest
		if aID > bID {
			aID, bID = bID, aID
		}

		if _, err := s.Accounts.GetForUpdate(ctx, tx, aID); err != nil {
			return err
		}
		if _, err := s.Accounts.GetForUpdate(ctx, tx, bID); err != nil {
			return err
		}

		src, err := s.Accounts.GetForUpdate(ctx, tx, source)
		if err != nil {
			return err
		}
		dst, err := s.Accounts.GetForUpdate(ctx, tx, dest)
		if err != nil {
			return err
		}

		if src.Balance.Cmp(amount.Decimal) < 0 {
			return domain.ErrInsufficientFunds
		}

		newSrc := domain.Money{Decimal: src.Balance.Sub(amount.Decimal)}
		newDst := domain.Money{Decimal: dst.Balance.Add(amount.Decimal)}

		if err := s.Accounts.UpdateBalance(ctx, tx, source, newSrc); err != nil {
			return fmt.Errorf("update source: %w", err)
		}
		if err := s.Accounts.UpdateBalance(ctx, tx, dest, newDst); err != nil {
			return fmt.Errorf("update dest: %w", err)
		}

		if err := s.Transactions.Insert(ctx, tx, repo.Transaction{
			ID:                   txID,
			SourceAccountID:      source,
			DestinationAccountID: dest,
			Amount:               amount,
		}); err != nil {
			return fmt.Errorf("insert transaction: %w", err)
		}

		return nil
	})
}
