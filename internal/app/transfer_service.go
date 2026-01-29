package app

import (
	"context"
	"fmt"

	"github.com/payments-core/internal/db/dal"
	"github.com/payments-core/internal/db/dto"
	"github.com/payments-core/internal/domain"
	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

type TransferService struct {
	DB                *gorm.DB
	Accounts          dal.AccountsLedger
	TransactionLedger dal.TransactionsLedger
	DBDriver          string
}

func NewTransferService(db *gorm.DB, dbDriver string, ar dal.AccountsLedger, tr dal.TransactionsLedger) *TransferService {
	return &TransferService{
		DB:                db,
		Accounts:          ar,
		TransactionLedger: tr,
		DBDriver:          dbDriver,
	}
}

func (s *TransferService) CreateAccount(ctx context.Context, id domain.AccountID, initial domain.Money) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}
	if initial.Decimal.IsNegative() {
		return domain.ErrInvalidInput
	}
	return s.Accounts.Create(ctx, id, initial)
}

func (s *TransferService) GetAccount(ctx context.Context, id domain.AccountID) (dto.Account, error) {
	if id <= 0 {
		return dto.Account{}, domain.ErrInvalidInput
	}
	return s.Accounts.Get(ctx, id)
}

func (s *TransferService) Transfer(ctx context.Context, source, dest domain.AccountID, amount domain.Money, txID string) error {
	if source <= 0 || dest <= 0 || source == dest || amount.IsNegativeOrZero() {
		return domain.ErrInvalidInput
	}

	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		
		aID, bID := source, dest
		// this is ensure we do not have deadlocks.
		// for ex : a->b :50$
		// there's another transfer 
		// b-> a : 25$ 
		// if the locks are not obtained in the same order ,it will result in deadlock
		if aID > bID {
			aID, bID = bID, aID
		}

		if _, err := s.Accounts.WithTx(tx).GetForUpdate(ctx, aID); err != nil {
			return err
		}
		if _, err := s.Accounts.WithTx(tx).GetForUpdate(ctx, bID); err != nil {
			return err
		}

		src, err := s.Accounts.WithTx(tx).GetForUpdate(ctx, source)
		if err != nil {
			return err
		}
		dst, err := s.Accounts.WithTx(tx).GetForUpdate(ctx, dest)
		if err != nil {
			return err
		}
		srcBal, _ := decimal.NewFromString(src.Balance)
		dstBal, _ := decimal.NewFromString(dst.Balance)
		if srcBal.Cmp(amount.Decimal) < 0 {
			return domain.ErrInsufficientFunds
		}

		newSrc := domain.Money{Decimal: srcBal.Sub(amount.Decimal)}
		newDst := domain.Money{Decimal: dstBal.Add(amount.Decimal)}

		if err := s.Accounts.WithTx(tx).UpdateBalance(ctx, source, newSrc); err != nil {
			return fmt.Errorf("update source: %w", err)
		}
		if err := s.Accounts.WithTx(tx).UpdateBalance(ctx, dest, newDst); err != nil {
			return fmt.Errorf("update dest: %w", err)
		}

		if err := s.TransactionLedger.WithTx(tx).Insert(ctx, dto.Transaction{
			ID:                   txID,
			SourceAccountID:      int64(source),
			DestinationAccountID: int64(dest),
			Amount:               amount.String(),
		}); err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		return nil
	})
}
