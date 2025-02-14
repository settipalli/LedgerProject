package services

import (
	"ledgerproject/config"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// TestData represents sample currency data for testing
const TestData = `{
    "currencies": [
        {"code": "USD", "name": "US Dollar", "number": "840"},
        {"code": "EUR", "name": "Euro", "number": "978"},
        {"code": "GBP", "name": "British Pound", "number": "826"}
    ]
}`

func setupTestData(t *testing.T) (string, *config.Config, func()) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "currency_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create data directory
	dataDir := filepath.Join(tmpDir, "data")
	if err := os.Mkdir(dataDir, 0755); err != nil {
		t.Fatalf("Failed to create data directory: %v", err)
	}

	// Create test JSON file
	filePath := filepath.Join(dataDir, "iso4217_currency.json")
	if err := os.WriteFile(filePath, []byte(TestData), 0644); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	cfg := &config.Config{
		ServerPort:   ":8080",
		CurrencyFile: filePath,
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cfg, cleanup
}

func TestNewCurrencyValidator(t *testing.T) {
	// Setup test data
	tmpDir, cfg, cleanup := setupTestData(t)
	defer cleanup()

	// Change working directory to temp directory
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
	}{
		{
			name:    "Successfully create new validator",
			config:  cfg,
			wantErr: false,
		},
		{
			name: "Invalid file path",
			config: &config.Config{
				CurrencyFile: "nonexistent/path/currency.json",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv, err := NewCurrencyValidator(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCurrencyValidator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && cv == nil {
				t.Error("NewCurrencyValidator() returned nil validator")
			}
		})
	}
}

func TestCurrencyValidator_IsValid(t *testing.T) {
	// Setup test data
	tmpDir, cfg, cleanup := setupTestData(t)
	defer cleanup()

	// Change working directory to temp directory
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	cv, err := NewCurrencyValidator(cfg)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "Valid currency USD",
			code: "USD",
			want: true,
		},
		{
			name: "Valid currency EUR",
			code: "EUR",
			want: true,
		},
		{
			name: "Invalid currency",
			code: "XXX",
			want: false,
		},
		{
			name: "Empty currency code",
			code: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cv.IsValid(tt.code); got != tt.want {
				t.Errorf("CurrencyValidator.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrencyValidator_Concurrent(t *testing.T) {
	// Setup test data
	tmpDir, cfg, cleanup := setupTestData(t)
	defer cleanup()

	// Change working directory to temp directory
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	cv, err := NewCurrencyValidator(cfg)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			// Test concurrent reads
			_ = cv.IsValid("USD")
			_ = cv.IsValid("EUR")
			_ = cv.IsValid("GBP")
		}()
	}

	wg.Wait()
}

func TestCurrencyValidator_FileNotFound(t *testing.T) {
	// Test with non-existent file
	originalWd, _ := os.Getwd()
	tmpDir, err := os.MkdirTemp("", "currency_test_error")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	_, err = NewCurrencyValidator(&config.Config{
		CurrencyFile: "nonexistent/path/currency.json",
	})
	if err == nil {
		t.Error("Expected error when currency file doesn't exist")
	}
}

func TestCurrencyValidator_InvalidJSON(t *testing.T) {
	// Setup invalid JSON data
	tmpDir, err := os.MkdirTemp("", "currency_test_invalid")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create data directory
	dataDir := filepath.Join(tmpDir, "data")
	if err := os.Mkdir(dataDir, 0755); err != nil {
		t.Fatalf("Failed to create data directory: %v", err)
	}

	// Create invalid JSON file
	filePath := filepath.Join(dataDir, "iso4217_currency.json")
	if err := os.WriteFile(filePath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	// Change working directory
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	_, err = NewCurrencyValidator(&config.Config{
		CurrencyFile: filePath,
	})
	if err == nil {
		t.Error("Expected error when JSON is invalid")
	}
}
