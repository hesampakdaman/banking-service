package service

import (
	"context"
	"errors"
	"testing"

	"github.com/hesampakdaman/banking-service/internal/adapters/storage"
	"github.com/hesampakdaman/banking-service/internal/domain"
	"gotest.tools/assert"
)

// mockRepository wraps MemoryRepository and allows simulating failures
type mockRepository struct {
	*storage.MemoryRepository
	failOnRecord bool
}

// Record overrides the normal Record function to simulate failure
func (m *mockRepository) Record(ctx context.Context, account domain.Account) error {
	if m.failOnRecord {
		return errors.New("simulated transaction failure")
	}
	return m.MemoryRepository.Record(ctx, account)
}

func TestBankService_Transfer(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: Two accounts exist
	fromID, err := service.CreateAccount(ctx, "Alice", 1000)
	assert.NilError(t, err)

	toID, err := service.CreateAccount(ctx, "Bob", 500)
	assert.NilError(t, err)

	// When: Transferring funds
	fromTxn, toTxn, err := service.Transfer(ctx, fromID, toID, 200)
	assert.NilError(t, err)

	// Then: Transactions should be recorded
	transactions := service.ListTransactions(ctx, fromID)
	assert.Equal(t, len(transactions), 1)
	assert.DeepEqual(t, transactions[0], fromTxn)

	transactions = service.ListTransactions(ctx, toID)
	assert.Equal(t, len(transactions), 1)
	assert.DeepEqual(t, transactions[0], toTxn)

	// And: Account balances should be updated
	fromAccount, err := service.GetAccount(ctx, fromID)
	assert.NilError(t, err)
	assert.Equal(t, fromAccount.Balance, 800.0)

	toAccount, err := service.GetAccount(ctx, toID)
	assert.NilError(t, err)
	assert.Equal(t, toAccount.Balance, 700.0)
}

func TestBankService_Transfer_InsufficientFunds(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: Two accounts exist
	fromID, err := service.CreateAccount(ctx, "Alice", 100)
	assert.NilError(t, err)

	toID, err := service.CreateAccount(ctx, "Bob", 500)
	assert.NilError(t, err)

	// When: Attempting to transfer more than available balance
	_, _, err = service.Transfer(ctx, fromID, toID, 200)

	// Then: Transfer should fail due to insufficient funds
	assert.Assert(t, errors.Is(err, domain.ErrInsufficientFunds))
}

func TestBankService_Transfer_InvalidAccount(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: One valid and one invalid account
	fromID, err := service.CreateAccount(ctx, "Alice", 1000)
	assert.NilError(t, err)

	invalidID := "non-existent-id"

	// When: Transferring to a non-existent account
	_, _, err = service.Transfer(ctx, fromID, invalidID, 100)

	// Then: Transfer should fail
	assert.Assert(t, errors.Is(err, domain.ErrInvalidAccountID))
}

func TestBankService_Transfer_Rollback(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: Two accounts exist
	fromID, err := service.CreateAccount(ctx, "Alice", 1000)
	assert.NilError(t, err)

	toID, err := service.CreateAccount(ctx, "Bob", 500)
	assert.NilError(t, err)

	// Inject failure in repo
	service.repo = &mockRepository{
		MemoryRepository: service.repo.(*storage.MemoryRepository),
		failOnRecord:     true, // Simulate failure on second transaction
	}

	// When: Transferring funds (with failure)
	_, _, err = service.Transfer(ctx, fromID, toID, 200)

	// Then: Transfer should fail and rollback should occur
	assert.ErrorContains(t, err, "simulated transaction failure")

	// And: Source account should have its balance restored
	fromAccount, err := service.GetAccount(ctx, fromID)
	assert.NilError(t, err)
	assert.Equal(t, fromAccount.Balance, 1000.0)

	// And: Destination account should remain unchanged
	toAccount, err := service.GetAccount(ctx, toID)
	assert.NilError(t, err)
	assert.Equal(t, toAccount.Balance, 500.0)
}
