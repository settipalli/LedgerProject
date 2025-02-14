package ledger

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"ledgerproject/models"
	"testing"
	"time"
)

func TestNewLedger(t *testing.T) {
	l := NewLedger()
	assert.NotNil(t, l, "NewLedger should return a non-nil ledger")

	// Test initial ledger balance
	err := l.(*ledger).VerifyLedgerBalance()
	assert.NoError(t, err, "New ledger should be balanced")
}

func TestLedger_CreateAccount(t *testing.T) {
	l := NewLedger()

	tests := []struct {
		name    string
		account models.Account
		wantErr bool
	}{
		{
			name: "Valid account creation",
			account: models.Account{
				ID:       "ACC001",
				Name:     "Test Account",
				Currency: "USD",
				Type:     "Asset",
				Balance: models.Money{
					Amount:   decimal.NewFromFloat(0),
					Currency: "USD",
				},
			},
			wantErr: false,
		},
		{
			name: "Duplicate account",
			account: models.Account{
				ID:       "ACC001",
				Name:     "Duplicate Account",
				Currency: "USD",
				Type:     "Asset",
			},
			wantErr: true,
		},
		{
			name: "Account with initial balance",
			account: models.Account{
				ID:       "ACC002",
				Name:     "Account with Balance",
				Currency: "EUR",
				Type:     "Asset",
				Balance: models.Money{
					Amount:   decimal.NewFromFloat(100.50),
					Currency: "EUR",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := l.CreateAccount(tt.account)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLedger_RecordTransaction(t *testing.T) {
	l := NewLedger()

	// Setup test accounts
	acc1 := models.Account{
		ID:       "ACC001",
		Name:     "Account 1",
		Currency: "USD",
		Type:     "Asset",
		Balance: models.Money{
			Amount:   decimal.NewFromFloat(1000),
			Currency: "USD",
		},
	}
	acc2 := models.Account{
		ID:       "ACC002",
		Name:     "Account 2",
		Currency: "USD",
		Type:     "Asset",
		Balance: models.Money{
			Amount:   decimal.NewFromFloat(500),
			Currency: "USD",
		},
	}
	acc3 := models.Account{
		ID:       "ACC003",
		Name:     "Account 3",
		Currency: "EUR",
		Type:     "Asset",
	}

	_ = l.CreateAccount(acc1)
	_ = l.CreateAccount(acc2)
	_ = l.CreateAccount(acc3)

	tests := []struct {
		name    string
		tx      models.Transaction
		wantErr bool
	}{
		{
			name: "Valid transaction",
			tx: models.Transaction{
				ID:            "TX001",
				Description:   "Test transaction",
				DebitAccount:  "ACC002",
				CreditAccount: "ACC001",
				Money: models.Money{
					Amount:   decimal.NewFromFloat(500),
					Currency: "USD",
				},
			},
			wantErr: false,
		},
		{
			name: "Non-existent debit account",
			tx: models.Transaction{
				ID:            "TX002",
				DebitAccount:  "NONEXISTENT",
				CreditAccount: "ACC001",
				Money: models.Money{
					Amount:   decimal.NewFromFloat(100),
					Currency: "USD",
				},
			},
			wantErr: true,
		},
		{
			name: "Currency mismatch",
			tx: models.Transaction{
				ID:            "TX003",
				DebitAccount:  "ACC001",
				CreditAccount: "ACC003",
				Money: models.Money{
					Amount:   decimal.NewFromFloat(100),
					Currency: "USD",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := l.RecordTransaction(tt.tx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLedger_GetAccountBalance(t *testing.T) {
	l := NewLedger()

	// Setup test account
	acc := models.Account{
		ID:       "ACC001",
		Name:     "Test Account",
		Currency: "USD",
		Type:     "Asset",
		Balance: models.Money{
			Amount:   decimal.NewFromFloat(1000),
			Currency: "USD",
		},
	}
	_ = l.CreateAccount(acc)

	tests := []struct {
		name      string
		accountID string
		want      models.Money
		wantErr   bool
	}{
		{
			name:      "Existing account",
			accountID: "ACC001",
			want: models.Money{
				Amount:   decimal.NewFromFloat(1000),
				Currency: "USD",
			},
			wantErr: false,
		},
		{
			name:      "Non-existent account",
			accountID: "NONEXISTENT",
			want:      models.Money{},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balance, err := l.GetAccountBalance(tt.accountID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, balance)
			}
		})
	}
}

func TestLedger_GetTransactionHistory(t *testing.T) {
	l := NewLedger()

	// Setup test accounts
	acc1 := models.Account{
		ID:   "ACC001",
		Name: "Test Account 1",
		Balance: models.Money{
			Amount:   decimal.NewFromFloat(1000),
			Currency: "USD",
		},
		Type:           "Asset",
		Currency:       "USD",
		CreateDateTime: time.Time{},
	}
	acc2 := models.Account{
		ID:   "ACC002",
		Name: "Test Account 2",
		Balance: models.Money{
			Amount:   decimal.NewFromFloat(500),
			Currency: "USD",
		},
		Type:           "Expense",
		Currency:       "USD",
		CreateDateTime: time.Time{},
	}
	_ = l.CreateAccount(acc1)
	_ = l.CreateAccount(acc2)

	// Record some transactions
	tx1 := models.Transaction{
		ID:            "TX001",
		DebitAccount:  "ACC001",
		CreditAccount: "ACC002",
		Money: models.Money{
			Amount:   decimal.NewFromFloat(100),
			Currency: "USD",
		},
	}
	tx2 := models.Transaction{
		ID:            "TX002",
		DebitAccount:  "ACC002",
		CreditAccount: "ACC001",
		Money: models.Money{
			Amount:   decimal.NewFromFloat(50),
			Currency: "USD",
		},
	}

	_ = l.RecordTransaction(tx1)
	_ = l.RecordTransaction(tx2)

	tests := []struct {
		name      string
		accountID string
		wantCount int
	}{
		{
			name:      "Account with transactions",
			accountID: "ACC001",
			wantCount: 2,
		},
		{
			name:      "Non-existent account",
			accountID: "NONEXISTENT",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			history := l.GetTransactionHistory(tt.accountID)
			assert.Equal(t, tt.wantCount, len(history))
		})
	}
}

func TestLedger_ConcurrentOperations(t *testing.T) {
	l := NewLedger()

	// Setup initial accounts
	acc1 := models.Account{
		ID:       "ACC001",
		Currency: "USD",
		Balance: models.Money{
			Amount:   decimal.NewFromFloat(1000),
			Currency: "USD",
		},
	}
	acc2 := models.Account{
		ID:       "ACC002",
		Currency: "USD",
		Balance: models.Money{
			Amount:   decimal.NewFromFloat(1000),
			Currency: "USD",
		},
	}

	_ = l.CreateAccount(acc1)
	_ = l.CreateAccount(acc2)

	// Test concurrent transactions
	const numGoroutines = 10
	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			tx := models.Transaction{
				ID:            fmt.Sprintf("TX%d", idx),
				DebitAccount:  "ACC001",
				CreditAccount: "ACC002",
				Money: models.Money{
					Amount:   decimal.NewFromFloat(10),
					Currency: "USD",
				},
			}
			_ = l.RecordTransaction(tx)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify final balances
	bal1, _ := l.GetAccountBalance("ACC001")
	bal2, _ := l.GetAccountBalance("ACC002")

	// ACC001 should have lost 100 (10 * numGoroutines)
	assert.True(t, bal1.Amount.Equal(decimal.NewFromFloat(900)))
	// ACC002 should have gained 100
	assert.True(t, bal2.Amount.Equal(decimal.NewFromFloat(1100)))
}

func TestLedger_RecordTransaction_BalanceChecks(t *testing.T) {
	l := NewLedger()

	// Setup test accounts with initial balances
	accounts := []models.Account{
		{
			ID:       "ASSET",
			Name:     "Asset Account",
			Type:     "ASSET",
			Currency: "USD",
			Balance: models.Money{
				Amount:   decimal.NewFromFloat(1000),
				Currency: "USD",
			},
		},
		{
			ID:       "LIABILITY",
			Name:     "Liability Account",
			Type:     "LIABILITY",
			Currency: "USD",
			Balance: models.Money{
				Amount:   decimal.NewFromFloat(-1000), // Offset to maintain initial balance
				Currency: "USD",
			},
		},
		{
			ID:       "EMPTY",
			Name:     "Empty Account",
			Type:     "ASSET",
			Currency: "USD",
			Balance: models.Money{
				Amount:   decimal.Zero,
				Currency: "USD",
			},
		},
	}

	// Create accounts
	for _, acc := range accounts {
		err := l.CreateAccount(acc)
		require.NoError(t, err)
	}

	// Verify initial balance
	err := l.(*ledger).VerifyLedgerBalance()
	require.NoError(t, err, "Initial ledger should be balanced")

	tests := []struct {
		name    string
		tx      models.Transaction
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid balanced transaction",
			tx: models.Transaction{
				ID:            "TX001",
				Description:   "Valid transfer",
				DebitAccount:  "ASSET",
				CreditAccount: "LIABILITY",
				Money: models.Money{
					Amount:   decimal.NewFromFloat(100),
					Currency: "USD",
				},
			},
			wantErr: false,
		},
		{
			name: "Insufficient funds",
			tx: models.Transaction{
				ID:            "TX002",
				Description:   "Transfer more than balance",
				DebitAccount:  "EMPTY",
				CreditAccount: "LIABILITY",
				Money: models.Money{
					Amount:   decimal.NewFromFloat(1000),
					Currency: "USD",
				},
			},
			wantErr: true,
			errMsg:  "insufficient funds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := l.RecordTransaction(tt.tx)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify ledger remains balanced after each transaction
			balanceErr := l.(*ledger).VerifyLedgerBalance()
			assert.NoError(t, balanceErr, "Ledger should remain balanced after transaction")
		})
	}
}
