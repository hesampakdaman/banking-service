package domain

import (
	"errors"
)

var (
	ErrAccountAlreadyExists       = errors.New("account already exists")
	ErrAccountTransactionMismatch = errors.New("account and transaction mismatch")
	ErrInsufficientFunds          = errors.New("insufficient funds")
	ErrInvalidAccountID           = errors.New("invalid account")
	ErrInvalidAmount              = errors.New("transaction amount must be positive")
	ErrInvalidOwner               = errors.New("owner name cannot be empty")
	ErrInvalidTransactionType     = errors.New("invalid transaction type")
	ErrNegativeBalance            = errors.New("initial balance cannot be negative")
	ErrSelfTransfer               = errors.New("cannot transfer funds to the same account")
)
