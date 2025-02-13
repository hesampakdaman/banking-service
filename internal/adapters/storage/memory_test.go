package storage

import (
	"context"
	"errors"
	"slices"
	"testing"

	"gotest.tools/assert"

	"github.com/hesampakdaman/banking-service/internal/domain"
)

func TestMemoryRepository_CreateAndGetAccount(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	// Given: A new account
	expected, _ := domain.NewAccount(domain.GetUUID(), "foo", 100.0)

	// When: The account is created
	_ = repo.CreateAccount(ctx, expected)

	// Then: It should be retrievable
	actual, err := repo.GetAccount(ctx, expected.ID)
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestMemoryRepository_CannotCreateDuplicateAccount(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	// Given: An account already exists
	account, _ := domain.NewAccount("123", "foo", 50.0)
	_ = repo.CreateAccount(ctx, account)

	// When: Trying to create an account with the same ID
	account, _ = domain.NewAccount("123", "bar", 0)
	err := repo.CreateAccount(ctx, account)

	// Then: It should return an error indicating account already exists
	assert.Assert(t, errors.Is(err, domain.ErrAccountAlreadyExists))
}

func TestMemoryRepository_GetNonExistentAccount(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	// Given: No accounts exist
	// When: We try to get a non-existent account
	_, err := repo.GetAccount(ctx, "non-existent-id")

	// Then: It should return an error indicating account not found
	assert.Assert(t, errors.Is(err, domain.ErrInvalidAccountID))
}

func TestMemoryRepository_ListAccounts(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	// Given: Multiple accounts exist
	account1, _ := domain.NewAccount(domain.GetUUID(), "foo", 100.0)
	account2, _ := domain.NewAccount(domain.GetUUID(), "bar", 200.0)
	_ = repo.CreateAccount(ctx, account1)
	_ = repo.CreateAccount(ctx, account2)

	// When: Listing accounts
	accounts := repo.ListAccounts(ctx)

	// Then: All accounts should be returned
	assert.Equal(t, len(accounts), 2)
	assert.Assert(t, slices.Contains(accounts, account1))
	assert.Assert(t, slices.Contains(accounts, account2))

}

func TestMemoryRepository_Record(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	// Given: An existing account
	account, _ := domain.NewAccount(domain.GetUUID(), "foo", 100.0)
	_ = repo.CreateAccount(ctx, account)

	// And: A deposit transaction
	transaction, _ := account.Deposit(50.0)

	// When: The transaction is recorded
	_ = repo.Record(ctx, account, transaction)

	// Then: It should appear in the list of transactions
	transactions := repo.ListTransactions(ctx, account.ID)
	assert.DeepEqual(t, []domain.Transaction{transaction}, transactions)

	// And: The account balance should be updated correctly
	updatedAccount, err := repo.GetAccount(ctx, account.ID)
	assert.NilError(t, err)
	assert.Equal(t, updatedAccount.Balance, 150.0)
}

func TestMemoryRepository_ListTransactionsForNonExistentAccount(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	// Given: No transactions exist
	// When: Listing transactions for a non-existent account
	transactions := repo.ListTransactions(ctx, "non-existent-id")

	// Then: It should return an empty slice with no error
	assert.Equal(t, len(transactions), 0)
}
