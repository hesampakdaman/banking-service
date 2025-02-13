package service

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/hesampakdaman/banking-service/internal/adapters/storage"
	"github.com/hesampakdaman/banking-service/internal/domain"
)

// fixture initializes a BankService with a test logger.
func fixture() *BankService {
	repo := storage.NewMemoryRepository()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	return NewBankService(repo, logger)
}

func TestMain(m *testing.M) {
	// Set a fixed time for all tests
	domain.GetTimeNow = func() time.Time {
		return time.Date(2025, 2, 12, 12, 0, 0, 0, time.UTC)
	}

	os.Exit(m.Run())
}
