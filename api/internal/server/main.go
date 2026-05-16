package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	port int
	echo *echo.Echo
	api  huma.API
}

type healthCheckOutput struct {
	Body struct {
		Status string `json:"status"`
	}
}

func New(port int, logger *slog.Logger) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logger.ErrorContext(c.Request().Context(), "panic recovered",
				slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
				slog.String("error", err.Error()),
				slog.String("stack", string(stack)),
			)
			return nil
		},
	}))
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogRequestID: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.InfoContext(c.Request().Context(), "access log",
				slog.String("component", "http-server"),
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
				slog.String("remote_ip", v.RemoteIP),
				slog.String("request_id", v.RequestID),
			)
			return nil
		},
	}))

	config := huma.DefaultConfig("Bifrost API", "v1beta1")
	api := humaecho.New(e, config)

	s := &Server{
		port: port,
		echo: e,
		api:  api,
	}

	huma.Register(api, huma.Operation{
		OperationID: "healthcheck",
		Method:      http.MethodGet,
		Path:        "/_/healthz",
		Summary:     "Health check",
	}, func(_ context.Context, _ *struct{}) (*healthCheckOutput, error) {
		resp := &healthCheckOutput{}
		resp.Body.Status = "ok"

		return resp, nil
	})

	return s
}

func (s *Server) Start() error {
	return s.echo.Start(fmt.Sprintf(":%d", s.port))
}
