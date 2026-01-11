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
	Short: "A simple CLI tool for daily journaling in Markdown",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for commands that don't need it
		if cmd.Name() == "init" || cmd.Name() == "help" || cmd.Name() == "completion" || cmd.Name() == "version" {
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

// GetConfig returns the loaded configuration
func GetConfig() *config.Config {
	return cfg
}

// GetJournal returns the journal instance
func GetJournal() *journal.Journal {
	return jnl
}
