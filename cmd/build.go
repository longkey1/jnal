package cmd

import (
	"fmt"

	"github.com/longkey1/jnal/internal/jnal"
	"github.com/longkey1/jnal/internal/server"
	"github.com/spf13/cobra"
)

func newBuildCommand(app **jnal.App) *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build static HTML files",
		Long:  `Generate static HTML files from journal entries.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := (*app).Config()
			jnl := (*app).Journal()

			builder, err := server.NewBuilder(cfg, jnl, cfg.Common.BaseDirectory)
			if err != nil {
				return fmt.Errorf("creating builder: %w", err)
			}

			if err := builder.Build(output); err != nil {
				return fmt.Errorf("building: %w", err)
			}

			fmt.Printf("Built to %s\n", output)
			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "public", "Output directory")

	return cmd
}
