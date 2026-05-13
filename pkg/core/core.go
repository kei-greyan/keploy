// Package core provides the central orchestration logic for Keploy,
// including test recording, replaying, and result reporting.
package core

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// Mode represents the operational mode of Keploy.
type Mode string

const (
	// ModeRecord captures outgoing network calls and stores them as test cases.
	ModeRecord Mode = "record"
	// ModeTest replays stored test cases and validates application responses.
	ModeTest Mode = "test"
	// ModeOff disables Keploy instrumentation entirely.
	ModeOff Mode = "off"
)

// Config holds the runtime configuration for a Keploy session.
type Config struct {
	// Mode determines whether Keploy is recording or replaying tests.
	Mode Mode
	// AppID is the unique identifier for the application under test.
	AppID string
	// TestSetsPath is the directory where test sets are stored or read from.
	TestSetsPath string
	// Delay is the time in seconds to wait before starting test replay.
	Delay uint64
	// APITimeout is the maximum duration (seconds) for a single API call during replay.
	APITimeout uint64
}

// Core is the main orchestrator that manages the lifecycle of recording
// and replaying test sessions.
type Core struct {
	mu     sync.Mutex
	logger *zap.Logger
	cfg    Config
	running bool
}

// New creates and returns a new Core instance with the provided configuration.
func New(logger *zap.Logger, cfg Config) (*Core, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger must not be nil")
	}
	if cfg.AppID == "" {
		return nil, fmt.Errorf("AppID must not be empty")
	}
	if cfg.Mode != ModeRecord && cfg.Mode != ModeTest && cfg.Mode != ModeOff {
		return nil, fmt.Errorf("invalid mode %q: must be one of record, test, off", cfg.Mode)
	}
	return &Core{
		logger: logger,
		cfg:    cfg,
	}, nil
}

// Start begins the Keploy session based on the configured mode.
// It is safe to call Start only once; subsequent calls return an error.
func (c *Core) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return fmt.Errorf("core is already running")
	}

	c.logger.Info("Starting Keploy core",
		zap.String("mode", string(c.cfg.Mode)),
		zap.String("appID", c.cfg.AppID),
	)

	switch c.cfg.Mode {
	case ModeRecord:
		if err := c.startRecord(ctx); err != nil {
			return fmt.Errorf("failed to start record mode: %w", err)
		}
	case ModeTest:
		if err := c.startTest(ctx); err != nil {
			return fmt.Errorf("failed to start test mode: %w", err)
		}
	case ModeOff:
		c.logger.Info("Keploy is in off mode; no instrumentation active")
	}

	c.running = true
	return nil
}

// Stop gracefully shuts down the active Keploy session.
func (c *Core) Stop(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return nil
	}

	c.logger.Info("Stopping Keploy core", zap.String("mode", string(c.cfg.Mode)))
	c.running = false
	return nil
}

// startRecord initialises the recording pipeline.
func (c *Core) startRecord(ctx context.Context) error {
	c.logger.Info("Record mode initialised",
		zap.String("testSetsPath", c.cfg.TestSetsPath),
	)
	// TODO: wire up eBPF / proxy hooks for traffic capture.
	return nil
}

// startTest initialises the test-replay pipeline.
func (c *Core) startTest(ctx context.Context) error {
	c.logger.Info("Test mode initialised",
		zap.String("testSetsPath", c.cfg.TestSetsPath),
		zap.Uint64("delay", c.cfg.Delay),
		zap.Uint64("apiTimeout", c.cfg.APITimeout),
	)
	// TODO: load test sets from disk and begin replay.
	return nil
}
