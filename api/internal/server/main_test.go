package server_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/adamkirk/bifrost/api/internal/server"
	"github.com/danielgtaylor/huma/v2"
)

type DummyController struct{}

func (c *DummyController) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "dummy.healthz",
		Method:        http.MethodGet,
		Path:          "/healthz",
		Summary:       "Check if the app is started up",
		DefaultStatus: http.StatusNoContent,
		Tags: []string{
			"Healthz",
		},
		Metadata: map[string]any{},
	}, c.Healthz)
}

type HealthzRequest struct{}

func (c *DummyController) Healthz(ctx context.Context, req *HealthzRequest) (*server.NoContent, error) {
	return &server.NoContent{
		Status: http.StatusNoContent,
	}, nil
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func waitForServer(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url) //nolint:noctx
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return fmt.Errorf("server did not become ready within %s", timeout)
}

func TestAPI(t *testing.T) {
	port := freePort(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	s := server.New(port, logger, logger, &DummyController{})
	go s.Start() //nolint:errcheck

	base := fmt.Sprintf("http://localhost:%d", port)
	if err := waitForServer(base+"/api/v1beta/healthz", 2*time.Second); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get(base + "/api/v1beta/healthz") //nolint:noctx
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
