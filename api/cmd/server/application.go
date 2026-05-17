package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/adamkirk/bifrost/api/internal/config"
	"github.com/adamkirk/bifrost/api/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Application struct {
	// These should always be set, and safe to rely upon during construction of
	// any other services.
	stdout io.Writer
	stderr io.Writer
	logger *slog.Logger
	cfg    *config.Config

	server *server.Server
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

func (a *Application) loadConfig(cmd *cobra.Command) error {
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

	a.cfg = cfg

	return nil
}

func (a *Application) setupLogger() error {
	var level slog.Level
	if err := level.UnmarshalText([]byte(a.cfg.Logging.Level)); err != nil {
		return fmt.Errorf("invalid log level %q: %w", a.cfg.Logging.Level, err)
	}

	opts := &slog.HandlerOptions{Level: level}

	var l *slog.Logger

	switch a.cfg.Logging.Format {
	case "json":
		l = slog.New(slog.NewJSONHandler(os.Stderr, opts))
	case "text":
		l = slog.New(slog.NewTextHandler(os.Stderr, opts))
	default:
		return fmt.Errorf("invalid log format %q: expected json or text", a.cfg.Logging.Format)
	}

	a.logger = l

	return nil
}

func (a *Application) GetServer() *server.Server {
	var accessLogger *slog.Logger

	if a.cfg.Server.AccessLogsEnabled {
		h := slog.NewJSONHandler(a.stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		accessLogger = slog.New(h)
		accessLogger = accessLogger.With("component", "http-server-access")
	}

	return server.New(
		a.cfg.GetServerPort(),
		a.logger.With("component", "http-server"),
		server.WithAccessLogger(accessLogger),
		server.WithApiVersionGroup(
			server.ApiVersionGroup{
				Version:     server.ApiVersionV1Beta,
				Controllers: a.GetV1BetaControllers(),
			},
		),
	)
}

func (a *Application) GetV1BetaControllers() []server.Controller {
	return []server.Controller{
		server.NewProbesController(),
		server.NewV1BetaEnvironmentsController(),
	}
}

func NewApplication(cmd *cobra.Command) (*Application, error) {
	app := &Application{
		stderr: cmd.OutOrStderr(),
		stdout: cmd.OutOrStdout(),
	}

	if err := app.loadConfig(cmd); err != nil {
		return nil, err
	}

	if err := app.setupLogger(); err != nil {
		return nil, err
	}

	return app, nil
}
