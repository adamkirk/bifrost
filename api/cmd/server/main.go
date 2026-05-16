package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type runEHandlerFunc func(cmd *cobra.Command, args []string) error
type runHandlerFunc func(cmd *cobra.Command, args []string)

var rootCmd = &cobra.Command{
	Use:   "bifrost-server",
	Short: "Bifrost API server.",
	RunE:  handleGroup,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version you're using.",
	Run:   errorHandlerWrapper(handleVersion, 1),
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server for the backend components.",
	Run:   errorHandlerWrapper(handleServe, 1),
}

func errorHandlerWrapper(f runEHandlerFunc, errorExitCode int) runHandlerFunc {
	return func(cmd *cobra.Command, args []string) {
		if err := f(cmd, args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(errorExitCode)
		}
	}
}

func handleGroup(cmd *cobra.Command, _ []string) error {
	fmt.Println("blah")
	return cmd.Help()
}

func init() {
	versionCmd.Flags().Bool("short", false, "Show only the version, excluding commit and date information.")
	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(serveCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
