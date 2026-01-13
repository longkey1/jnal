package journal

import (
	"fmt"
	"os"
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
	relativePath := date.Format(j.cfg.Common.PathFormat)
	return filepath.Join(j.cfg.Common.BaseDirectory, relativePath)
}

// GetBaseDir returns the base directory path
func (j *Journal) GetBaseDir() string {
	return j.cfg.Common.BaseDirectory
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


// ListEntries returns all journal entries in the base directory
func (j *Journal) ListEntries() (Entries, error) {
	var entries Entries

	err := filepath.Walk(j.cfg.Common.BaseDirectory, func(path string, info os.FileInfo, err error) error {
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
		return nil, fmt.Errorf("walking directory %s: %w", j.cfg.Common.BaseDirectory, err)
	}

	return entries, nil
}


// buildEntryContent builds the initial content for a new entry
func (j *Journal) buildEntryContent(date time.Time) (string, error) {
	dateStr := date.Format(j.cfg.Common.DateFormat)

	return j.executeTemplate(j.cfg.New.FileTemplate, map[string]interface{}{
		"Date": dateStr,
		"Env":  getEnvMap(),
	})
}
