package version

import (
	"strings"
	"testing"
)

func TestInfo(t *testing.T) {
	info := Info()

	for _, want := range []string{"jnal", Version, CommitSHA, BuildTime} {
		if !strings.Contains(info, want) {
			t.Errorf("Info() = %q, want it to contain %q", info, want)
		}
	}
}

func TestShort(t *testing.T) {
	if got := Short(); got != Version {
		t.Errorf("Short() = %v, want %v", got, Version)
	}
}
