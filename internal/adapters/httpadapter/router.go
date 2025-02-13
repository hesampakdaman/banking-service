package httpadapter

import (
	"net/http"

	"github.com/hesampakdaman/banking-service/internal/adapters/httpadapter/handlers"
	"github.com/hesampakdaman/banking-service/internal/service"
)

func NewRouter(bankService *service.BankService) *http.ServeMux {
	handler := handlers.NewHTTPHandler(bankService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /accounts", handler.CreateAccountHandler)
	mux.HandleFunc("GET /accounts/{id}", handler.GetAccountHandler)
	mux.HandleFunc("GET /accounts", handler.ListAccountsHandler)
	mux.HandleFunc("POST /accounts/{id}/transactions", handler.CreateTransactionHandler)
	mux.HandleFunc("GET /accounts/{id}/transactions", handler.ListTransactionsHandler)
	mux.HandleFunc("POST /transfer", handler.TransferHandler)

	return mux
}
