package models

import (
    "encoding/json"
    "fmt"
    "github.com/shopspring/decimal"
)

type Money struct {
    Amount   decimal.Decimal `json:"amount"`
    Currency string          `json:"currency"`
}

func (m Money) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        Amount   string `json:"amount"`
        Currency string `json:"currency"`
    }{
        Amount:   m.Amount.String(),
        Currency: m.Currency,
    })
}

func (m *Money) UnmarshalJSON(data []byte) error {
    var temp struct {
        Amount   string `json:"amount"`
        Currency string `json:"currency"`
    }
    if err := json.Unmarshal(data, &temp); err != nil {
        return err
    }

    amount, err := decimal.NewFromString(temp.Amount)
    if err != nil {
        return fmt.Errorf("invalid amount format: %v", err)
    }

    if temp.Currency == "" {
        return fmt.Errorf("currency is required and cannot be empty")
    }

    m.Amount = amount
    m.Currency = temp.Currency
    return nil
}
