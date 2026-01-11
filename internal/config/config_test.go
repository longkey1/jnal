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
				General: GeneralConfig{
					BaseDirectory: "/home/user/journal",
				},
			},
			wantErr: false,
		},
		{
			name:    "missing base_directory",
			config:  Config{},
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

func TestGeneralConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  GeneralConfig
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  GeneralConfig{BaseDirectory: "/home/user/journal", Sort: "desc"},
			wantErr: false,
		},
		{
			name:    "missing base_directory",
			config:  GeneralConfig{},
			wantErr: true,
		},
		{
			name:    "invalid sort",
			config:  GeneralConfig{BaseDirectory: "/home/user/journal", Sort: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneralConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
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
			config:  ServeConfig{Port: 8080},
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

	if cfg.General.DateFormat != "2006-01-02" {
		t.Errorf("General.DateFormat = %v, want 2006-01-02", cfg.General.DateFormat)
	}
	if cfg.General.Sort != DefaultSort {
		t.Errorf("General.Sort = %v, want %v", cfg.General.Sort, DefaultSort)
	}
	if cfg.Serve.Port != DefaultPort {
		t.Errorf("Serve.Port = %v, want %v", cfg.Serve.Port, DefaultPort)
	}
}
