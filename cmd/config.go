package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/longkey1/jnal/internal/jnal"
	"github.com/spf13/cobra"
)

func newConfigCommand(app **jnal.App, cfgFile *string) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Display current configuration",
		Long:  "Display the current configuration values loaded from the config file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := (*app).Config()

			// Display config file path
			configFilePath := *cfgFile
			if configFilePath == "" {
				configFilePath = ".jnal.toml"
			}
			fmt.Printf("ConfigFile: %s\n", configFilePath)

			// Common section
			fmt.Println("[common]")
			baseDir := cfg.Common.BaseDirectory
			if !filepath.IsAbs(baseDir) {
				absPath, err := filepath.Abs(baseDir)
				if err == nil {
					baseDir = absPath
				}
			}
			fmt.Printf("  base_directory: %s\n", baseDir)
			fmt.Printf("  date_format: %s\n", cfg.Common.DateFormat)
			fmt.Printf("  path_format: %s\n", cfg.Common.PathFormat)

			// New section
			fmt.Println("[new]")
			fmt.Printf("  file_template: %s\n", cfg.New.FileTemplate)

			// Build section
			fmt.Println("[build]")
			fmt.Printf("  title: %s\n", cfg.Build.Title)
			fmt.Printf("  sort: %s\n", cfg.Build.Sort)
			if cfg.Build.CSS != "" {
				fmt.Printf("  css: %s\n", cfg.Build.CSS)
			}
			fmt.Printf("  heading_shift: %d\n", cfg.Build.GetHeadingShift())
			fmt.Printf("  hard_wraps: %t\n", cfg.Build.GetHardWraps())
			fmt.Printf("  linkify: %t\n", cfg.Build.GetLinkify())
			fmt.Printf("  link_target_blank: %t\n", cfg.Build.GetLinkTargetBlank())

			// Serve section
			fmt.Println("[serve]")
			fmt.Printf("  port: %d\n", cfg.Serve.Port)

			return nil
		},
	}
}
