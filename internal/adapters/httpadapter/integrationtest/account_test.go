package integrationtest

import (
	"net/http"
	"testing"

	"github.com/hesampakdaman/banking-service/internal/domain"
	"gotest.tools/assert"
)

func TestGetAccount_NotFound(t *testing.T) {
	server := setupTestServer(t)

	// When: Trying to retrieve a non-existent account
	resp := getJSON(t, server.URL+"/accounts/non-existent-id")

	// Then: The response should indicate not found
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
}

func TestCreateAndGetAccount(t *testing.T) {
	server := setupTestServer(t)

	// Given: A new account request with valid data
	resp := postJSON(t, server.URL+"/accounts", map[string]interface{}{
		"owner":           "Alice",
		"initial_balance": 1000,
	})

	// When: The account is created
	var createResp map[string]string
	parseJSON(t, resp, &createResp)

	accountID, exists := createResp["account_id"]
	assert.Assert(t, exists, "account_id should exist in response")

	// And: We retrieve the created account
	resp = getJSON(t, server.URL+"/accounts/"+accountID)

	// Then: The response should contain the correct actual details
	var actual domain.Account
	parseJSON(t, resp, &actual)

	expected, _ := domain.NewAccount(accountID, "Alice", 1000.0)
	assert.DeepEqual(t, expected, actual)
}

func TestCreateAccount_InvalidInput(t *testing.T) {
	server := setupTestServer(t)

	tests := []struct {
		name       string
		request    map[string]interface{}
		wantStatus int
	}{
		{
			name: "Missing owner",
			request: map[string]interface{}{
				"initial_balance": 1000,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Negative balance",
			request: map[string]interface{}{
				"owner":           "Alice",
				"initial_balance": -500,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Empty body",
			request:    map[string]interface{}{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := postJSON(t, server.URL+"/accounts", tc.request)
			assert.Equal(t, resp.StatusCode, tc.wantStatus)
		})
	}
}

func TestCreateAccount_Duplicate(t *testing.T) {
	server := setupTestServer(t)

	// Override UUID generation to return a fixed ID
	originalUUIDFunc := domain.GetUUID
	defer func() { domain.GetUUID = originalUUIDFunc }()
	domain.GetUUID = func() string { return "fixed-uuid" }

	// Given: A valid account request
	resp := postJSON(t, server.URL+"/accounts", map[string]interface{}{
		"owner":           "Alice",
		"initial_balance": 1000,
	})
	assert.Equal(t, resp.StatusCode, http.StatusCreated)

	// When: Trying to create an account with the same UUID
	resp = postJSON(t, server.URL+"/accounts", map[string]interface{}{
		"owner":           "Bob",
		"initial_balance": 500,
	})

	// Then: The response should indicate conflict
	assert.Equal(t, resp.StatusCode, http.StatusConflict)
}

func TestListAccounts(t *testing.T) {
	server := setupTestServer(t)

	// Given: A single account exists
	postJSON(t, server.URL+"/accounts", map[string]interface{}{
		"owner":           "Alice",
		"initial_balance": 1000,
	})

	// When: We retrieve the list of accounts
	resp := getJSON(t, server.URL+"/accounts")

	// Then: The response should contain exactly one account
	var actual []domain.Account
	parseJSON(t, resp, &actual)

	expected, _ := domain.NewAccount(actual[0].ID, "Alice", 1000.0)
	assert.DeepEqual(t, []domain.Account{expected}, actual)
}
