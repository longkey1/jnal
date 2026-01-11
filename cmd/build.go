package cmd

import (
	"fmt"

	"github.com/longkey1/jnal/internal/server"
	"github.com/spf13/cobra"
)

var buildOutput string

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build static HTML files",
	Long:  `Generate static HTML files from journal entries.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		builder, err := server.NewBuilder(cfg, jnl, cfg.General.BaseDirectory)
		if err != nil {
			return fmt.Errorf("creating builder: %w", err)
		}

		if err := builder.Build(buildOutput); err != nil {
			return fmt.Errorf("building: %w", err)
		}

		fmt.Printf("Built to %s\n", buildOutput)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVarP(&buildOutput, "output", "o", "public", "Output directory")
}
