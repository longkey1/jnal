package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/longkey1/jnal/internal/config"
	"github.com/spf13/cobra"
)

const defaultConfigTemplate = `# jnal configuration file

[common]
# base_directory = "."  # Default: current directory
date_format = "2006-01-02"
path_format = "2006-01-02.md"

[new]
file_template = "# {{ .Date }}\n"

[build]
title = "Journal"
sort = "desc"
# heading_shift = 4  # Shift heading levels in HTML output (0 to disable)
# css = "https://cdn.jsdelivr.net/npm/water.css@2/out/water.css"

[serve]
port = 8080
`

func newInitCommand() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize jnal configuration",
		Long:  `Create a default configuration file at .jnal.toml in the current directory`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := filepath.Join(".", config.DefaultConfigFileName)

			// Check if config already exists
			if _, err := os.Stat(configPath); err == nil && !force {
				return fmt.Errorf("config file already exists at %s (use --force to overwrite)", configPath)
			}

			// Write config file
			if err := os.WriteFile(configPath, []byte(defaultConfigTemplate), config.FilePermission); err != nil {
				return fmt.Errorf("writing config file: %w", err)
			}

			fmt.Printf("Created config file at %s\n", configPath)
			fmt.Printf("Edit the file to customize your journal settings.\n")

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing config file")

	return cmd
}
