package ledger

import (
	"context"
	"ledgerproject/models"
)

type LedgerService interface {
	CreateAccount(account models.Account) error
	RecordTransaction(tx models.Transaction) error
	GetAccountBalance(accountID string) (models.Money, error)
	GetTransactionHistory(accountID string) []models.Transaction
	VerifyLedgerBalance() error
	PerformPeriodicBalanceCheck(context.Context)
}
