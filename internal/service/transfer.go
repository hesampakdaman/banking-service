package service

import (
	"context"

	"github.com/hesampakdaman/banking-service/internal/domain"
)

func (s *BankService) Transfer(ctx context.Context, fromAccountID, toAccountID string, amount float64) (domain.Transaction, domain.Transaction, error) {
	logger := s.logger.With("from_account_id", fromAccountID, "to_account_id", toAccountID, "amount", amount)

	logger.InfoContext(ctx, "Processing transfer")

	// Fetch both accounts from repository
	fromAccount, err := s.repo.GetAccount(ctx, fromAccountID)
	if err != nil {
		logger.WarnContext(ctx, "Transfer failed (invalid source account)", "reason", err.Error())
		return domain.Transaction{}, domain.Transaction{}, err
	}

	toAccount, err := s.repo.GetAccount(ctx, toAccountID)
	if err != nil {
		logger.WarnContext(ctx, "Transfer failed (invalid destination account)", "reason", err.Error())
		return domain.Transaction{}, domain.Transaction{}, err
	}

	// Attempt transfer
	fromTxn, toTxn, err := fromAccount.Transfer(&toAccount, amount)
	if err != nil {
		logger.WarnContext(ctx, "Transfer denied", "reason", err.Error())
		return domain.Transaction{}, domain.Transaction{}, err
	}

	// Record both transactions, ensuring consistency
	if err := s.repo.Record(ctx, fromAccount); err != nil {
		logger.ErrorContext(ctx, "Failed to record source transaction", "error", err.Error())
		return domain.Transaction{}, domain.Transaction{}, err
	}

	if err := s.repo.Record(ctx, toAccount); err != nil {
		// **Rollback:** Attempt to revert withdrawal
		logger.ErrorContext(ctx, "Failed to record destination transaction, attempting rollback", "error", err.Error())

		_, rollbackErr := fromAccount.Deposit(amount)

		if rollbackErr != nil {
			logger.ErrorContext(ctx, "Rollback failed, system may be in an inconsistent state", "rollback_error", rollbackErr.Error())
		} else {
			if recErr := s.repo.Record(ctx, fromAccount); recErr != nil {
				logger.ErrorContext(ctx, "Failed to record rollback transaction", "rollback_error", recErr.Error())
			} else {
				logger.WarnContext(ctx, "Rollback successful")
			}
		}
		return domain.Transaction{}, domain.Transaction{}, err
	}

	logger.InfoContext(ctx, "Transfer successful")
	return fromTxn, toTxn, nil
}
