package domain

// Account represents a bank account entity.
type Account struct {
	ID      string  `json:"id"`
	Owner   string  `json:"owner"`
	Balance float64 `json:"balance"`
}

func NewAccount(ID string, owner string, initialBalance float64) (Account, error) {
	if ID == "" {
		return Account{}, ErrInvalidAccountID
	}
	if owner == "" {
		return Account{}, ErrInvalidOwner
	}
	if initialBalance < 0 {
		return Account{}, ErrNegativeBalance
	}

	return Account{
		ID:      ID,
		Owner:   owner,
		Balance: initialBalance,
	}, nil
}

func (a *Account) Deposit(amount float64) (Transaction, error) {
	if amount <= 0 {
		return Transaction{}, ErrInvalidAmount
	}

	a.Balance += amount

	return NewTransaction(a.ID, Deposit, amount)
}

func (a *Account) Withdraw(amount float64) (Transaction, error) {
	if amount <= 0 {
		return Transaction{}, ErrInvalidAmount
	}
	if amount > a.Balance {
		return Transaction{}, ErrInsufficientFunds
	}

	a.Balance -= amount

	return NewTransaction(a.ID, Withdrawal, amount)
}

func (a *Account) Transfer(to *Account, amount float64) (Transaction, Transaction, error) {
	if a.ID == to.ID {
		return Transaction{}, Transaction{}, ErrSelfTransfer
	}
	if amount <= 0 {
		return Transaction{}, Transaction{}, ErrInvalidAmount
	}
	if amount > a.Balance {
		return Transaction{}, Transaction{}, ErrInsufficientFunds
	}

	a.Balance -= amount
	to.Balance += amount

	return Transaction{
			ID:        GetUUID(),
			AccountID: a.ID,
			Type:      Withdrawal,
			Amount:    amount,
			Timestamp: GetTimeNow(),
		}, Transaction{
			ID:        GetUUID(),
			AccountID: to.ID,
			Type:      Deposit,
			Amount:    amount,
			Timestamp: GetTimeNow(),
		}, nil
}
