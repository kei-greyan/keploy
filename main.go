// Package main is the entry point for the Keploy application.
// Keploy is an open source API testing platform that captures and replays
// API calls to generate test cases and data mocks automatically.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/keploy/keploy/v2/cmd"
	"go.uber.org/zap"
)

func main() {
	// Initialize a temporary logger for startup errors
	logger, err := zap.NewDevelopment()
	if err != nil {
		// If we can't create a logger, fall back to stderr
		os.Stderr.WriteString("failed to initialize logger: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer logger.Sync() //nolint:errcheck // best-effort sync on exit

	// Create a root context that is cancelled on OS interrupt signals.
	// This allows all components to perform graceful shutdown.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer stop()

	// Execute the root CLI command. All sub-commands (record, test, mock, etc.)
	// are registered inside the cmd package.
	if err := cmd.Execute(ctx, logger); err != nil {
		logger.Error("keploy exited with error", zap.Error(err))
		os.Exit(1)
	}
}
