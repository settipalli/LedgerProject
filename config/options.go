package config

import (
	"go.uber.org/fx"
	"time"
)

func Development() fx.Option {
	return fx.Options(
		fx.Provide(func() *Config {
			return &Config{
				ServerPort:        ":8080",
				CurrencyFile:      "data/iso4217_currency_dev.json",
				ReadTimeout:       15 * time.Second,
				WriteTimeout:      15 * time.Second,
				IdleTimeout:       60 * time.Second,
				ReadHeaderTimeout: 5 * time.Second,
				MaxHeaderBytes:    1 << 20,
			}
		}),
	)
}

func Test() fx.Option {
	return fx.Options(
		fx.Provide(func() *Config {
			return &Config{
				ServerPort:        ":8081",
				CurrencyFile:      "data/iso4217_currency_test.json",
				ReadTimeout:       5 * time.Second, // Shorter timeouts for testing
				WriteTimeout:      5 * time.Second,
				IdleTimeout:       30 * time.Second,
				ReadHeaderTimeout: 2 * time.Second,
				MaxHeaderBytes:    1 << 20,
			}
		}),
	)
}

func Production() fx.Option {
	return fx.Options(
		fx.Provide(func() *Config {
			return &Config{
				ServerPort:        ":80",
				CurrencyFile:      "data/iso4217_currency.json",
				ReadTimeout:       30 * time.Second, // Longer timeouts for production
				WriteTimeout:      30 * time.Second,
				IdleTimeout:       120 * time.Second,
				ReadHeaderTimeout: 10 * time.Second,
				MaxHeaderBytes:    1 << 20,
			}
		}),
	)
}
