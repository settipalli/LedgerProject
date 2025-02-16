package config

import "time"

type Config struct {
	ServerPort        string
	CurrencyFile      string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	MaxHeaderBytes    int
}

func NewConfig() *Config {
	return &Config{
		ServerPort:        ":8080",
		CurrencyFile:      "data/iso4217_currency_dev.json",
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}
}
