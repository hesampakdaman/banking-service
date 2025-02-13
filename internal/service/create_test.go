package service

import (
	"context"
	"errors"
	"testing"

	"github.com/hesampakdaman/banking-service/internal/domain"
	"gotest.tools/assert"
)

func TestBankService_CreateAccount(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: A valid account request
	accountID, err := service.CreateAccount(ctx, "foo", 1000)
	assert.NilError(t, err)

	// When: We retrieve the created account
	account, err := service.GetAccount(ctx, accountID)
	assert.NilError(t, err)

	// Then: The account should exist with correct balance
	expected, _ := domain.NewAccount(accountID, "foo", 1000)
	assert.DeepEqual(t, expected, account)
}

func TestBankService_CreateAccount_NegativeBalance(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: An attempt to create an account with a negative balance
	_, err := service.CreateAccount(ctx, "foo", -100)

	// Then: It should fail with ErrNegativeBalance
	assert.Assert(t, errors.Is(err, domain.ErrNegativeBalance))
}

func TestBankService_CreateAccount_MissingOwner(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: An attempt to create an account with an empty owner
	_, err := service.CreateAccount(ctx, "", 500)

	// Then: It should fail with ErrInvalidOwner
	assert.Assert(t, errors.Is(err, domain.ErrInvalidOwner))
}

func TestBankService_CreateAccount_Duplicate(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Override GetUUID() for deterministic testing
	originalUUID := domain.GetUUID
	domain.GetUUID = func() string { return "fixed-uuid" }
	defer func() { domain.GetUUID = originalUUID }()

	// Given: A valid account is created
	_, err := service.CreateAccount(ctx, "foo", 1000)
	assert.NilError(t, err)

	// When: Trying to create another account (which will get the same "fixed-uuid")
	_, err = service.CreateAccount(ctx, "bar", 500)

	// Then: It should fail with ErrAccountAlreadyExists
	assert.Assert(t, errors.Is(err, domain.ErrAccountAlreadyExists))
}

func TestBankService_Withdraw(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: An account with sufficient balance
	accountID, err := service.CreateAccount(ctx, "foo", 1000)
	assert.NilError(t, err)

	// When: A valid withdrawal is made
	_, err = service.CreateTransaction(ctx, accountID, domain.Withdrawal, 200)
	assert.NilError(t, err)

	// Then: The balance should be updated
	account, err := service.GetAccount(ctx, accountID)
	assert.NilError(t, err)
	assert.Equal(t, account.Balance, 800.0)
}

func TestBankService_Withdraw_NonExistentAccount(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// When: Trying to withdraw from a non-existent account
	_, err := service.CreateTransaction(ctx, "non-existent-id", domain.Withdrawal, 100)

	// Then: Should return an error
	assert.Assert(t, errors.Is(err, domain.ErrInvalidAccountID))
}

func TestBankService_Withdraw_NegativeAmount(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: An existing account
	accountID, err := service.CreateAccount(ctx, "foo", 1000)
	assert.NilError(t, err)

	// When: Trying to withdraw a negative amount
	_, err = service.CreateTransaction(ctx, accountID, domain.Withdrawal, -100)

	// Then: Should return an error
	assert.Assert(t, errors.Is(err, domain.ErrInvalidAmount))
}

func TestBankService_Withdraw_ZeroAmount(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: An existing account
	accountID, err := service.CreateAccount(ctx, "foo", 1000)
	assert.NilError(t, err)

	// When: Trying to withdraw zero
	_, err = service.CreateTransaction(ctx, accountID, domain.Withdrawal, 0)

	// Then: Should return an error
	assert.Assert(t, errors.Is(err, domain.ErrInvalidAmount))
}

func TestBankService_Withdraw_InsufficientFunds(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: An account with limited funds
	accountID, err := service.CreateAccount(ctx, "foo", 100)
	assert.NilError(t, err)

	// When: Trying to withdraw more than available balance
	_, err = service.CreateTransaction(ctx, accountID, domain.Withdrawal, 500)

	// Then: Should return an error
	assert.Assert(t, errors.Is(err, domain.ErrInsufficientFunds))
}

func TestBankService_Deposit(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: An account with an initial balance
	accountID, err := service.CreateAccount(ctx, "foo", 1000)
	assert.NilError(t, err)

	// When: A valid deposit is made
	_, err = service.CreateTransaction(ctx, accountID, domain.Deposit, 500)
	assert.NilError(t, err)

	// Then: The balance should be updated
	account, err := service.GetAccount(ctx, accountID)
	assert.NilError(t, err)
	assert.Equal(t, account.Balance, 1500.0)
}

func TestBankService_Deposit_NonExistentAccount(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// When: Trying to deposit to a non-existent account
	_, err := service.CreateTransaction(ctx, "non-existent-id", domain.Deposit, 100)

	// Then: Should return an error
	assert.Assert(t, errors.Is(err, domain.ErrInvalidAccountID))
}

func TestBankService_Deposit_NegativeAmount(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: An existing account
	accountID, err := service.CreateAccount(ctx, "foo", 1000)
	assert.NilError(t, err)

	// When: Trying to deposit a negative amount
	_, err = service.CreateTransaction(ctx, accountID, domain.Deposit, -100)

	// Then: Should return an error
	assert.Assert(t, errors.Is(err, domain.ErrInvalidAmount))
}

func TestBankService_Deposit_ZeroAmount(t *testing.T) {
	service := fixture()
	ctx := context.Background()

	// Given: An existing account
	accountID, err := service.CreateAccount(ctx, "foo", 1000)
	assert.NilError(t, err)

	// When: Trying to deposit zero
	_, err = service.CreateTransaction(ctx, accountID, domain.Deposit, 0)

	// Then: Should return an error
	assert.Assert(t, errors.Is(err, domain.ErrInvalidAmount))
}
