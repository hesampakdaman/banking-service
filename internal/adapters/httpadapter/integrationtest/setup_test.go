package integrationtest

import (
	"io"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/hesampakdaman/banking-service/internal/adapters/httpadapter"
	"github.com/hesampakdaman/banking-service/internal/adapters/storage"
	"github.com/hesampakdaman/banking-service/internal/service"
)

func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	repo := storage.NewMemoryRepository()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	bankService := service.NewBankService(repo, logger)
	router := httpadapter.NewRouter(bankService)

	testServer := httptest.NewServer(router)

	t.Cleanup(func() {
		testServer.Close()
	})

	return testServer
}
