package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("jnal %s (Built on %s from Git SHA %s)\n", version, date, commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// SetVersionInfo sets version information
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
}
