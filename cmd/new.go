package cmd

import (
	"fmt"

	"github.com/longkey1/jnal/internal/jnal"
	"github.com/longkey1/jnal/internal/util"
	"github.com/spf13/cobra"
)

func newNewCommand(app **jnal.App) *cobra.Command {
	var date string

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a journal entry",
		Long:  `Create a new journal entry for the specified date (or today if not specified).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse the date
			targetDate, err := util.Parse(date)
			if err != nil {
				return fmt.Errorf("invalid date %q: %w", date, err)
			}

			// Create entry if it doesn't exist
			entryPath, err := (*app).Journal().CreateEntry(targetDate)
			if err != nil {
				return fmt.Errorf("creating entry: %w", err)
			}

			fmt.Println(entryPath)

			return nil
		},
	}

	cmd.Flags().StringVarP(&date, "date", "d",
		util.Format(util.Today()),
		"Date for the journal entry (format: yyyy-mm-dd)")

	return cmd
}
