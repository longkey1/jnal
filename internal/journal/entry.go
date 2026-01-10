package journal

import (
	"time"
)

// Entry represents a journal entry
type Entry struct {
	Path    string
	Date    time.Time
	Content string
}

// Entries is a slice of Entry
type Entries []Entry

// SortByDateDesc sorts entries by date in descending order (newest first)
func (e Entries) SortByDateDesc() {
	for i := 0; i < len(e)-1; i++ {
		for j := i + 1; j < len(e); j++ {
			if e[j].Date.After(e[i].Date) {
				e[i], e[j] = e[j], e[i]
			}
		}
	}
}

// SortByDateAsc sorts entries by date in ascending order (oldest first)
func (e Entries) SortByDateAsc() {
	for i := 0; i < len(e)-1; i++ {
		for j := i + 1; j < len(e); j++ {
			if e[j].Date.Before(e[i].Date) {
				e[i], e[j] = e[j], e[i]
			}
		}
	}
}
