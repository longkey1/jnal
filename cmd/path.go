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
	pathBase  bool
)

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show file or directory path",
	Long:  `Show the path of a journal entry file or the base directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var targetPath string

		if pathBase {
			targetPath = jnl.GetBaseDir()
		} else {
			targetDate, err := dateutil.Parse(pathDate)
			if err != nil {
				return fmt.Errorf("invalid date %q: %w", pathDate, err)
			}
			targetPath = jnl.GetEntryPath(targetDate)
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

	pathCmd.Flags().BoolVarP(&pathBase, "base", "b", false, "Show base directory path")
	pathCmd.Flags().BoolVarP(&pathCheck, "check", "c", false, "Check if the path exists")
	pathCmd.Flags().StringVarP(&pathDate, "date", "d", dateutil.Format(dateutil.Today()), "Target date (format: yyyy-mm-dd)")
}
