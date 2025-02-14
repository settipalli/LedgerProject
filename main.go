package main

import (
	"context"
	"go.uber.org/fx"
	"ledgerproject/api"
	"ledgerproject/config"
	"ledgerproject/ledger"
	"ledgerproject/services"
	"os"
	"strings"
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
			services.NewCurrencyValidator,
			ledger.NewLedger,
			api.NewServer,
		),

		// Register lifecycle hooks
		fx.Invoke(registerHooks),
	)

	app.Run()
}

func registerHooks(lc fx.Lifecycle, server *api.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start server in a goroutine
			go server.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Add cleanup logic if needed
			return nil
		},
	})
}
