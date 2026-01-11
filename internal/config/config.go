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
	General GeneralConfig `mapstructure:"general"`
	Serve   ServeConfig   `mapstructure:"serve"`
}

// GeneralConfig represents the general configuration
type GeneralConfig struct {
	BaseDirectory string `mapstructure:"base_directory"`
	DateFormat    string `mapstructure:"date_format"`
	PathFormat    string `mapstructure:"path_format"`
	FileTemplate  string `mapstructure:"file_template"`
	Title         string `mapstructure:"title"`
	Sort          string `mapstructure:"sort"`
	CSS           string `mapstructure:"css"`
}

// ServeConfig represents the serve command configuration
type ServeConfig struct {
	Port int `mapstructure:"port"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if err := c.General.Validate(); err != nil {
		return fmt.Errorf("general config: %w", err)
	}

	if err := c.Serve.Validate(); err != nil {
		return fmt.Errorf("serve config: %w", err)
	}

	return nil
}

// Validate validates the general configuration
func (g *GeneralConfig) Validate() error {
	if g.BaseDirectory == "" {
		return fmt.Errorf("base_directory is required")
	}

	validSorts := map[string]bool{SortDesc: true, SortAsc: true}
	if g.Sort != "" && !validSorts[g.Sort] {
		return fmt.Errorf("invalid sort: %s (must be one of: desc, asc)", g.Sort)
	}

	return nil
}

// Validate validates the serve configuration
func (s *ServeConfig) Validate() error {
	if s.Port < 0 || s.Port > 65535 {
		return fmt.Errorf("port must be between 0 and 65535")
	}

	return nil
}

// SetDefaults sets default values for the configuration
func (c *Config) SetDefaults() {
	c.General.SetDefaults()
	c.Serve.SetDefaults()
}

// SetDefaults sets default values for the general configuration
func (g *GeneralConfig) SetDefaults() {
	if g.DateFormat == "" {
		g.DateFormat = "2006-01-02"
	}
	if g.PathFormat == "" {
		g.PathFormat = "2006-01-02.md"
	}
	if g.FileTemplate == "" {
		g.FileTemplate = "# {{ .Date }}\n"
	}
	if g.Title == "" {
		g.Title = "Journal"
	}
	if g.Sort == "" {
		g.Sort = DefaultSort
	}
}

// SetDefaults sets default values for the serve configuration
func (s *ServeConfig) SetDefaults() {
	if s.Port == 0 {
		s.Port = DefaultPort
	}
}
