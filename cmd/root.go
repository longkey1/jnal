package cmd

import (
	"fmt"
	"os"

	"github.com/longkey1/jnal/internal/config"
	"github.com/longkey1/jnal/internal/journal"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
	jnl     *journal.Journal
)

var rootCmd = &cobra.Command{
	Use:   "jnal",
	Short: "Text file based journal command",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for commands that don't need it
		if cmd.Name() == "init" || cmd.Name() == "help" || cmd.Name() == "completion" {
			return nil
		}

		var err error
		cfg, err = config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		jnl = journal.New(cfg)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/jnal/config.toml)")
}

// SetVersionInfo sets version information for the root command
func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

// GetConfig returns the loaded configuration
func GetConfig() *config.Config {
	return cfg
}

// GetJournal returns the journal instance
func GetJournal() *journal.Journal {
	return jnl
}
