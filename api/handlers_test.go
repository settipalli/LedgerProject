package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"ledgerproject/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Test helpers
func setupTest() (*Server, *MockLedger) {
	mockLedger := new(MockLedger)
	server := &Server{
		router: mux.NewRouter(),
		ledger: mockLedger,
	}
	return server, mockLedger
}

// CreateAccountHandler tests
func TestCreateAccountHandler(t *testing.T) {
	t.Run("successful account creation", func(t *testing.T) {
		server, mockLedger := setupTest()

		account := models.Account{
			ID:       "ACC123",
			Name:     "Test Account",
			Currency: "USD",
			Type:     "asset",
			Balance: models.Money{
				Amount:   decimal.NewFromInt(100),
				Currency: "USD",
			},
		}

		mockLedger.On("CreateAccount", account).Return(nil)

		body, _ := json.Marshal(account)
		req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		server.CreateAccountHandler(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockLedger.AssertExpectations(t)
	})

	t.Run("invalid json body", func(t *testing.T) {
		server, _ := setupTest()

		req := httptest.NewRequest("POST", "/accounts", bytes.NewBufferString("invalid json"))
		rr := httptest.NewRecorder()

		server.CreateAccountHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("ledger error", func(t *testing.T) {
		server, mockLedger := setupTest()

		account := models.Account{
			ID:       "ACC123",
			Name:     "Test Account",
			Currency: "USD",
			Type:     "asset",
			Balance: models.Money{
				Amount:   decimal.NewFromInt(100),
				Currency: "USD",
			},
		}

		mockLedger.On("CreateAccount", account).Return(fmt.Errorf("account already exists"))

		body, _ := json.Marshal(account)
		req := httptest.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		server.CreateAccountHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockLedger.AssertExpectations(t)
	})
}

// RecordTransactionHandler tests
func TestRecordTransactionHandler(t *testing.T) {
	t.Run("successful transaction recording", func(t *testing.T) {
		server, mockLedger := setupTest()

		tx := models.Transaction{
			ID:            "TX123",
			Description:   "Test Transaction",
			DebitAccount:  "ACC1",
			CreditAccount: "ACC2",
			Amount: models.Money{
				Amount:   decimal.NewFromInt(50),
				Currency: "USD",
			},
		}

		mockLedger.On("RecordTransaction", tx).Return(nil)

		body, _ := json.Marshal(tx)
		req := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		server.RecordTransactionHandler(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockLedger.AssertExpectations(t)
	})

	t.Run("invalid json body", func(t *testing.T) {
		server, _ := setupTest()

		req := httptest.NewRequest("POST", "/transactions", bytes.NewBufferString("invalid json"))
		rr := httptest.NewRecorder()

		server.RecordTransactionHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("ledger error", func(t *testing.T) {
		server, mockLedger := setupTest()

		tx := models.Transaction{
			ID:            "TX123",
			DebitAccount:  "ACC1",
			CreditAccount: "ACC2",
			Amount: models.Money{
				Amount:   decimal.NewFromInt(50),
				Currency: "USD",
			},
		}

		mockLedger.On("RecordTransaction", tx).Return(fmt.Errorf("insufficient funds"))

		body, _ := json.Marshal(tx)
		req := httptest.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		server.RecordTransactionHandler(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockLedger.AssertExpectations(t)
	})
}

// GetBalanceHandler tests
func TestGetBalanceHandler(t *testing.T) {
	t.Run("successful balance retrieval", func(t *testing.T) {
		server, mockLedger := setupTest()

		accountID := "ACC123"
		balance := models.Money{
			Amount:   decimal.NewFromInt(100),
			Currency: "USD",
		}

		mockLedger.On("GetAccountBalance", accountID).Return(balance, nil)

		req := httptest.NewRequest("GET", "/accounts/"+accountID+"/balance", nil)
		req = mux.SetURLVars(req, map[string]string{"accountId": accountID})
		rr := httptest.NewRecorder()

		server.GetBalanceHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response map[string]models.Money
		json.NewDecoder(rr.Body).Decode(&response)
		assert.Equal(t, balance, response["balance"])
		mockLedger.AssertExpectations(t)
	})

	t.Run("account not found", func(t *testing.T) {
		server, mockLedger := setupTest()

		accountID := "NONEXISTENT"
		mockLedger.On("GetAccountBalance", accountID).Return(models.Money{}, fmt.Errorf("account not found"))

		req := httptest.NewRequest("GET", "/accounts/"+accountID+"/balance", nil)
		req = mux.SetURLVars(req, map[string]string{"accountId": accountID})
		rr := httptest.NewRecorder()

		server.GetBalanceHandler(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		mockLedger.AssertExpectations(t)
	})
}

// GetTransactionHistoryHandler tests
func TestGetTransactionHistoryHandler(t *testing.T) {
	t.Run("successful history retrieval", func(t *testing.T) {
		server, mockLedger := setupTest()

		accountID := "ACC123"
		history := []models.Transaction{
			{
				ID:            "TX1",
				DateTime:      time.Now().UTC(),
				Description:   "First Transaction",
				DebitAccount:  "ACC123",
				CreditAccount: "ACC456",
				Amount: models.Money{
					Amount:   decimal.NewFromInt(100),
					Currency: "USD",
				},
			},
			{
				ID:            "TX2",
				DateTime:      time.Now().UTC(),
				Description:   "Second Transaction",
				DebitAccount:  "ACC789",
				CreditAccount: "ACC123",
				Amount: models.Money{
					Amount:   decimal.NewFromInt(50),
					Currency: "USD",
				},
			},
		}

		mockLedger.On("GetTransactionHistory", accountID).Return(history)

		req := httptest.NewRequest("GET", "/accounts/"+accountID+"/history", nil)
		req = mux.SetURLVars(req, map[string]string{"accountId": accountID})
		rr := httptest.NewRecorder()

		server.GetTransactionHistoryHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response []models.Transaction
		json.NewDecoder(rr.Body).Decode(&response)
		assert.Equal(t, len(history), len(response))
		assert.Equal(t, history[0].ID, response[0].ID)
		assert.Equal(t, history[1].ID, response[1].ID)
		mockLedger.AssertExpectations(t)
	})

	t.Run("empty history", func(t *testing.T) {
		server, mockLedger := setupTest()

		accountID := "ACC123"
		mockLedger.On("GetTransactionHistory", accountID).Return([]models.Transaction{})

		req := httptest.NewRequest("GET", "/accounts/"+accountID+"/history", nil)
		req = mux.SetURLVars(req, map[string]string{"accountId": accountID})
		rr := httptest.NewRecorder()

		server.GetTransactionHistoryHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response []models.Transaction
		json.NewDecoder(rr.Body).Decode(&response)
		assert.Empty(t, response)
		mockLedger.AssertExpectations(t)
	})
}
