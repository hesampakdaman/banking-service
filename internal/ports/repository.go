package ports

import (
	"context"

	"github.com/hesampakdaman/banking-service/internal/domain"
)

// Repository defines storage operations for accounts and transactions.
type Repository interface {
	// Account-related operations
	CreateAccount(ctx context.Context, account domain.Account) error
	GetAccount(ctx context.Context, accountID string) (domain.Account, error)
	ListAccounts(ctx context.Context) []domain.Account

	// Transaction-related operations
	Record(ctx context.Context, account domain.Account) error
	ListTransactions(ctx context.Context, accountID string) []domain.Transaction
}
