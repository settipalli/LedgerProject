package ledger

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"ledgerproject/config"
	"ledgerproject/logger"
	"ledgerproject/models"
	"ledgerproject/services"
	"testing"
	"time"
)

type testSetup struct {
	ledger      LedgerService
	validator   *services.CurrencyValidator
	validCurr   string
	invalidCurr string
}

// setupTestLogger initializes a test logger
func setupTestLogger(t *testing.T) *zap.Logger {
	testLogger := zaptest.NewLogger(t)
	// Initialize the package-level logger
	logger.Init(true) // Set to development mode
	return testLogger
}

func setupTest(t *testing.T) *testSetup {
	// Create test config with test currency file
	cfg := &config.Config{
		CurrencyFile: "../data/iso4217_currency_test.json",
	}

	// Setup test logger
	testLogger := setupTestLogger(t)
	defer testLogger.Sync()

	// Initialize currency validator
	validator, err := services.NewCurrencyValidator(cfg)
	require.NoError(t, err)

	// Create test ledger with validator
	testLedger := NewLedger(validator)

	return &testSetup{
		ledger:      testLedger,
		validator:   validator,
		validCurr:   "USD", // USD should be present in the test currency json file
		invalidCurr: "INVALID",
	}
}

func TestCreateAccount(t *testing.T) {
	setup := setupTest(t)

	tests := []struct {
		name    string
		account models.Account
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid account with valid currency",
			account: models.Account{
				ID:       "ACC001",
				Name:     "Test Account",
				Type:     "Expense",
				Currency: setup.validCurr,
				Balance: models.Money{
					Amount:   decimal.Zero,
					Currency: setup.validCurr,
				},
			},
			wantErr: false,
		},
		{
			name: "Valid account and valid currency with balance",
			account: models.Account{
				ID:       "ACC002",
				Name:     "Test Account",
				Currency: setup.validCurr,
				Type:     "Asset",
				Balance: models.Money{
					Amount:   decimal.NewFromFloat(10101.00),
					Currency: setup.validCurr,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid currency",
			account: models.Account{
				ID:       "ACC003",
				Name:     "Invalid Currency Account",
				Currency: setup.invalidCurr,
			},
			wantErr: true,
			errMsg:  "invalid currency code: INVALID",
		},
		{
			name: "Currency mismatch between account and balance",
			account: models.Account{
				ID:       "ACC004",
				Name:     "Mismatched Currency Account",
				Currency: setup.validCurr,
				Balance: models.Money{
					Amount:   decimal.NewFromInt(100),
					Currency: "EUR", // assuming EUR is valid (present in test currency json file) and different from account currency
				},
			},
			wantErr: true,
			errMsg:  "balance currency (EUR) does not match account currency (USD)",
		},
		{
			name: "Duplicate account ID",
			account: models.Account{
				ID:       "ACC001", // same as first test case
				Name:     "Duplicate Account",
				Currency: setup.validCurr,
			},
			wantErr: true,
			errMsg:  "account ACC001 already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setup.ledger.CreateAccount(tt.account)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)

				// Verify account was created correctly
				balance, err := setup.ledger.GetAccountBalance(tt.account.ID)
				assert.NoError(t, err)
				assert.Equal(t, tt.account.Currency, balance.Currency)
			}
		})
	}
}

func TestRecordTransaction(t *testing.T) {
	setup := setupTest(t)

	// Create test accounts
	accounts := []models.Account{
		{
			ID:       "ACC001",
			Name:     "Account 1",
			Currency: setup.validCurr,
			Balance: models.Money{
				Amount:   decimal.NewFromInt(1000),
				Currency: setup.validCurr,
			},
		},
		{
			ID:       "ACC002",
			Name:     "Account 2",
			Currency: setup.validCurr,
			Balance: models.Money{
				Amount:   decimal.Zero,
				Currency: setup.validCurr,
			},
		},
	}

	// Create accounts
	for _, acc := range accounts {
		err := setup.ledger.CreateAccount(acc)
		require.NoError(t, err)
	}

	tests := []struct {
		name    string
		tx      models.Transaction
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid transaction",
			tx: models.Transaction{
				ID:            "TX001",
				Description:   "Test Transaction",
				DebitAccount:  "ACC001",
				CreditAccount: "ACC002",
				Amount: models.Money{
					Amount:   decimal.NewFromInt(500),
					Currency: setup.validCurr,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid currency transaction",
			tx: models.Transaction{
				ID:            "TX002",
				Description:   "Invalid Currency Transaction",
				DebitAccount:  "ACC001",
				CreditAccount: "ACC002",
				Amount: models.Money{
					Amount:   decimal.NewFromInt(100),
					Currency: setup.invalidCurr,
				},
			},
			wantErr: true,
			errMsg:  "currency mismatch between accounts and transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setup.ledger.RecordTransaction(tt.tx)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)

				// Verify transaction was recorded correctly
				history := setup.ledger.GetTransactionHistory(tt.tx.DebitAccount)
				assert.Condition(t, func() bool {
					for _, hist := range history {
						if hist.ID == tt.tx.ID &&
							hist.Description == tt.tx.Description &&
							hist.DebitAccount == tt.tx.DebitAccount &&
							hist.CreditAccount == tt.tx.CreditAccount &&
							hist.Amount == tt.tx.Amount {
							return true
						}
					}
					return false
				})
			}
		})
	}
}

func TestVerifyLedgerBalance(t *testing.T) {
	setup := setupTest(t)

	// Create two accounts with initial balances
	accounts := []models.Account{
		{
			ID:       "ACC001",
			Name:     "Account 1",
			Currency: setup.validCurr,
			Balance: models.Money{
				Amount:   decimal.NewFromInt(1000),
				Currency: setup.validCurr,
			},
		},
		{
			ID:       "ACC002",
			Name:     "Account 2",
			Currency: setup.validCurr,
			Balance: models.Money{
				Amount:   decimal.NewFromInt(-1000),
				Currency: setup.validCurr,
			},
		},
	}

	// Create accounts
	for _, acc := range accounts {
		err := setup.ledger.CreateAccount(acc)
		require.NoError(t, err)
	}

	// Verify ledger balance
	err := setup.ledger.VerifyLedgerBalance()
	assert.NoError(t, err)
}

func TestPeriodicBalanceCheck(t *testing.T) {
	setup := setupTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start periodic balance check
	go setup.ledger.PerformPeriodicBalanceCheck(ctx)

	// Create balanced accounts
	accounts := []models.Account{
		{
			ID:       "ACC001",
			Name:     "Account 1",
			Currency: setup.validCurr,
			Balance: models.Money{
				Amount:   decimal.NewFromInt(1000),
				Currency: setup.validCurr,
			},
		},
		{
			ID:       "ACC002",
			Name:     "Account 2",
			Currency: setup.validCurr,
			Balance: models.Money{
				Amount:   decimal.NewFromInt(-1000),
				Currency: setup.validCurr,
			},
		},
	}

	// Create accounts
	for _, acc := range accounts {
		err := setup.ledger.CreateAccount(acc)
		require.NoError(t, err)
	}

	// Let the periodic check run at least once
	time.Sleep(100 * time.Millisecond)
}
