package ports

import (
	"context"

	"github.com/hesampakdaman/banking-service/internal/domain"
)

// BankService defines business operations for accounts and transactions.
type BankService interface {
	CreateAccount(ctx context.Context, owner string, initialBalance float64) (string, error)
	GetAccount(ctx context.Context, accountID string) (domain.Account, error)
	ListAccounts(ctx context.Context) []domain.Account
	CreateTransaction(ctx context.Context, accountID string, txnType domain.TransactionType, amount float64) (domain.Transaction, error)
	ListTransactions(ctx context.Context, accountID string) []domain.Transaction
	Transfer(ctx context.Context, fromAccountID, toAccountID string, amount float64) (domain.Transaction, domain.Transaction, error)
}
