package service

import (
	"context"

	"github.com/hesampakdaman/banking-service/internal/domain"
)

func (s *BankService) GetAccount(ctx context.Context, accountID string) (domain.Account, error) {
	logger := s.logger.With("account_id", accountID)

	logger.InfoContext(ctx, "Retrieving account")

	account, err := s.repo.GetAccount(ctx, accountID)
	if err != nil {
		logger.WarnContext(ctx, "Failed to retrieve account", "reason", err.Error())
		return domain.Account{}, err
	}

	logger.InfoContext(ctx, "Successfully retrieved account")
	return account, nil
}

func (s *BankService) ListAccounts(ctx context.Context) []domain.Account {
	s.logger.InfoContext(ctx, "Listing all accounts")

	accounts := s.repo.ListAccounts(ctx)

	s.logger.InfoContext(ctx, "Successfully listed accounts", "count", len(accounts))
	return accounts
}

func (s *BankService) ListTransactions(ctx context.Context, accountID string) []domain.Transaction {
	logger := s.logger.With("account_id", accountID)
	logger.InfoContext(ctx, "Listing all transactions for account")

	transactions := s.repo.ListTransactions(ctx, accountID)

	logger.InfoContext(ctx, "Successfully listed transactions for account", "count", len(transactions))
	return transactions
}
