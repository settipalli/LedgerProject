package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"ledgerproject/config"
	"os"
	"strings"
)

var log *zap.Logger

// NewLogger creates a new zap logger instance based on the provided config
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	isDevelopment := strings.ToLower(os.Getenv("APP_ENV")) != "production"
	if err := Init(isDevelopment); err != nil {
		return nil, err
	}
	return Get(), nil
}

func Init(isDevelopment bool) error {
	var cfg zap.Config
	if isDevelopment {
		// Colored console output with debug level
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		// JSON format with info level
		cfg = zap.NewProductionConfig()
	}

	var err error
	log, err = cfg.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(log)
	return nil
}

func Get() *zap.Logger {
	return log
}

// Sync flushes any buffered log entries
func Sync() error {
	return log.Sync()
}
