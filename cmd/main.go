package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/hesampakdaman/banking-service/internal/adapters/httpadapter"
	"github.com/hesampakdaman/banking-service/internal/adapters/storage"
	"github.com/hesampakdaman/banking-service/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(log.Writer(), nil))

	// Initialize repository & service layer
	repo := storage.NewMemoryRepository()
	bankService := service.NewBankService(repo, logger)

	// Initialize http server
	mux := httpadapter.NewRouter(bankService)
	loggedMux := httpadapter.LoggingMiddleware(mux, logger)
	server := &http.Server{
		Addr:    ":8080",
		Handler: loggedMux,
	}

	logger.Info("Starting banking-service on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
