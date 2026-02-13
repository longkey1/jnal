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
			name: "valid config with base_directory",
			config: Config{
				Common: CommonConfig{
					BaseDirectory: "/home/user/journal",
				},
			},
			wantErr: false,
		},
		{
			name:    "empty config is valid (base_directory defaults to current directory)",
			config:  Config{},
			wantErr: false,
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

func TestCommonConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  CommonConfig
		wantErr bool
	}{
		{
			name:    "valid config with base_directory",
			config:  CommonConfig{BaseDirectory: "/home/user/journal"},
			wantErr: false,
		},
		{
			name:    "empty base_directory is valid (defaults to current directory)",
			config:  CommonConfig{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CommonConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  BuildConfig
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  BuildConfig{Sort: "desc"},
			wantErr: false,
		},
		{
			name:    "empty values are valid",
			config:  BuildConfig{},
			wantErr: false,
		},
		{
			name:    "invalid sort",
			config:  BuildConfig{Sort: "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
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

	if cfg.Common.BaseDirectory != "." {
		t.Errorf("Common.BaseDirectory = %v, want .", cfg.Common.BaseDirectory)
	}
	if cfg.Common.DateFormat != "2006-01-02" {
		t.Errorf("Common.DateFormat = %v, want 2006-01-02", cfg.Common.DateFormat)
	}
	// FileTemplate has no default - empty means create an empty file
	if cfg.New.FileTemplate != "" {
		t.Errorf("New.FileTemplate = %v, want empty string", cfg.New.FileTemplate)
	}
	if cfg.Build.Sort != DefaultSort {
		t.Errorf("Build.Sort = %v, want %v", cfg.Build.Sort, DefaultSort)
	}
	if cfg.Serve.Port != DefaultPort {
		t.Errorf("Serve.Port = %v, want %v", cfg.Serve.Port, DefaultPort)
	}
}
