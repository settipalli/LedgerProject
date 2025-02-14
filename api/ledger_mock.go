// used by server_test.go and handlers_test.go
package api

import (
	"context"
	"github.com/stretchr/testify/mock"
	"ledgerproject/models"
)

// MockLedger implements the LedgerService interface for testing
type MockLedger struct {
	mock.Mock
}

func (m *MockLedger) CreateAccount(account models.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockLedger) RecordTransaction(tx models.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockLedger) GetAccountBalance(accountID string) (models.Money, error) {
	args := m.Called(accountID)
	return args.Get(0).(models.Money), args.Error(1)
}

func (m *MockLedger) GetTransactionHistory(accountID string) []models.Transaction {
	args := m.Called(accountID)
	return args.Get(0).([]models.Transaction)
}

func (m *MockLedger) VerifyLedgerBalance() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockLedger) PerformPeriodicBalanceCheck(ctx context.Context) {
	m.Called(ctx)
}
