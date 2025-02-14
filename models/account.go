package models

import "time"

type Account struct {
    ID             string    `json:"id"`
    Name           string    `json:"name"`
    Balance        Money     `json:"balance"`
    Type           string    `json:"type"`
    Currency       string    `json:"currency"`
    CreateDateTime time.Time `json:"datetime"`
}
