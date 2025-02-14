package models

import (
    "encoding/json"
    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "testing"
)

func TestMoney_Creation(t *testing.T) {
    tests := []struct {
        name     string
        amount   string
        currency string
        wantErr  bool
    }{
        {
            name:     "valid positive amount",
            amount:   "123.45",
            currency: "USD",
            wantErr:  false,
        },
        {
            name:     "valid negative amount",
            amount:   "-123.45",
            currency: "EUR",
            wantErr:  false,
        },
        {
            name:     "zero amount",
            amount:   "0.00",
            currency: "GBP",
            wantErr:  false,
        },
        {
            name:     "high precision amount",
            amount:   "123.45678901234567890",
            currency: "BTC",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            amount, err := decimal.NewFromString(tt.amount)
            require.NoError(t, err)

            money := Money{
                Amount:   amount,
                Currency: tt.currency,
            }

            assert.Equal(t, tt.currency, money.Currency)
            assert.True(t, money.Amount.Equal(amount))
        })
    }
}

func TestMoney_MarshalJSON(t *testing.T) {
    tests := []struct {
        name    string
        money   Money
        want    string
        wantErr bool
    }{
        {
            name: "positive amount",
            money: Money{
                Amount:   decimal.NewFromFloat(123.45),
                Currency: "USD",
            },
            want:    `{"amount":"123.45","currency":"USD"}`,
            wantErr: false,
        },
        {
            name: "negative amount",
            money: Money{
                Amount:   decimal.NewFromFloat(-123.45),
                Currency: "EUR",
            },
            want:    `{"amount":"-123.45","currency":"EUR"}`,
            wantErr: false,
        },
        {
            name: "zero amount",
            money: Money{
                Amount:   decimal.Zero,
                Currency: "GBP",
            },
            want:    `{"amount":"0","currency":"GBP"}`,
            wantErr: false,
        },
        {
            name: "high precision",
            money: Money{
                Amount:   decimal.NewFromFloat(123.45678901234567890),
                Currency: "BTC",
            },
            want:    `{"amount":"123.45678901234568","currency":"BTC"}`,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := json.Marshal(tt.money)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            assert.JSONEq(t, tt.want, string(got))
        })
    }
}

func TestMoney_UnmarshalJSON(t *testing.T) {
    tests := []struct {
        name    string
        json    string
        want    Money
        wantErr bool
    }{
        {
            name: "valid positive amount",
            json: `{"amount":"123.45","currency":"USD"}`,
            want: Money{
                Amount:   decimal.NewFromFloat(123.45),
                Currency: "USD",
            },
            wantErr: false,
        },
        {
            name: "valid negative amount",
            json: `{"amount":"-123.45","currency":"EUR"}`,
            want: Money{
                Amount:   decimal.NewFromFloat(-123.45),
                Currency: "EUR",
            },
            wantErr: false,
        },
        {
            name: "zero amount",
            json: `{"amount":"0","currency":"GBP"}`,
            want: Money{
                Amount:   decimal.Zero,
                Currency: "GBP",
            },
            wantErr: false,
        },
        {
            name:    "invalid amount format",
            json:    `{"amount":"not-a-number","currency":"USD"}`,
            wantErr: true,
        },
        {
            name:    "missing amount",
            json:    `{"currency":"USD"}`,
            wantErr: true,
        },
        {
            name:    "missing currency",
            json:    `{"amount":"123.45"}`,
            wantErr: true,
        },
        {
            name:    "invalid json",
            json:    `{"amount":}`,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var got Money
            err := json.Unmarshal([]byte(tt.json), &got)

            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            assert.True(t, got.Amount.Equal(tt.want.Amount))
            assert.Equal(t, tt.want.Currency, got.Currency)
        })
    }
}

func TestMoney_Comparisons(t *testing.T) {
    m1 := Money{Amount: decimal.NewFromFloat(100.00), Currency: "USD"}
    m2 := Money{Amount: decimal.NewFromFloat(100.00), Currency: "USD"}
    m3 := Money{Amount: decimal.NewFromFloat(200.00), Currency: "USD"}
    m4 := Money{Amount: decimal.NewFromFloat(100.00), Currency: "EUR"}

    t.Run("equal amounts same currency", func(t *testing.T) {
        assert.True(t, m1.Amount.Equal(m2.Amount))
        assert.Equal(t, m1.Currency, m2.Currency)
    })

    t.Run("different amounts same currency", func(t *testing.T) {
        assert.False(t, m1.Amount.Equal(m3.Amount))
        assert.Equal(t, m1.Currency, m3.Currency)
    })

    t.Run("equal amounts different currency", func(t *testing.T) {
        assert.True(t, m1.Amount.Equal(m4.Amount))
        assert.NotEqual(t, m1.Currency, m4.Currency)
    })
}

func TestMoney_String(t *testing.T) {
    tests := []struct {
        name  string
        money Money
        want  string
    }{
        {
            name: "positive amount",
            money: Money{
                Amount:   decimal.NewFromFloat(123.45),
                Currency: "USD",
            },
            want: "123.45 USD",
        },
        {
            name: "negative amount",
            money: Money{
                Amount:   decimal.NewFromFloat(-123.45),
                Currency: "EUR",
            },
            want: "-123.45 EUR",
        },
        {
            name: "zero amount",
            money: Money{
                Amount:   decimal.Zero,
                Currency: "GBP",
            },
            want: "0 GBP",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := tt.money.String()
            assert.Equal(t, tt.want, got)
        })
    }
}

// Add this method to money.go if not already present
func (m Money) String() string {
    return m.Amount.String() + " " + m.Currency
}

// Benchmark tests
func BenchmarkMoney_MarshalJSON(b *testing.B) {
    money := Money{
        Amount:   decimal.NewFromFloat(123.45),
        Currency: "USD",
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := json.Marshal(money)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkMoney_UnmarshalJSON(b *testing.B) {
    data := []byte(`{"amount":"123.45","currency":"USD"}`)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var money Money
        err := json.Unmarshal(data, &money)
        if err != nil {
            b.Fatal(err)
        }
    }
}
