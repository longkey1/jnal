package jnal

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/longkey1/jnal/internal/config"
)

func newTestConfig(baseDir string) *config.Config {
	cfg := &config.Config{
		Common: config.CommonConfig{
			BaseDirectory: baseDir,
		},
	}
	cfg.SetDefaults()
	return cfg
}

func TestJournal_GetEntryPath(t *testing.T) {
	tests := []struct {
		name       string
		baseDir    string
		pathFormat string
		date       time.Time
		want       string
	}{
		{
			name:       "default flat format",
			baseDir:    "/home/user/journal",
			pathFormat: "2006-01-02.md",
			date:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			want:       "/home/user/journal/2024-01-15.md",
		},
		{
			name:       "nested year/month format",
			baseDir:    "/home/user/journal",
			pathFormat: "2006/01/2006-01-02.md",
			date:       time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			want:       "/home/user/journal/2024/12/2024-12-31.md",
		},
		{
			name:       "single digit day is zero padded",
			baseDir:    "/journal",
			pathFormat: "2006-01-02.md",
			date:       time.Date(2023, 3, 5, 0, 0, 0, 0, time.UTC),
			want:       "/journal/2023-03-05.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newTestConfig(tt.baseDir)
			cfg.Common.PathFormat = tt.pathFormat
			j := NewJournal(cfg)

			got := j.GetEntryPath(tt.date)
			if got != tt.want {
				t.Errorf("GetEntryPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJournal_GetBaseDir(t *testing.T) {
	cfg := newTestConfig("/home/user/journal")
	j := NewJournal(cfg)

	if got := j.GetBaseDir(); got != "/home/user/journal" {
		t.Errorf("GetBaseDir() = %v, want /home/user/journal", got)
	}
}

func TestJournal_EntryExists(t *testing.T) {
	baseDir := t.TempDir()
	cfg := newTestConfig(baseDir)
	j := NewJournal(cfg)

	existing := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	if err := os.WriteFile(filepath.Join(baseDir, "2024-01-15.md"), []byte("content"), 0644); err != nil {
		t.Fatalf("writing test file: %v", err)
	}

	tests := []struct {
		name string
		date time.Time
		want bool
	}{
		{
			name: "existing entry",
			date: existing,
			want: true,
		},
		{
			name: "missing entry",
			date: time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := j.EntryExists(tt.date); got != tt.want {
				t.Errorf("EntryExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJournal_CreateEntry(t *testing.T) {
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		pathFormat   string
		fileTemplate string
		env          map[string]string
		date         time.Time
		wantContent  string
		wantErr      bool
	}{
		{
			name:         "template with date",
			pathFormat:   "2006-01-02.md",
			fileTemplate: "# {{ .Date }}",
			date:         testDate,
			wantContent:  "# 2024-01-15\n",
		},
		{
			name:         "empty template creates near-empty file",
			pathFormat:   "2006-01-02.md",
			fileTemplate: "",
			date:         testDate,
			wantContent:  "\n",
		},
		{
			name:         "template with environment variable",
			pathFormat:   "2006-01-02.md",
			fileTemplate: "# {{ .Date }} by {{ .Env.JNAL_TEST_AUTHOR }}",
			env:          map[string]string{"JNAL_TEST_AUTHOR": "tester"},
			date:         testDate,
			wantContent:  "# 2024-01-15 by tester\n",
		},
		{
			name:         "nested path format creates directories",
			pathFormat:   "2006/01/2006-01-02.md",
			fileTemplate: "# {{ .Date }}",
			date:         testDate,
			wantContent:  "# 2024-01-15\n",
		},
		{
			name:         "invalid template returns error",
			pathFormat:   "2006-01-02.md",
			fileTemplate: "{{ .Date",
			date:         testDate,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			baseDir := t.TempDir()
			cfg := newTestConfig(baseDir)
			cfg.Common.PathFormat = tt.pathFormat
			cfg.New.FileTemplate = tt.fileTemplate
			j := NewJournal(cfg)

			path, err := j.CreateEntry(tt.date)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			wantPath := j.GetEntryPath(tt.date)
			if path != wantPath {
				t.Errorf("CreateEntry() path = %v, want %v", path, wantPath)
			}

			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("reading created entry: %v", err)
			}
			if string(data) != tt.wantContent {
				t.Errorf("CreateEntry() content = %q, want %q", string(data), tt.wantContent)
			}
		})
	}
}

func TestJournal_CreateEntry_ExistingEntryIsNotOverwritten(t *testing.T) {
	baseDir := t.TempDir()
	cfg := newTestConfig(baseDir)
	cfg.New.FileTemplate = "# {{ .Date }}"
	j := NewJournal(cfg)

	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	existingPath := j.GetEntryPath(testDate)
	existingContent := "already written"
	if err := os.WriteFile(existingPath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("writing existing entry: %v", err)
	}

	path, err := j.CreateEntry(testDate)
	if err != nil {
		t.Fatalf("CreateEntry() error = %v", err)
	}
	if path != existingPath {
		t.Errorf("CreateEntry() path = %v, want %v", path, existingPath)
	}

	data, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatalf("reading entry: %v", err)
	}
	if string(data) != existingContent {
		t.Errorf("CreateEntry() overwrote existing entry: content = %q, want %q", string(data), existingContent)
	}
}

func TestJournal_ListEntries(t *testing.T) {
	baseDir := t.TempDir()
	cfg := newTestConfig(baseDir)
	j := NewJournal(cfg)

	// Entries in nested and flat locations
	subDir := filepath.Join(baseDir, "2024", "01")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("creating sub directory: %v", err)
	}

	files := map[string]string{
		filepath.Join(baseDir, "2023-12-31.md"):   "old entry",
		filepath.Join(subDir, "2024-01-15.md"):    "nested entry",
		filepath.Join(baseDir, "notes.md"):        "markdown without date, skipped",
		filepath.Join(baseDir, "2024-01-16.txt"):  "not markdown, skipped",
		filepath.Join(baseDir, "2024-01-17-x.md"): "entry with suffix",
	}
	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("writing test file %s: %v", path, err)
		}
	}

	entries, err := j.ListEntries()
	if err != nil {
		t.Fatalf("ListEntries() error = %v", err)
	}

	wantDates := map[string]bool{
		"2023-12-31": false,
		"2024-01-15": false,
		"2024-01-17": false,
	}
	if len(entries) != len(wantDates) {
		t.Fatalf("ListEntries() returned %d entries, want %d", len(entries), len(wantDates))
	}
	for _, e := range entries {
		key := e.Date.Format("2006-01-02")
		if _, ok := wantDates[key]; !ok {
			t.Errorf("ListEntries() unexpected entry date %s (path %s)", key, e.Path)
			continue
		}
		wantDates[key] = true
	}
	for d, seen := range wantDates {
		if !seen {
			t.Errorf("ListEntries() missing entry for date %s", d)
		}
	}
}

func TestJournal_ListEntries_MissingBaseDirectory(t *testing.T) {
	cfg := newTestConfig(filepath.Join(t.TempDir(), "does-not-exist"))
	j := NewJournal(cfg)

	if _, err := j.ListEntries(); err == nil {
		t.Error("ListEntries() error = nil, want error for missing base directory")
	}
}

func TestGetEnvMap(t *testing.T) {
	t.Setenv("JNAL_TEST_ENV_KEY", "env-value")

	envMap := getEnvMap()
	if got := envMap["JNAL_TEST_ENV_KEY"]; got != "env-value" {
		t.Errorf("getEnvMap()[JNAL_TEST_ENV_KEY] = %q, want %q", got, "env-value")
	}
}
