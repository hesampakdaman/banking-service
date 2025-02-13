package integrationtest

import (
	"net/http"
	"testing"
	"time"

	"github.com/hesampakdaman/banking-service/internal/domain"
	"gotest.tools/assert"
)

func TestCreateTransaction_Withdrawal(t *testing.T) {
	server := setupTestServer(t)

	// Given: An account with sufficient balance
	createResp := postJSON(t, server.URL+"/accounts", map[string]interface{}{
		"owner":           "Alice",
		"initial_balance": 1000,
	})

	var accountData map[string]string
	parseJSON(t, createResp, &accountData)
	accountID := accountData["account_id"]

	// When: A valid withdrawal transaction is created
	resp := postJSON(t, server.URL+"/accounts/"+accountID+"/transactions", map[string]interface{}{
		"type":   "withdrawal",
		"amount": 200.0,
	})

	// Then: The response should contain the transaction ID
	assert.Equal(t, resp.StatusCode, http.StatusCreated)

	var txnResp map[string]string
	parseJSON(t, resp, &txnResp)

	txnID, exists := txnResp["transaction_id"]
	assert.Assert(t, exists, "transaction_id should exist in response")
	assert.Assert(t, txnID != "", "transaction_id should not be empty")
}

func TestCreateTransaction_Deposit(t *testing.T) {
	server := setupTestServer(t)

	// Given: An account with an initial balance
	createResp := postJSON(t, server.URL+"/accounts", map[string]interface{}{
		"owner":           "Alice",
		"initial_balance": 1000,
	})

	var accountData map[string]string
	parseJSON(t, createResp, &accountData)
	accountID := accountData["account_id"]

	// When: A valid deposit transaction is created
	resp := postJSON(t, server.URL+"/accounts/"+accountID+"/transactions", map[string]interface{}{
		"type":   "deposit",
		"amount": 500.0,
	})

	// Then: The response should contain the transaction ID
	assert.Equal(t, resp.StatusCode, http.StatusCreated)

	var txnResp map[string]string
	parseJSON(t, resp, &txnResp)

	txnID, exists := txnResp["transaction_id"]
	assert.Assert(t, exists, "transaction_id should exist in response")
	assert.Assert(t, txnID != "", "transaction_id should not be empty")
}

func TestCreateTransaction_Errors(t *testing.T) {
	server := setupTestServer(t)

	// Given: a single test account
	resp := postJSON(t, server.URL+"/accounts", map[string]interface{}{
		"owner":           "Alice",
		"initial_balance": 1000,
	})
	var createResp map[string]string
	parseJSON(t, resp, &createResp)
	testAccountID := createResp["account_id"]

	tt := []struct {
		name       string
		accountID  string // Empty means use testAccountID
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name:       "Invalid Transaction Type",
			payload:    map[string]interface{}{"type": "invalid_type", "amount": 100},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Non-Existent Account",
			accountID:  "non-existent-id",
			payload:    map[string]interface{}{"type": "deposit", "amount": 100},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Insufficient Funds",
			payload:    map[string]interface{}{"type": "withdrawal", "amount": 5000},
			wantStatus: http.StatusConflict,
		},
		{
			name:       "Negative Amount",
			payload:    map[string]interface{}{"type": "deposit", "amount": -100},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Zero Amount",
			payload:    map[string]interface{}{"type": "withdrawal", "amount": 0},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			accountID := tc.accountID
			if accountID == "" {
				accountID = testAccountID
			}

			resp := postJSON(t, server.URL+"/accounts/"+accountID+"/transactions", tc.payload)
			assert.Equal(t, resp.StatusCode, tc.wantStatus)
		})
	}
}

func TestTransfer_Success(t *testing.T) {
	server := setupTestServer(t)

	// Given: Two accounts
	fromResp := postJSON(t, server.URL+"/accounts", map[string]interface{}{
		"owner":           "Alice",
		"initial_balance": 1000,
	})
	toResp := postJSON(t, server.URL+"/accounts", map[string]interface{}{
		"owner":           "Bob",
		"initial_balance": 500,
	})

	var fromAccount, toAccount map[string]string
	parseJSON(t, fromResp, &fromAccount)
	parseJSON(t, toResp, &toAccount)

	fromID := fromAccount["account_id"]
	toID := toAccount["account_id"]

	// When: Transferring 200 from Alice → Bob
	transferResp := postJSON(t, server.URL+"/transfer", map[string]interface{}{
		"from_account_id": fromID,
		"to_account_id":   toID,
		"amount":          200,
	})

	// Then: Verify the response contains transaction IDs
	assert.Equal(t, transferResp.StatusCode, http.StatusCreated)

	var txnResp map[string]string
	parseJSON(t, transferResp, &txnResp)

	assert.Assert(t, txnResp["withdrawal_transaction_id"] != "", "Missing withdrawal transaction ID")
	assert.Assert(t, txnResp["deposit_transaction_id"] != "", "Missing deposit transaction ID")
}

func TestTransfer_Errors(t *testing.T) {
	server := setupTestServer(t)

	// Given: A valid accounts
	validResp := postJSON(t, server.URL+"/accounts", map[string]interface{}{
		"owner":           "Alice",
		"initial_balance": 1000,
	})
	var validAccount map[string]string
	parseJSON(t, validResp, &validAccount)
	accID1 := validAccount["account_id"]

	secondResp := postJSON(t, server.URL+"/accounts", map[string]interface{}{
		"owner":           "Bob",
		"initial_balance": 500,
	})
	var secondAccount map[string]string
	parseJSON(t, secondResp, &secondAccount)
	accID2 := secondAccount["account_id"]

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "Non-existent source account",
			payload: map[string]interface{}{
				"from_account_id": "invalid-id",
				"to_account_id":   accID1,
				"amount":          100,
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "Non-existent destination account",
			payload: map[string]interface{}{
				"from_account_id": accID1,
				"to_account_id":   "invalid-id",
				"amount":          100,
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "Same account transfer",
			payload: map[string]interface{}{
				"from_account_id": accID1,
				"to_account_id":   accID1,
				"amount":          100,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Insufficient funds",
			payload: map[string]interface{}{
				"from_account_id": accID1,
				"to_account_id":   accID2,
				"amount":          5000, // More than balance
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "Invalid amount (negative)",
			payload: map[string]interface{}{
				"from_account_id": accID1,
				"to_account_id":   accID2,
				"amount":          -50,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Zero transfer amount",
			payload: map[string]interface{}{
				"from_account_id": accID1,
				"to_account_id":   accID2,
				"amount":          0,
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := postJSON(t, server.URL+"/transfer", tt.payload)
			assert.Equal(t, resp.StatusCode, tt.wantStatus)
		})
	}
}

func TestListTransactions(t *testing.T) {
	server := setupTestServer(t)

	// Given: A single test account
	resp := postJSON(t, server.URL+"/accounts", map[string]interface{}{
		"owner":           "Alice",
		"initial_balance": 1000,
	})
	var createResp map[string]string
	parseJSON(t, resp, &createResp)
	accountID := createResp["account_id"]

	// And: A single deposit transaction
	resp = postJSON(t, server.URL+"/accounts/"+accountID+"/transactions", map[string]interface{}{
		"type":   "deposit",
		"amount": 500,
	})
	var txnResp map[string]string
	parseJSON(t, resp, &txnResp)
	transactionID := txnResp["transaction_id"]

	// When: Listing transactions for the account
	resp = getJSON(t, server.URL+"/accounts/"+accountID+"/transactions")

	// Then: The response should contain the deposit transaction
	var transactions []domain.Transaction
	parseJSON(t, resp, &transactions)

	expected := []domain.Transaction{
		{
			ID:        transactionID,
			AccountID: accountID,
			Type:      domain.Deposit,
			Amount:    500.0,
		},
	}
	transactions[0].Timestamp = time.Time{} // ignore time field
	assert.DeepEqual(t, expected, transactions)
}
