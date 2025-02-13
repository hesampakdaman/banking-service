package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *httpHandler) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Owner          string  `json:"owner"`
		InitialBalance float64 `json:"initial_balance"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	accountID, err := h.service.CreateAccount(r.Context(), req.Owner, req.InitialBalance)
	if err != nil {
		http.Error(w, err.Error(), domainErrToStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]string{"account_id": accountID}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *httpHandler) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	accountID := r.PathValue("id")

	account, err := h.service.GetAccount(r.Context(), accountID)
	if err != nil {
		http.Error(w, err.Error(), domainErrToStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *httpHandler) ListAccountsHandler(w http.ResponseWriter, r *http.Request) {
	accounts := h.service.ListAccounts(r.Context())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
