package jnal

import (
	"fmt"
	"os"

	"github.com/longkey1/jnal/internal/config"
)

// App manages the application state
type App struct {
	config  *config.Config
	journal *Journal
}

// NewApp creates a new App instance with the given config path
func NewApp(configPath string) (*App, error) {
	// JNAL_CONFIG environment variable processing
	if configPath == "" {
		configPath = os.Getenv("JNAL_CONFIG")
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	jnl := NewJournal(cfg)

	return &App{
		config:  cfg,
		journal: jnl,
	}, nil
}

// Config returns the application configuration
func (a *App) Config() *config.Config {
	return a.config
}

// Journal returns the journal instance
func (a *App) Journal() *Journal {
	return a.journal
}
