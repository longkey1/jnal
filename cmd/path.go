package cmd

import (
	"fmt"
	"os"

	"github.com/longkey1/jnal/internal/dateutil"
	"github.com/spf13/cobra"
)

var (
	pathCheck bool
	pathDate  string
	pathType  string
)

const (
	PathTypeFile = "file"
	PathTypeBase = "base"
)

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show file or directory path",
	Long:  `Show the path of a journal entry file or the base directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDate, err := dateutil.Parse(pathDate)
		if err != nil {
			return fmt.Errorf("invalid date %q: %w", pathDate, err)
		}

		var targetPath string
		switch pathType {
		case PathTypeBase:
			targetPath = jnl.GetBaseDir()
		case PathTypeFile:
			targetPath = jnl.GetEntryPath(targetDate)
		default:
			return fmt.Errorf("invalid type %q: must be %q or %q", pathType, PathTypeFile, PathTypeBase)
		}

		if pathCheck {
			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				return fmt.Errorf("path does not exist: %s", targetPath)
			}
		}

		fmt.Println(targetPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)

	pathCmd.Flags().BoolVarP(&pathCheck, "check", "c", false, "Check if the path exists")
	pathCmd.Flags().StringVarP(&pathDate, "date", "d", dateutil.Format(dateutil.Today()), "Target date (format: yyyy-mm-dd)")
	pathCmd.Flags().StringVarP(&pathType, "type", "t", PathTypeFile, fmt.Sprintf("Path type [%s, %s]", PathTypeFile, PathTypeBase))
}
