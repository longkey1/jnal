package cmd

import (
	"fmt"

	"github.com/longkey1/jnal/internal/jnal"
	"github.com/spf13/cobra"
)

// NewRootCommand creates the root command with all subcommands
func NewRootCommand() *cobra.Command {
	var (
		cfgFile string
		app     *jnal.App
	)

	cmd := &cobra.Command{
		Use:   "jnal",
		Short: "A simple CLI tool for daily journaling in Markdown",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip config loading for commands that don't need it
			if shouldSkipConfig(cmd) {
				return nil
			}

			var err error
			app, err = jnal.NewApp(cfgFile)
			if err != nil {
				return fmt.Errorf("initializing app: %w", err)
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is .jnal.toml in current directory)")

	// Add subcommands
	cmd.AddCommand(newNewCommand(&app))
	cmd.AddCommand(newBuildCommand(&app))
	cmd.AddCommand(newServeCommand(&app))
	cmd.AddCommand(newPathCommand(&app))
	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newVersionCommand())

	return cmd
}

// shouldSkipConfig returns true if the command doesn't need config loading
func shouldSkipConfig(cmd *cobra.Command) bool {
	skipCommands := map[string]bool{
		"init":       true,
		"version":    true,
		"help":       true,
		"completion": true,
	}
	return skipCommands[cmd.Name()]
}
