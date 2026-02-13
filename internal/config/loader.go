package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// DefaultConfigFileName is the default configuration file name
const DefaultConfigFileName = ".jnal.toml"

// DefaultConfigPath returns the default configuration file path (current directory)
func DefaultConfigPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current directory: %w", err)
	}
	return cwd + "/" + DefaultConfigFileName, nil
}

// Load loads configuration from the specified file or default location
// If no config file exists, returns default configuration
func Load(configFile string) (*Config, error) {
	v := viper.New()

	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		// Use current directory with .jnal.toml
		v.AddConfigPath(".")
		v.SetConfigName(".jnal")
		v.SetConfigType("toml")
	}

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		// If config file not found, use default configuration
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			cfg := &Config{}
			cfg.SetDefaults()
			if err := cfg.Validate(); err != nil {
				return nil, fmt.Errorf("validating config: %w", err)
			}
			return cfg, nil
		}
		// For explicit config file, return error if not found
		if configFile != "" && os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s", configFile)
		}
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	cfg.SetDefaults()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return &cfg, nil
}

