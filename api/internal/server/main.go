package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
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

func New(port int) *Server {
	e := echo.New()
	e.HideBanner = true

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
