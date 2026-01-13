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
				Common: CommonConfig{
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

func TestCommonConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  CommonConfig
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  CommonConfig{BaseDirectory: "/home/user/journal"},
			wantErr: false,
		},
		{
			name:    "missing base_directory",
			config:  CommonConfig{},
			wantErr: true,
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

	if cfg.Common.DateFormat != "2006-01-02" {
		t.Errorf("Common.DateFormat = %v, want 2006-01-02", cfg.Common.DateFormat)
	}
	if cfg.New.FileTemplate != "# {{ .Date }}\n" {
		t.Errorf("New.FileTemplate = %v, want # {{ .Date }}\\n", cfg.New.FileTemplate)
	}
	if cfg.Build.Sort != DefaultSort {
		t.Errorf("Build.Sort = %v, want %v", cfg.Build.Sort, DefaultSort)
	}
	if cfg.Serve.Port != DefaultPort {
		t.Errorf("Serve.Port = %v, want %v", cfg.Serve.Port, DefaultPort)
	}
}
