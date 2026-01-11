package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/longkey1/jnal/internal/config"
	"github.com/spf13/cobra"
)

var initForce bool

const defaultConfigTemplate = `# jnal configuration file

base_directory = "%s"
date_format = "2006-01-02"
path_format = "2006-01-02.md"
file_template = "# {{ .Date }}\n"

[serve]
port = 8080
sort = "desc"
# css = "https://cdn.jsdelivr.net/npm/water.css@2/out/water.css"
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize jnal configuration",
	Long:  `Create a default configuration file at ~/.config/jnal/config.toml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := config.DefaultConfigPath()
		if err != nil {
			return fmt.Errorf("getting config path: %w", err)
		}

		// Check if config already exists
		if _, err := os.Stat(configPath); err == nil && !initForce {
			return fmt.Errorf("config file already exists at %s (use --force to overwrite)", configPath)
		}

		// Create config directory if it doesn't exist
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, config.DirPermission); err != nil {
			return fmt.Errorf("creating config directory: %w", err)
		}

		// Get default journal directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("getting home directory: %w", err)
		}
		defaultJournalDir := filepath.Join(homeDir, "journal")

		// Write config file
		content := fmt.Sprintf(defaultConfigTemplate, defaultJournalDir)
		if err := os.WriteFile(configPath, []byte(content), config.FilePermission); err != nil {
			return fmt.Errorf("writing config file: %w", err)
		}

		fmt.Printf("Created config file at %s\n", configPath)
		fmt.Printf("Edit the file to customize your journal settings.\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Overwrite existing config file")
}
