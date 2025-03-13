package domain

// Account represents a bank account entity.
type Account struct {
	ID              string  `json:"id"`
	Owner           string  `json:"owner"`
	Balance         float64 `json:"balance"`
	NewTransactions []Transaction
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

	txn, err := a.recordTransaction(Deposit, amount)
	if err != nil {
		return txn, err
	}

	a.Balance += amount

	return txn, err
}

func (a *Account) Withdraw(amount float64) (Transaction, error) {
	if amount <= 0 {
		return Transaction{}, ErrInvalidAmount
	}
	if amount > a.Balance {
		return Transaction{}, ErrInsufficientFunds
	}

	txn, err := a.recordTransaction(Withdrawal, amount)
	if err != nil {
		return txn, err
	}

	a.Balance -= amount

	return txn, err
}

func (a *Account) Transfer(to *Account, amount float64) (Transaction, Transaction, error) {
	if a.ID == to.ID {
		return Transaction{}, Transaction{}, ErrSelfTransfer
	}

	fromTxn, err := a.Withdraw(amount)
	if err != nil {
		return Transaction{}, Transaction{}, err
	}

	toTxn, err := to.Deposit(amount)
	if err != nil {
		return Transaction{}, Transaction{}, err
	}

	return fromTxn, toTxn, nil
}

func (a *Account) recordTransaction(ttype TransactionType, amount float64) (Transaction, error) {
	txn, err := NewTransaction(a.ID, ttype, amount)
	if err != nil {
		return txn, err
	}

	a.NewTransactions = append(a.NewTransactions, txn)

	return txn, err
}
