package main

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"ledgerproject/api"
	"ledgerproject/config"
	"ledgerproject/ledger"
	"ledgerproject/logger"
	"ledgerproject/services"
	"os"
	"strings"
	"time"
)

// getEnvironmentOption returns the appropriate fx.Option based on the environment
func getEnvironmentOption() fx.Option {
	env := strings.ToLower(os.Getenv("APP_ENV"))

	switch env {
	case "test":
		return config.Test()
	case "prod", "production":
		return config.Production()
	default:
		// Default to development if environment is not specified or unknown
		return config.Development()
	}
}

func main() {
	app := fx.New(
		// Use environment-specific config based on APP_ENV
		getEnvironmentOption(),

		// Provide core dependencies
		fx.Provide(
			logger.NewLogger,
			services.NewCurrencyValidator,
			ledger.NewLedger,
			api.NewServer,
		),

		// Register lifecycle hooks
		fx.Invoke(registerHooks),

		// Configure fx logging
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		logger.Get().Fatal("Failed to start application", zap.Error(err))
	}

	defer func() {
		if err := app.Stop(context.Background()); err != nil {
			logger.Get().Fatal("Failed to stop application", zap.Error(err))
		}
	}()

	<-app.Done()
}

func registerHooks(lc fx.Lifecycle, server *api.Server, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting server")
			// Start server in a goroutine
			go func() {
				if err := server.Start(); err != nil {
					log.Error("Failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping server")

			// Graceful shutdown with timeout
			shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Error("Server shutdown error", zap.Error(err))
				return err
			}

			if err := logger.Sync(); err != nil {
				log.Error("Failed to sync logger", zap.Error(err))
			}

			return nil
		},
	})
}
