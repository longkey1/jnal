package journal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/longkey1/jnal/internal/config"
	"github.com/longkey1/jnal/internal/dateutil"
)

// Journal manages journal entries
type Journal struct {
	cfg *config.Config
}

// New creates a new Journal instance
func New(cfg *config.Config) *Journal {
	return &Journal{cfg: cfg}
}

// GetEntryPath returns the file path for a journal entry on the given date
func (j *Journal) GetEntryPath(date time.Time) string {
	filename := date.Format(j.cfg.FileNameFormat)
	return filepath.Join(j.cfg.BaseDirectory, filename)
}

// GetBaseDir returns the base directory path
func (j *Journal) GetBaseDir() string {
	return j.cfg.BaseDirectory
}

// EntryExists checks if a journal entry exists for the given date
func (j *Journal) EntryExists(date time.Time) bool {
	path := j.GetEntryPath(date)
	_, err := os.Stat(path)
	return err == nil
}

// CreateEntry creates a new journal entry for the given date
// Returns the file path and any error encountered
func (j *Journal) CreateEntry(date time.Time) (string, error) {
	entryPath := j.GetEntryPath(date)

	// Check if entry already exists
	if _, err := os.Stat(entryPath); err == nil {
		return entryPath, nil
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(entryPath)
	if err := os.MkdirAll(dir, config.DirPermission); err != nil {
		return "", fmt.Errorf("creating directory %s: %w", dir, err)
	}

	// Create the file
	file, err := os.OpenFile(entryPath, os.O_WRONLY|os.O_CREATE, config.FilePermission)
	if err != nil {
		return "", fmt.Errorf("creating file %s: %w", entryPath, err)
	}
	defer file.Close()

	// Write template content
	content, err := j.buildEntryContent(date)
	if err != nil {
		return "", fmt.Errorf("building entry content: %w", err)
	}

	if _, err := fmt.Fprintln(file, content); err != nil {
		return "", fmt.Errorf("writing entry content: %w", err)
	}

	return entryPath, nil
}

// OpenEntry opens a journal entry with the configured editor
func (j *Journal) OpenEntry(date time.Time) error {
	cmd, err := j.buildOpenCommand(date)
	if err != nil {
		return fmt.Errorf("building open command: %w", err)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("executing open command: %w", err)
	}

	return nil
}

// ListEntries returns all journal entries in the base directory
func (j *Journal) ListEntries() (Entries, error) {
	var entries Entries

	err := filepath.Walk(j.cfg.BaseDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Only process .md files
		if filepath.Ext(path) != ".md" {
			return nil
		}

		// Extract date from filename
		date, err := dateutil.ExtractFromFilename(info.Name())
		if err != nil {
			// Skip files without valid date in filename
			return nil
		}

		entries = append(entries, Entry{
			Path: path,
			Date: date,
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking directory %s: %w", j.cfg.BaseDirectory, err)
	}

	return entries, nil
}

// buildOpenCommand builds the command to open a journal entry
func (j *Journal) buildOpenCommand(date time.Time) (*exec.Cmd, error) {
	entryPath := j.GetEntryPath(date)
	dateStr := date.Format(j.cfg.DateFormat)

	cmdStr, err := j.executeTemplate(j.cfg.OpenCommand, map[string]interface{}{
		"BaseDir": j.cfg.BaseDirectory,
		"Date":    dateStr,
		"File":    entryPath,
		"Env":     getEnvMap(),
	})
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}

// buildEntryContent builds the initial content for a new entry
func (j *Journal) buildEntryContent(date time.Time) (string, error) {
	dateStr := date.Format(j.cfg.DateFormat)

	return j.executeTemplate(j.cfg.FileTemplate, map[string]interface{}{
		"Date": dateStr,
		"Env":  getEnvMap(),
	})
}
