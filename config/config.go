package config

type Config struct {
	ServerPort   string
	CurrencyFile string
}

func NewConfig() *Config {
	return &Config{
		ServerPort:   ":8080",
		CurrencyFile: "data/iso4217_currency_dev.json",
	}
}
