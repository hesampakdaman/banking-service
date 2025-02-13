package domain

import (
	"time"
)

// TransactionType represents the type of transaction.
type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
)

// Transaction represents a bank transaction entity.
type Transaction struct {
	ID        string          `json:"id"`
	AccountID string          `json:"account_id"`
	Type      TransactionType `json:"type"`
	Amount    float64         `json:"amount"`
	Timestamp time.Time       `json:"timestamp"`
}

func NewTransaction(accountID string, txnType TransactionType, amount float64) (Transaction, error) {
	if accountID == "" {
		return Transaction{}, ErrInvalidAccountID
	}

	if txnType != Deposit && txnType != Withdrawal {
		return Transaction{}, ErrInvalidTransactionType
	}

	if amount <= 0 {
		return Transaction{}, ErrInvalidAmount
	}

	return Transaction{
		ID:        GetUUID(),
		AccountID: accountID,
		Type:      txnType,
		Amount:    amount,
		Timestamp: GetTimeNow(),
	}, nil
}
