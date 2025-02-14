package config

import (
	"go.uber.org/fx"
)

func Development() fx.Option {
	return fx.Options(
		fx.Provide(func() *Config {
			return &Config{
				ServerPort:   ":8080",
				CurrencyFile: "data/iso4217_currency_dev.json",
			}
		}),
	)
}

func Test() fx.Option {
	return fx.Options(
		fx.Provide(func() *Config {
			return &Config{
				ServerPort:   ":8081",
				CurrencyFile: "data/iso4217_currency_test.json",
			}
		}),
	)
}

func Production() fx.Option {
	return fx.Options(
		fx.Provide(func() *Config {
			return &Config{
				ServerPort:   ":80",
				CurrencyFile: "data/iso4217_currency.json",
			}
		}),
	)
}
