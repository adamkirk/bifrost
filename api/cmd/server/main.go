package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/adamkirk/bifrost/api/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
	Level: slog.LevelError,
}))

func buildLogger(cfg *config.Config) (*slog.Logger, error) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(cfg.Logging.Level)); err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", cfg.Logging.Level, err)
	}

	opts := &slog.HandlerOptions{Level: level}

	switch cfg.Logging.Format {
	case "json":
		return slog.New(slog.NewJSONHandler(os.Stderr, opts)), nil
	case "text":
		return slog.New(slog.NewTextHandler(os.Stderr, opts)), nil
	default:
		return nil, fmt.Errorf("invalid log format %q: expected json or text", cfg.Logging.Format)
	}
}

type runEHandlerFunc func(cmd *cobra.Command, args []string) error
type runEHandlerFuncWithConfig func(cmd *cobra.Command, args []string, cfg *config.Config) error
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
	Run:   errorHandlerWrapper(withConfig(handleServe), 1),
}

func bindEnvs(v *viper.Viper, prefix string, t reflect.Type) {
	for field := range t.Fields() {
		key := field.Tag.Get("mapstructure")
		if key == "" {
			key = strings.ToLower(field.Name)
		}
		if prefix != "" {
			key = prefix + "." + key
		}
		if field.Type.Kind() == reflect.Struct {
			bindEnvs(v, key, field.Type)
		} else {
			_ = v.BindEnv(key)
		}
	}
}

func withConfig(f runEHandlerFuncWithConfig) runEHandlerFunc {
	return func(cmd *cobra.Command, args []string) error {
		v := viper.New()
		v.SetConfigFile("bifrost.server.yml")
		v.SetEnvPrefix("BIFROST_SERVER")
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		bindEnvs(v, "", reflect.TypeFor[config.Config]())

		_ = v.BindPFlag("logging.level", cmd.Root().PersistentFlags().Lookup("log-level"))
		_ = v.BindPFlag("logging.format", cmd.Root().PersistentFlags().Lookup("log-format"))

		if err := v.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}

		cfg := config.Default()
		if err := v.Unmarshal(cfg); err != nil {
			return err
		}

		var err error
		logger, err = buildLogger(cfg)

		if err != nil {
			return err
		}

		return f(cmd, args, cfg)
	}
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
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error).")
	rootCmd.PersistentFlags().String("log-format", "json", "Log format (json, text).")

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
