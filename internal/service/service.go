package service

import (
	"log/slog"

	"github.com/hesampakdaman/banking-service/internal/ports"
)

// BankService provides business logic for accounts and transactions.
type BankService struct {
	repo   ports.Repository
	logger *slog.Logger
}

func NewBankService(repo ports.Repository, logger *slog.Logger) *BankService {
	logger = logger.With("component", "BankService")
	return &BankService{repo: repo, logger: logger}
}
