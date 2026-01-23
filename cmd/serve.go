package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/longkey1/jnal/internal/config"
	"github.com/longkey1/jnal/internal/jnal"
	"github.com/longkey1/jnal/internal/server"
	"github.com/spf13/cobra"
)

func newServeCommand(app **jnal.App) *cobra.Command {
	var (
		port       int
		sort       string
		liveReload bool
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start a local preview server",
		Long: `Start a local HTTP server to preview journal entries.
The server watches for file changes and automatically reloads content.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := (*app).Config()
			jnl := (*app).Journal()

			// Override config with command line flags
			if cmd.Flags().Changed("port") {
				cfg.Serve.Port = port
			}
			if cmd.Flags().Changed("sort") {
				cfg.Build.Sort = sort
			}

			// Validate config
			if err := cfg.Validate(); err != nil {
				return fmt.Errorf("invalid config: %w", err)
			}

			srv, err := server.New(cfg, jnl, cfg.Common.BaseDirectory, liveReload)
			if err != nil {
				return fmt.Errorf("creating server: %w", err)
			}

			// Setup context with signal handling
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigCh
				fmt.Println("\nShutting down...")
				cancel()
			}()

			return srv.Start(ctx)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", config.DefaultPort, "Port to listen on")
	cmd.Flags().StringVarP(&sort, "sort", "s", config.DefaultSort, "Sort order: desc (newest first), asc (oldest first)")
	cmd.Flags().BoolVarP(&liveReload, "live-reload", "l", false, "Enable live reload on file changes")

	return cmd
}
