package models

import "time"

type Transaction struct {
	ID            string    `json:"id"`
	DateTime      time.Time `json:"datetime"`
	Description   string    `json:"description"`
	DebitAccount  string    `json:"debit_account"`
	CreditAccount string    `json:"credit_account"`
	Amount        Money     `json:"amount"`
}
