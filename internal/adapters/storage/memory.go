package storage

import (
	"context"
	"slices"
	"sync"

	"github.com/hesampakdaman/banking-service/internal/domain"
	"github.com/hesampakdaman/banking-service/internal/ports"
)

// MemoryRepository provides an in-memory implementation of Repository.
type MemoryRepository struct {
	mu           sync.RWMutex
	accounts     map[string]domain.Account
	transactions map[string][]domain.Transaction
}

func NewMemoryRepository() ports.Repository {
	return &MemoryRepository{
		accounts:     make(map[string]domain.Account),
		transactions: make(map[string][]domain.Transaction),
	}
}

func (r *MemoryRepository) CreateAccount(ctx context.Context, account domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.accounts[account.ID]; exists {
		return domain.ErrAccountAlreadyExists
	}

	r.accounts[account.ID] = account
	return nil
}

func (r *MemoryRepository) GetAccount(ctx context.Context, accountID string) (domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, exists := r.accounts[accountID]
	if !exists {
		return domain.Account{}, domain.ErrInvalidAccountID
	}

	return account, nil
}

func (r *MemoryRepository) ListAccounts(ctx context.Context) []domain.Account {
	r.mu.RLock()
	defer r.mu.RUnlock()

	accounts := make([]domain.Account, 0, len(r.accounts))
	for _, account := range r.accounts {
		accounts = append(accounts, account)
	}

	return accounts
}

func (r *MemoryRepository) Record(ctx context.Context, account domain.Account, txn domain.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.accounts[account.ID]
	if !exists {
		return domain.ErrInvalidAccountID
	}

	if account.ID != txn.AccountID {
		return domain.ErrAccountTransactionMismatch
	}

	r.accounts[account.ID] = account
	r.transactions[txn.AccountID] = append(r.transactions[txn.AccountID], txn)

	return nil
}

func (r *MemoryRepository) ListTransactions(ctx context.Context, accountID string) []domain.Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transactions, exists := r.transactions[accountID]
	if !exists {
		return nil
	}

	sortedTransactions := make([]domain.Transaction, len(transactions))
	copy(sortedTransactions, transactions)
	slices.SortFunc(sortedTransactions, func(a, b domain.Transaction) int {
		switch {
		case a.Timestamp.Before(b.Timestamp):
			return -1
		case a.Timestamp.After(b.Timestamp):
			return 1
		default:
			return 0
		}
	})
	return sortedTransactions
}
