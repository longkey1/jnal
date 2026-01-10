package config

import (
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				BaseDirectory:  "/home/user/journal",
				FileNameFormat: "2006-01-02.md",
				OpenCommand:    "vim {{ .File }}",
			},
			wantErr: false,
		},
		{
			name: "missing base_directory",
			config: Config{
				FileNameFormat: "2006-01-02.md",
				OpenCommand:    "vim {{ .File }}",
			},
			wantErr: true,
		},
		{
			name: "missing file_name_format",
			config: Config{
				BaseDirectory: "/home/user/journal",
				OpenCommand:   "vim {{ .File }}",
			},
			wantErr: true,
		},
		{
			name: "missing open_command",
			config: Config{
				BaseDirectory:  "/home/user/journal",
				FileNameFormat: "2006-01-02.md",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServeConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ServeConfig
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  ServeConfig{Port: 8080, Group: "month", Sort: "desc"},
			wantErr: false,
		},
		{
			name:    "empty values are valid",
			config:  ServeConfig{},
			wantErr: false,
		},
		{
			name:    "invalid port negative",
			config:  ServeConfig{Port: -1},
			wantErr: true,
		},
		{
			name:    "invalid port too high",
			config:  ServeConfig{Port: 70000},
			wantErr: true,
		},
		{
			name:    "invalid group",
			config:  ServeConfig{Group: "invalid"},
			wantErr: true,
		},
		{
			name:    "invalid sort",
			config:  ServeConfig{Sort: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServeConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_SetDefaults(t *testing.T) {
	cfg := &Config{}
	cfg.SetDefaults()

	if cfg.DateFormat != "2006-01-02" {
		t.Errorf("DateFormat = %v, want 2006-01-02", cfg.DateFormat)
	}
	if cfg.FileNameFormat != "2006-01-02.md" {
		t.Errorf("FileNameFormat = %v, want 2006-01-02.md", cfg.FileNameFormat)
	}
	if cfg.Serve.Port != DefaultPort {
		t.Errorf("Serve.Port = %v, want %v", cfg.Serve.Port, DefaultPort)
	}
	if cfg.Serve.Group != DefaultGroup {
		t.Errorf("Serve.Group = %v, want %v", cfg.Serve.Group, DefaultGroup)
	}
	if cfg.Serve.Sort != DefaultSort {
		t.Errorf("Serve.Sort = %v, want %v", cfg.Serve.Sort, DefaultSort)
	}
}
