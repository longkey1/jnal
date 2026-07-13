package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing config file: %v", err)
	}
	return path
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		content string
		check   func(t *testing.T, cfg *Config)
		wantErr bool
	}{
		{
			name: "full config",
			content: `
[common]
base_directory = "/home/user/journal"
date_format = "2006/01/02"
path_format = "2006/01/2006-01-02.md"

[new]
file_template = "# {{ .Date }}"

[build]
title = "My Journal"
sort = "asc"
heading_shift = 2
hard_wraps = false
linkify = false
link_target_blank = false

[serve]
port = 9090
`,
			check: func(t *testing.T, cfg *Config) {
				if cfg.Common.BaseDirectory != "/home/user/journal" {
					t.Errorf("BaseDirectory = %v, want /home/user/journal", cfg.Common.BaseDirectory)
				}
				if cfg.Common.DateFormat != "2006/01/02" {
					t.Errorf("DateFormat = %v, want 2006/01/02", cfg.Common.DateFormat)
				}
				if cfg.Common.PathFormat != "2006/01/2006-01-02.md" {
					t.Errorf("PathFormat = %v, want 2006/01/2006-01-02.md", cfg.Common.PathFormat)
				}
				if cfg.New.FileTemplate != "# {{ .Date }}" {
					t.Errorf("FileTemplate = %v, want # {{ .Date }}", cfg.New.FileTemplate)
				}
				if cfg.Build.Title != "My Journal" {
					t.Errorf("Title = %v, want My Journal", cfg.Build.Title)
				}
				if cfg.Build.Sort != SortAsc {
					t.Errorf("Sort = %v, want %v", cfg.Build.Sort, SortAsc)
				}
				if got := cfg.Build.GetHeadingShift(); got != 2 {
					t.Errorf("GetHeadingShift() = %v, want 2", got)
				}
				if cfg.Build.GetHardWraps() {
					t.Error("GetHardWraps() = true, want false")
				}
				if cfg.Build.GetLinkify() {
					t.Error("GetLinkify() = true, want false")
				}
				if cfg.Build.GetLinkTargetBlank() {
					t.Error("GetLinkTargetBlank() = true, want false")
				}
				if cfg.Serve.Port != 9090 {
					t.Errorf("Port = %v, want 9090", cfg.Serve.Port)
				}
			},
		},
		{
			name: "minimal config applies defaults",
			content: `
[common]
base_directory = "/home/user/journal"
`,
			check: func(t *testing.T, cfg *Config) {
				if cfg.Common.DateFormat != "2006-01-02" {
					t.Errorf("DateFormat = %v, want 2006-01-02", cfg.Common.DateFormat)
				}
				if cfg.Common.PathFormat != "2006-01-02.md" {
					t.Errorf("PathFormat = %v, want 2006-01-02.md", cfg.Common.PathFormat)
				}
				if cfg.Build.Title != "Journal" {
					t.Errorf("Title = %v, want Journal", cfg.Build.Title)
				}
				if cfg.Build.Sort != DefaultSort {
					t.Errorf("Sort = %v, want %v", cfg.Build.Sort, DefaultSort)
				}
				if got := cfg.Build.GetHeadingShift(); got != DefaultHeadingShift {
					t.Errorf("GetHeadingShift() = %v, want %v", got, DefaultHeadingShift)
				}
				if !cfg.Build.GetHardWraps() {
					t.Error("GetHardWraps() = false, want true")
				}
				if !cfg.Build.GetLinkify() {
					t.Error("GetLinkify() = false, want true")
				}
				if !cfg.Build.GetLinkTargetBlank() {
					t.Error("GetLinkTargetBlank() = false, want true")
				}
				if cfg.Serve.Port != DefaultPort {
					t.Errorf("Port = %v, want %v", cfg.Serve.Port, DefaultPort)
				}
			},
		},
		{
			name: "missing base_directory fails validation",
			content: `
[build]
title = "My Journal"
`,
			wantErr: true,
		},
		{
			name: "invalid sort fails validation",
			content: `
[common]
base_directory = "/home/user/journal"

[build]
sort = "sideways"
`,
			wantErr: true,
		},
		{
			name:    "malformed toml",
			content: `[common`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeConfigFile(t, tt.content)

			cfg, err := Load(path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			tt.check(t, cfg)
		})
	}
}

func TestLoad_MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.toml")
	if _, err := Load(path); err == nil {
		t.Error("Load() error = nil, want error for missing config file")
	}
}

func TestDefaultConfigPath(t *testing.T) {
	path, err := DefaultConfigPath()
	if err != nil {
		t.Fatalf("DefaultConfigPath() error = %v", err)
	}

	want := filepath.Join(".config", "jnal", "config.toml")
	if !strings.HasSuffix(path, want) {
		t.Errorf("DefaultConfigPath() = %v, want suffix %v", path, want)
	}
	if !filepath.IsAbs(path) {
		t.Errorf("DefaultConfigPath() = %v, want absolute path", path)
	}
}
