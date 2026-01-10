package dateutil

import (
	"fmt"
	"regexp"
	"time"
)

// ISO8601Date is the standard date format used for parsing
const ISO8601Date = "2006-01-02"

// DatePattern matches yyyy-mm-dd format in strings
var DatePattern = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

// Parse parses a date string in ISO8601 format (yyyy-mm-dd)
func Parse(dateStr string) (time.Time, error) {
	t, err := time.Parse(ISO8601Date, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format %q: expected yyyy-mm-dd: %w", dateStr, err)
	}
	return t, nil
}

// ExtractFromFilename extracts the first yyyy-mm-dd pattern from a filename
func ExtractFromFilename(filename string) (time.Time, error) {
	match := DatePattern.FindString(filename)
	if match == "" {
		return time.Time{}, fmt.Errorf("no date pattern found in filename %q", filename)
	}
	return Parse(match)
}

// Today returns today's date with time set to midnight
func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// Format formats a time.Time to ISO8601 date string
func Format(t time.Time) string {
	return t.Format(ISO8601Date)
}
