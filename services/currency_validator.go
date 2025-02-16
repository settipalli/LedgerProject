package services

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"ledgerproject/config"
	"ledgerproject/logger"
	"os"
	"sync"
)

type CurrencyValidator struct {
	validCurrencies map[string]struct{}
	mu              sync.RWMutex
	config          *config.Config
}

type currencyData struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Number string `json:"number"`
}

func NewCurrencyValidator(config *config.Config) (*CurrencyValidator, error) {
	cv := &CurrencyValidator{
		validCurrencies: make(map[string]struct{}),
		config:          config,
	}

	if err := cv.loadCurrencies(); err != nil {
		return nil, fmt.Errorf("failed to load currencies: %v", err)
	}

	return cv, nil
}

func (cv *CurrencyValidator) loadCurrencies() error {
	log := logger.Get()
	// Read the ISO 4217 currency codes from JSON file
	file, err := os.ReadFile(cv.config.CurrencyFile)
	if err != nil {
		log.Error("Failed to read currency file",
			zap.Error(err),
			zap.String("file", cv.config.CurrencyFile))
		return fmt.Errorf("error reading currency file: %v", err)
	}

	var data struct {
		Currencies []currencyData `json:"currencies"`
	}

	if err := json.Unmarshal(file, &data); err != nil {
		log.Error("Failed to unmarshal currency data",
			zap.Error(err))
		return fmt.Errorf("error unmarshaling currencies: %v", err)
	}

	cv.mu.Lock()
	defer cv.mu.Unlock()

	for _, currency := range data.Currencies {
		cv.validCurrencies[currency.Code] = struct{}{}
	}

	log.Info("Currency data loaded successfully",
		zap.Int("currency_count", len(cv.validCurrencies)))
	return nil
}

func (cv *CurrencyValidator) IsValid(code string) bool {
	log := logger.Get()
	cv.mu.RLock()
	defer cv.mu.RUnlock()

	_, exists := cv.validCurrencies[code]

	log.Info("Currency code is valid", zap.String("currency_code", code))
	return exists
}
