package ledger

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"ledgerproject/models"
	"ledgerproject/services"
	"log"
	"sync"
	"time"
)

type ledger struct {
	accounts          map[string]*models.Account
	transactions      []models.Transaction
	currencyValidator *services.CurrencyValidator
	mu                sync.RWMutex
}

func NewLedger(cv *services.CurrencyValidator) LedgerService {
	l := &ledger{
		accounts:          make(map[string]*models.Account),
		transactions:      []models.Transaction{},
		currencyValidator: cv,
	}

	// Start periodic balance checking
	ctx := context.Background()
	go l.PerformPeriodicBalanceCheck(ctx)

	return l
}

func (l *ledger) CreateAccount(account models.Account) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Validate if account already exists
	if _, exists := l.accounts[account.ID]; exists {
		return fmt.Errorf("account %s already exists", account.ID)
	}

	// Validate currency
	if !l.currencyValidator.IsValid(account.Currency) {
		return fmt.Errorf("invalid currency code: %s", account.Currency)
	}

	// Initialize balance if zero
	if account.Balance.Amount.IsZero() {
		account.Balance.Amount = decimal.Zero
		account.Balance.Currency = account.Currency
	} else {
		// Validate that balance currency matches account currency
		if account.Balance.Currency != account.Currency {
			return fmt.Errorf("balance currency (%s) does not match account currency (%s)",
				account.Balance.Currency, account.Currency)
		}
	}

	account.CreateDateTime = time.Now().UTC()
	l.accounts[account.ID] = &account
	return nil
}

func (l *ledger) RecordTransaction(tx models.Transaction) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Store initial balances for verification
	initialDebitBalance := decimal.Zero
	initialCreditBalance := decimal.Zero

	debitAcc, exists := l.accounts[tx.DebitAccount]
	if !exists {
		return fmt.Errorf("debit account %s does not exist", tx.DebitAccount)
	}
	initialDebitBalance = debitAcc.Balance.Amount

	creditAcc, exists := l.accounts[tx.CreditAccount]
	if !exists {
		return fmt.Errorf("credit account %s does not exist", tx.CreditAccount)
	}
	initialCreditBalance = creditAcc.Balance.Amount

	if debitAcc.Currency != tx.Amount.Currency || creditAcc.Currency != tx.Amount.Currency {
		return fmt.Errorf("currency mismatch between accounts and transaction")
	}

	// Check if debit account has sufficient funds
	if debitAcc.Balance.Amount.LessThan(tx.Amount.Amount) {
		return fmt.Errorf("insufficient funds in debit account %s", tx.DebitAccount)
	}

	// Perform the transaction
	debitAcc.Balance.Amount = debitAcc.Balance.Amount.Sub(tx.Amount.Amount)
	creditAcc.Balance.Amount = creditAcc.Balance.Amount.Add(tx.Amount.Amount)

	// Verify the books are balanced
	totalChange := initialDebitBalance.Sub(debitAcc.Balance.Amount).
		Sub(creditAcc.Balance.Amount.Sub(initialCreditBalance))

	if !totalChange.IsZero() {
		// Rollback the transaction
		debitAcc.Balance.Amount = initialDebitBalance
		creditAcc.Balance.Amount = initialCreditBalance
		return fmt.Errorf("transaction failed: books would be unbalanced by %s %s",
			totalChange.String(), tx.Amount.Currency)
	}

	// Record the transaction
	tx.DateTime = time.Now().UTC()
	l.transactions = append(l.transactions, tx)

	return nil
}

func (l *ledger) GetAccountBalance(accountID string) (models.Money, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	account, exists := l.accounts[accountID]
	if !exists {
		return models.Money{}, fmt.Errorf("account %s does not exist", accountID)
	}

	return account.Balance, nil
}

func (l *ledger) GetTransactionHistory(accountID string) []models.Transaction {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var history []models.Transaction
	for _, tx := range l.transactions {
		if tx.DebitAccount == accountID || tx.CreditAccount == accountID {
			history = append(history, tx)
		}
	}
	return history
}

func (l *ledger) PerformPeriodicBalanceCheck(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := l.VerifyLedgerBalance(); err != nil {
				// Log the error or trigger an alert
				log.Printf("CRITICAL: Ledger balance check failed: %v", err)
				// Could also trigger notifications to administrators
			}
		}
	}
}

func (l *ledger) VerifyLedgerBalance() error {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Group accounts by currency
	balancesByCurrency := make(map[string]decimal.Decimal)

	for _, account := range l.accounts {
		curr := account.Currency
		if _, exists := balancesByCurrency[curr]; !exists {
			balancesByCurrency[curr] = decimal.Zero
		}
		balancesByCurrency[curr] = balancesByCurrency[curr].Add(account.Balance.Amount)
	}

	// Check if each currency's total balance is zero
	for currency, total := range balancesByCurrency {
		if !total.IsZero() {
			return fmt.Errorf("ledger is unbalanced for %s: total balance is %s",
				currency, total.String())
		}
	}

	return nil
}
