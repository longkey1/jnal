package jnal

import (
	"testing"
	"time"
)

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func datesOf(entries Entries) []time.Time {
	dates := make([]time.Time, len(entries))
	for i, e := range entries {
		dates[i] = e.Date
	}
	return dates
}

func entriesFromDates(dates []time.Time) Entries {
	entries := make(Entries, len(dates))
	for i, d := range dates {
		entries[i] = Entry{Date: d}
	}
	return entries
}

func TestEntries_SortByDateDesc(t *testing.T) {
	tests := []struct {
		name  string
		input []time.Time
		want  []time.Time
	}{
		{
			name:  "unsorted entries",
			input: []time.Time{date(2024, 1, 15), date(2024, 3, 1), date(2023, 12, 31)},
			want:  []time.Time{date(2024, 3, 1), date(2024, 1, 15), date(2023, 12, 31)},
		},
		{
			name:  "already sorted descending",
			input: []time.Time{date(2024, 2, 1), date(2024, 1, 1)},
			want:  []time.Time{date(2024, 2, 1), date(2024, 1, 1)},
		},
		{
			name:  "single entry",
			input: []time.Time{date(2024, 1, 1)},
			want:  []time.Time{date(2024, 1, 1)},
		},
		{
			name:  "empty entries",
			input: []time.Time{},
			want:  []time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries := entriesFromDates(tt.input)
			entries.SortByDateDesc()
			got := datesOf(entries)
			for i := range tt.want {
				if !got[i].Equal(tt.want[i]) {
					t.Errorf("SortByDateDesc()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestEntries_SortByDateAsc(t *testing.T) {
	tests := []struct {
		name  string
		input []time.Time
		want  []time.Time
	}{
		{
			name:  "unsorted entries",
			input: []time.Time{date(2024, 1, 15), date(2023, 12, 31), date(2024, 3, 1)},
			want:  []time.Time{date(2023, 12, 31), date(2024, 1, 15), date(2024, 3, 1)},
		},
		{
			name:  "already sorted ascending",
			input: []time.Time{date(2024, 1, 1), date(2024, 2, 1)},
			want:  []time.Time{date(2024, 1, 1), date(2024, 2, 1)},
		},
		{
			name:  "single entry",
			input: []time.Time{date(2024, 1, 1)},
			want:  []time.Time{date(2024, 1, 1)},
		},
		{
			name:  "empty entries",
			input: []time.Time{},
			want:  []time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries := entriesFromDates(tt.input)
			entries.SortByDateAsc()
			got := datesOf(entries)
			for i := range tt.want {
				if !got[i].Equal(tt.want[i]) {
					t.Errorf("SortByDateAsc()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
