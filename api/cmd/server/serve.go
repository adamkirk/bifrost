package main

import (
	"github.com/adamkirk/bifrost/api/internal/server"
	"github.com/spf13/cobra"
)

func handleServe(_ *cobra.Command, _ []string) error {
	s := server.New(8080)
	return s.Start()
}
