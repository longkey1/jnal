package config

import (
	"fmt"
	"os"
)

// Permission constants
const (
	DirPermission  os.FileMode = 0755
	FilePermission os.FileMode = 0644
)

// Default values
const (
	DefaultPort = 8080
	DefaultSort = "desc"
)

// Sort options
const (
	SortDesc = "desc"
	SortAsc  = "asc"
)

// Config represents the application configuration
type Config struct {
	BaseDirectory string      `mapstructure:"base_directory"`
	DateFormat    string      `mapstructure:"date_format"`
	PathFormat    string      `mapstructure:"path_format"`
	FileTemplate  string      `mapstructure:"file_template"`
	Serve         ServeConfig `mapstructure:"serve"`
}

// ServeConfig represents the serve command configuration
type ServeConfig struct {
	Port int    `mapstructure:"port"`
	Sort string `mapstructure:"sort"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.BaseDirectory == "" {
		return fmt.Errorf("base_directory is required")
	}

	// Validate serve config
	if err := c.Serve.Validate(); err != nil {
		return fmt.Errorf("serve config: %w", err)
	}

	return nil
}

// Validate validates the serve configuration
func (s *ServeConfig) Validate() error {
	if s.Port < 0 || s.Port > 65535 {
		return fmt.Errorf("port must be between 0 and 65535")
	}

	validSorts := map[string]bool{SortDesc: true, SortAsc: true}
	if s.Sort != "" && !validSorts[s.Sort] {
		return fmt.Errorf("invalid sort: %s (must be one of: desc, asc)", s.Sort)
	}

	return nil
}

// SetDefaults sets default values for the configuration
func (c *Config) SetDefaults() {
	if c.DateFormat == "" {
		c.DateFormat = "2006-01-02"
	}
	if c.PathFormat == "" {
		c.PathFormat = "2006-01-02.md"
	}
	if c.FileTemplate == "" {
		c.FileTemplate = "# {{ .Date }}\n"
	}
	c.Serve.SetDefaults()
}

// SetDefaults sets default values for the serve configuration
func (s *ServeConfig) SetDefaults() {
	if s.Port == 0 {
		s.Port = DefaultPort
	}
	if s.Sort == "" {
		s.Sort = DefaultSort
	}
}
