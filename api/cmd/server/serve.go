package main

import (
	"log/slog"

	"github.com/adamkirk/bifrost/api/internal/config"
	"github.com/adamkirk/bifrost/api/internal/server"
	"github.com/spf13/cobra"
)

func handleServe(cmd *cobra.Command, _ []string, cfg *config.Config) error {
	var accessLogger *slog.Logger

	if cfg.Server.AccessLogsEnabled {
		h := slog.NewJSONHandler(cmd.OutOrStderr(), &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		accessLogger = slog.New(h)
		accessLogger = accessLogger.With("component", "http-server-access")
	}

	s := server.New(cfg.GetServerPort(), logger.With("component", "http-server"), accessLogger)
	return s.Start()
}
