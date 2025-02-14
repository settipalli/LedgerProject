package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"ledgerproject/config"
	"ledgerproject/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewServer(t *testing.T) {
	t.Run("successful server creation", func(t *testing.T) {
		mockLedger := new(MockLedger)
		cfg := &config.Config{
			ServerPort:   ":8080",
			CurrencyFile: "test_currency.json",
		}

		server := NewServer(mockLedger, cfg)

		assert.NotNil(t, server)
		assert.NotNil(t, server.router)
		assert.Equal(t, mockLedger, server.ledger)
		assert.Equal(t, cfg, server.config)
	})
}

func TestSetupRoutes(t *testing.T) {
	mockLedger := new(MockLedger)
	cfg := &config.Config{ServerPort: ":8080"}
	server := NewServer(mockLedger, cfg)

	// Helper function to test route existence and method
	testRoute := func(path, method string) {
		route := findRoute(server.router, path, method)
		assert.NotNil(t, route, "Route %s %s should exist", method, path)
	}

	// Test all expected routes
	testRoute("/accounts", "POST")
	testRoute("/transactions", "POST")
	testRoute("/accounts/{accountId}/balance", "GET")
	testRoute("/accounts/{accountId}/history", "GET")
}

func TestRouteHandlers(t *testing.T) {
	mockLedger := new(MockLedger)
	cfg := &config.Config{ServerPort: ":8080"}
	server := NewServer(mockLedger, cfg)

	t.Run("create account route", func(t *testing.T) {
		account := models.Account{
			ID:       "TEST1",
			Name:     "Test Account",
			Currency: "USD",
			Type:     "asset",
			Balance: models.Money{
				Amount:   decimal.NewFromInt(100),
				Currency: "USD",
			},
		}

		mockLedger.On("CreateAccount", mock.MatchedBy(func(a models.Account) bool {
			return a.ID == account.ID
		})).Return(nil)

		body, _ := json.Marshal(account)
		req := httptest.NewRequest("POST", "/accounts", strings.NewReader(string(body)))
		rr := httptest.NewRecorder()

		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockLedger.AssertExpectations(t)
	})

	t.Run("record transaction route", func(t *testing.T) {
		tx := models.Transaction{
			ID:            "TX1",
			DebitAccount:  "ACC1",
			CreditAccount: "ACC2",
			Money: models.Money{
				Amount:   decimal.NewFromInt(100),
				Currency: "USD",
			},
		}

		mockLedger.On("RecordTransaction", mock.MatchedBy(func(t models.Transaction) bool {
			return t.ID == tx.ID
		})).Return(nil)

		body, _ := json.Marshal(tx)
		req := httptest.NewRequest("POST", "/transactions", strings.NewReader(string(body)))
		rr := httptest.NewRecorder()

		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockLedger.AssertExpectations(t)
	})

	t.Run("get balance route", func(t *testing.T) {
		accountID := "ACC1"
		balance := models.Money{
			Amount:   decimal.NewFromInt(100),
			Currency: "USD",
		}

		mockLedger.On("GetAccountBalance", accountID).Return(balance, nil)

		req := httptest.NewRequest("GET", fmt.Sprintf("/accounts/%s/balance", accountID), nil)
		rr := httptest.NewRecorder()

		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockLedger.AssertExpectations(t)
	})

	t.Run("get transaction history route", func(t *testing.T) {
		accountID := "ACC1"
		history := []models.Transaction{
			{
				ID:            "TX1",
				DebitAccount:  accountID,
				CreditAccount: "ACC2",
				Money: models.Money{
					Amount:   decimal.NewFromInt(100),
					Currency: "USD",
				},
			},
		}

		mockLedger.On("GetTransactionHistory", accountID).Return(history)

		req := httptest.NewRequest("GET", fmt.Sprintf("/accounts/%s/history", accountID), nil)
		rr := httptest.NewRecorder()

		server.router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockLedger.AssertExpectations(t)
	})
}

func TestServerStart(t *testing.T) {
	t.Run("server start with valid port", func(t *testing.T) {
		mockLedger := new(MockLedger)
		cfg := &config.Config{ServerPort: ":0"} // Use port 0 for testing (OS will assign free port)
		server := NewServer(mockLedger, cfg)

		// Start server in goroutine
		go func() {
			err := server.Start()
			// Should only return with error
			assert.NotNil(t, err)
		}()

		// Give the server a moment to start
		// Make a test request to ensure it's running
		req := httptest.NewRequest("GET", "/health", nil)
		rr := httptest.NewRecorder()
		server.router.ServeHTTP(rr, req)
	})

	t.Run("server start with invalid port", func(t *testing.T) {
		mockLedger := new(MockLedger)
		cfg := &config.Config{ServerPort: "invalid-port"}
		server := NewServer(mockLedger, cfg)

		err := server.Start()
		assert.Error(t, err)
	})
}

// Helper function to find a route in the router
func findRoute(router *mux.Router, path, method string) *mux.Route {
	var foundRoute *mux.Route
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}

		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}

		if tpl == path && contains(methods, method) {
			foundRoute = route
			return fmt.Errorf("route found") // Use error to stop walking
		}

		return nil
	})
	return foundRoute
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Integration
func TestServerIntegration(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	mockLedger := new(MockLedger)
	cfg := &config.Config{ServerPort: ":8081"}
	server := NewServer(mockLedger, cfg)

	// Start server
	go server.Start()

	// Test complete flow: Create account -> Record transaction -> Get balance -> Get history
	t.Run("complete flow", func(t *testing.T) {
		// Setup mock expectations
		account1 := models.Account{
			ID:       "INT1",
			Name:     "Integration Test Account 1",
			Currency: "USD",
			Type:     "asset",
			Balance: models.Money{
				Amount:   decimal.NewFromInt(100),
				Currency: "USD",
			},
		}

		account2 := models.Account{
			ID:       "INT2",
			Name:     "Integration Test Account 2",
			Currency: "USD",
			Type:     "liability",
			Balance: models.Money{
				Amount:   decimal.NewFromInt(-100),
				Currency: "USD",
			},
		}

		tx := models.Transaction{
			ID:            "TXI1",
			Description:   "Test Transaction",
			DebitAccount:  "INT1",
			CreditAccount: "INT2",
			Money: models.Money{
				Amount:   decimal.NewFromInt(100),
				Currency: "USD",
			},
		}

		balance := models.Money{
			Amount:   decimal.NewFromInt(100),
			Currency: "USD",
		}

		// Setup mock behaviors
		mockLedger.On("CreateAccount", mock.MatchedBy(func(a models.Account) bool {
			return a.ID == account1.ID || a.ID == account2.ID
		})).Return(nil)

		mockLedger.On("RecordTransaction", mock.MatchedBy(func(t models.Transaction) bool {
			return t.ID == tx.ID
		})).Return(nil)

		mockLedger.On("GetAccountBalance", account1.ID).Return(balance, nil)
		mockLedger.On("GetTransactionHistory", account1.ID).Return([]models.Transaction{tx})

		// Create HTTP client
		client := &http.Client{}
		baseURL := fmt.Sprintf("http://localhost%s", cfg.ServerPort)

		// 1. Create first account
		t.Log("Creating first account...")
		accountBody1, err := json.Marshal(account1)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", baseURL+"/accounts", bytes.NewBuffer(accountBody1))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()

		// 2. Create second account
		t.Log("Creating second account...")
		accountBody2, err := json.Marshal(account2)
		require.NoError(t, err)

		req, err = http.NewRequest("POST", baseURL+"/accounts", bytes.NewBuffer(accountBody2))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()

		// 3. Record transaction between accounts
		t.Log("Recording transaction...")
		txBody, err := json.Marshal(tx)
		require.NoError(t, err)

		req, err = http.NewRequest("POST", baseURL+"/transactions", bytes.NewBuffer(txBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()

		// 4. Get balance for first account
		t.Log("Getting account balance...")
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/accounts/%s/balance", baseURL, account1.ID), nil)
		require.NoError(t, err)

		resp, err = client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var balanceResponse map[string]models.Money
		err = json.NewDecoder(resp.Body).Decode(&balanceResponse)
		require.NoError(t, err)
		assert.Equal(t, balance, balanceResponse["balance"])
		resp.Body.Close()

		// 5. Get transaction history for first account
		t.Log("Getting transaction history...")
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/accounts/%s/history", baseURL, account1.ID), nil)
		require.NoError(t, err)

		resp, err = client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var historyResponse []models.Transaction
		err = json.NewDecoder(resp.Body).Decode(&historyResponse)
		require.NoError(t, err)
		assert.Len(t, historyResponse, 1)
		assert.Equal(t, tx.ID, historyResponse[0].ID)
		resp.Body.Close()

		// Verify all mock expectations were met
		mockLedger.AssertExpectations(t)
	})
}
