package version

import (
	"fmt"
	"runtime"
)

// These variables are set at build time using ldflags
var (
	Version   = "dev"
	CommitSHA = "unknown"
	BuildTime = "unknown"
	GoVersion = runtime.Version()
)

// Info returns version information as a string
func Info() string {
	return fmt.Sprintf("jnal %s (Built on %s from Git SHA %s)",
		Version, BuildTime, CommitSHA)
}

// Short returns a short version string
func Short() string {
	return Version
}
