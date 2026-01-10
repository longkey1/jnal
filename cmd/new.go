package cmd

import (
	"fmt"

	"github.com/longkey1/jnal/internal/dateutil"
	"github.com/spf13/cobra"
)

var newDate string

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create and open a journal entry",
	Long: `Create a new journal entry for the specified date (or today if not specified).
If an entry already exists for that date, it will be opened.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse the date
		targetDate, err := dateutil.Parse(newDate)
		if err != nil {
			return fmt.Errorf("invalid date %q: %w", newDate, err)
		}

		// Create entry if it doesn't exist
		entryPath, err := jnl.CreateEntry(targetDate)
		if err != nil {
			return fmt.Errorf("creating entry: %w", err)
		}

		fmt.Printf("Opening %s\n", entryPath)

		// Open the entry with the configured editor
		if err := jnl.OpenEntry(targetDate); err != nil {
			return fmt.Errorf("opening entry: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&newDate, "date", "d", dateutil.Format(dateutil.Today()), "Date for the journal entry (format: yyyy-mm-dd)")
}
