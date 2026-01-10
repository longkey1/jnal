package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/longkey1/jnal/internal/config"
	"github.com/longkey1/jnal/internal/server"
	"github.com/spf13/cobra"
)

var (
	servePort  int
	serveGroup string
	serveSort  string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a local preview server",
	Long: `Start a local HTTP server to preview journal entries.
The server watches for file changes and automatically reloads content.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Override config with command line flags
		serveCfg := cfg.Serve
		if cmd.Flags().Changed("port") {
			serveCfg.Port = servePort
		}
		if cmd.Flags().Changed("group") {
			serveCfg.Group = serveGroup
		}
		if cmd.Flags().Changed("sort") {
			serveCfg.Sort = serveSort
		}

		// Validate serve config
		if err := serveCfg.Validate(); err != nil {
			return fmt.Errorf("invalid serve config: %w", err)
		}

		srv, err := server.New(&serveCfg, jnl, cfg.BaseDirectory)
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

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntVarP(&servePort, "port", "p", config.DefaultPort, "Port to listen on")
	serveCmd.Flags().StringVarP(&serveGroup, "group", "g", config.DefaultGroup, "Group entries by: none, year, month, week")
	serveCmd.Flags().StringVarP(&serveSort, "sort", "s", config.DefaultSort, "Sort order: desc (newest first), asc (oldest first)")
}
