package main

import (
	"errors"
	"fmt"
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

		defaults := config.Default()
		v.SetDefault("server.port", defaults.Server.Port)

		if err := v.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}

		cfg := &config.Config{}
		if err := v.Unmarshal(cfg); err != nil {
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
