package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hesampakdaman/banking-service/internal/domain"
)

func (h *httpHandler) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	accountID := r.PathValue("id")

	var req struct {
		Type   string  `json:"type"` // "deposit" or "withdrawal"
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var txnType domain.TransactionType
	switch req.Type {
	case "deposit":
		txnType = domain.Deposit
	case "withdrawal":
		txnType = domain.Withdrawal
	default:
		http.Error(w, "Invalid transaction type (must be 'deposit' or 'withdrawal')", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.CreateTransaction(r.Context(), accountID, txnType, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), domainErrToStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"transaction_id": transaction.ID}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *httpHandler) TransferHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromAccountID string  `json:"from_account_id"`
		ToAccountID   string  `json:"to_account_id"`
		Amount        float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fromTxn, toTxn, err := h.service.Transfer(r.Context(), req.FromAccountID, req.ToAccountID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), domainErrToStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"withdrawal_transaction_id": fromTxn.ID,
		"deposit_transaction_id":    toTxn.ID,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *httpHandler) ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	accountID := r.PathValue("id")

	transactions := h.service.ListTransactions(r.Context(), accountID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
