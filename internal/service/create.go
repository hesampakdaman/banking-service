package service

import (
	"context"

	"github.com/hesampakdaman/banking-service/internal/domain"
)

func (s *BankService) CreateAccount(ctx context.Context, owner string, initialBalance float64) (string, error) {
	logger := s.logger.With("owner", owner, "balance", initialBalance)

	logger.InfoContext(ctx, "Creating account")

	account, err := domain.NewAccount(domain.GetUUID(), owner, initialBalance)
	if err != nil {
		logger.WarnContext(ctx, "Failed to create account", "reason", err.Error())
		return "", err
	}

	logger = logger.With("account_id", account.ID)
	if err := s.repo.CreateAccount(ctx, account); err != nil {
		logger.ErrorContext(ctx, "Failed to create account", "error", err.Error())
		return "", err
	}

	logger.InfoContext(ctx, "Successfully created account")
	return account.ID, nil
}

func (s *BankService) CreateTransaction(ctx context.Context, accountID string, txnType domain.TransactionType, amount float64) (domain.Transaction, error) {
	logger := s.logger.With("account_id", accountID, "amount", amount, "transaction_type", txnType)

	logger.InfoContext(ctx, "Processing transaction")

	account, err := s.repo.GetAccount(ctx, accountID)
	if err != nil {
		logger.WarnContext(ctx, "Transaction failed (invalid account)", "reason", err.Error())
		return domain.Transaction{}, err
	}

	var transaction domain.Transaction
	if txnType == domain.Deposit {
		transaction, err = account.Deposit(amount)
	} else {
		transaction, err = account.Withdraw(amount)
	}

	if err != nil {
		logger.WarnContext(ctx, "Transaction denied", "reason", err.Error())
		return domain.Transaction{}, err
	}

	if err := s.repo.Record(ctx, account, transaction); err != nil {
		logger.ErrorContext(ctx, "Failed to record transaction", "error", err.Error())
		return domain.Transaction{}, err
	}

	logger.InfoContext(ctx, "Transaction successful")
	return transaction, nil
}
