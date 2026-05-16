package main

import (
	"log/slog"
	"os"

	"github.com/adamkirk/bifrost/api/internal/config"
	"github.com/adamkirk/bifrost/api/internal/server"
	"github.com/spf13/cobra"
)

func handleServe(_ *cobra.Command, _ []string, cfg *config.Config) error {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	s := server.New(cfg.GetServerPort(), logger)
	return s.Start()
}
