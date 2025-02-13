package service

import (
	"context"
	"slices"
	"testing"

	"github.com/hesampakdaman/banking-service/internal/domain"
	"gotest.tools/assert"
)

func TestBankService_ListAccounts(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: Multiple accounts exist
	account1, err := service.CreateAccount(ctx, "Foo", 1000)
	assert.NilError(t, err)

	account2, err := service.CreateAccount(ctx, "Bar", 500)
	assert.NilError(t, err)

	// When: Listing accounts
	accounts := service.ListAccounts(ctx)

	// Then: Both accounts should be returned
	assert.Equal(t, len(accounts), 2)

	// And: Accounts should be correct
	expectedFoo, _ := domain.NewAccount(account1, "Foo", 1000)
	expectedBar, _ := domain.NewAccount(account2, "Bar", 500)

	assert.Assert(t, slices.Contains(accounts, expectedFoo))
	assert.Assert(t, slices.Contains(accounts, expectedBar))
}

func TestBankService_ListTransactions(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: An account with deposits and withdrawals
	accountID, err := service.CreateAccount(ctx, "foo", 1000)
	assert.NilError(t, err)

	depositTxn, _ := service.CreateTransaction(ctx, accountID, domain.Deposit, 200)
	withdrawTxn, _ := service.CreateTransaction(ctx, accountID, domain.Withdrawal, 100)

	// When: Listing transactions
	transactions := service.ListTransactions(ctx, accountID)

	// Then: The transactions should be recorded correctly
	assert.DeepEqual(t, transactions, []domain.Transaction{depositTxn, withdrawTxn})
}
