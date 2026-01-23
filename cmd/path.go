package cmd

import (
	"fmt"
	"os"

	"github.com/longkey1/jnal/internal/jnal"
	"github.com/longkey1/jnal/internal/util"
	"github.com/spf13/cobra"
)

func newPathCommand(app **jnal.App) *cobra.Command {
	var (
		check bool
		date  string
		base  bool
	)

	cmd := &cobra.Command{
		Use:   "path",
		Short: "Show file or directory path",
		Long:  `Show the path of a journal entry file or the base directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var targetPath string

			if base {
				targetPath = (*app).Journal().GetBaseDir()
			} else {
				targetDate, err := util.Parse(date)
				if err != nil {
					return fmt.Errorf("invalid date %q: %w", date, err)
				}
				targetPath = (*app).Journal().GetEntryPath(targetDate)
			}

			if check {
				if _, err := os.Stat(targetPath); os.IsNotExist(err) {
					return fmt.Errorf("path does not exist: %s", targetPath)
				}
			}

			fmt.Println(targetPath)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&base, "base", "b", false, "Show base directory path")
	cmd.Flags().BoolVarP(&check, "check", "c", false, "Check if the path exists")
	cmd.Flags().StringVarP(&date, "date", "d",
		util.Format(util.Today()),
		"Target date (format: yyyy-mm-dd)")

	return cmd
}
