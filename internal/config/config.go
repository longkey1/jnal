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
	DefaultPort         = 8080
	DefaultSort         = "desc"
	DefaultHeadingShift = 4
)

// Sort options
const (
	SortDesc = "desc"
	SortAsc  = "asc"
)

// Config represents the application configuration
type Config struct {
	Common CommonConfig `mapstructure:"common"`
	New    NewConfig    `mapstructure:"new"`
	Build  BuildConfig  `mapstructure:"build"`
	Serve  ServeConfig  `mapstructure:"serve"`
}

// CommonConfig represents common configuration shared across commands
type CommonConfig struct {
	BaseDirectory string `mapstructure:"base_directory"`
	DateFormat    string `mapstructure:"date_format"`
	PathFormat    string `mapstructure:"path_format"`
}

// NewConfig represents the new command configuration
type NewConfig struct {
	FileTemplate string `mapstructure:"file_template"`
}

// BuildConfig represents the build command configuration (HTML content generation)
type BuildConfig struct {
	Title        string `mapstructure:"title"`
	Sort         string `mapstructure:"sort"`
	CSS          string `mapstructure:"css"`
	HeadingShift *int   `mapstructure:"heading_shift"`
	HardWraps    *bool  `mapstructure:"hard_wraps"`
}

// ServeConfig represents the serve command configuration (content delivery)
type ServeConfig struct {
	Port int `mapstructure:"port"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if err := c.Common.Validate(); err != nil {
		return fmt.Errorf("common config: %w", err)
	}

	if err := c.Build.Validate(); err != nil {
		return fmt.Errorf("build config: %w", err)
	}

	if err := c.Serve.Validate(); err != nil {
		return fmt.Errorf("serve config: %w", err)
	}

	return nil
}

// Validate validates the common configuration
func (c *CommonConfig) Validate() error {
	if c.BaseDirectory == "" {
		return fmt.Errorf("base_directory is required")
	}

	return nil
}

// Validate validates the build configuration
func (b *BuildConfig) Validate() error {
	validSorts := map[string]bool{SortDesc: true, SortAsc: true}
	if b.Sort != "" && !validSorts[b.Sort] {
		return fmt.Errorf("invalid sort: %s (must be one of: desc, asc)", b.Sort)
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
	c.Common.SetDefaults()
	c.New.SetDefaults()
	c.Build.SetDefaults()
	c.Serve.SetDefaults()
}

// SetDefaults sets default values for the common configuration
func (c *CommonConfig) SetDefaults() {
	if c.DateFormat == "" {
		c.DateFormat = "2006-01-02"
	}
	if c.PathFormat == "" {
		c.PathFormat = "2006-01-02.md"
	}
}

// SetDefaults sets default values for the new configuration
func (n *NewConfig) SetDefaults() {
	if n.FileTemplate == "" {
		n.FileTemplate = "# {{ .Date }}\n"
	}
}

// SetDefaults sets default values for the build configuration
func (b *BuildConfig) SetDefaults() {
	if b.Title == "" {
		b.Title = "Journal"
	}
	if b.Sort == "" {
		b.Sort = DefaultSort
	}
	if b.HeadingShift == nil {
		defaultShift := DefaultHeadingShift
		b.HeadingShift = &defaultShift
	}
	if b.HardWraps == nil {
		defaultHardWraps := true
		b.HardWraps = &defaultHardWraps
	}
}

// GetHeadingShift returns the heading shift value (0 means disabled)
func (b *BuildConfig) GetHeadingShift() int {
	if b.HeadingShift == nil {
		return DefaultHeadingShift
	}
	return *b.HeadingShift
}

// GetHardWraps returns the hard wraps setting (default: true)
func (b *BuildConfig) GetHardWraps() bool {
	if b.HardWraps == nil {
		return true
	}
	return *b.HardWraps
}

// SetDefaults sets default values for the serve configuration
func (s *ServeConfig) SetDefaults() {
	if s.Port == 0 {
		s.Port = DefaultPort
	}
}
