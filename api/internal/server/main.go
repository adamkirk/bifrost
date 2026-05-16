package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const OptDisableNotFound = "DisableNotFound"
const OptDisableAllDefaultResponses = "DisableAllDefaults"
const OptDisableDefaultAuthentication = "DisableAuthentication"

type ApiVersion string

const ApiVersionV1Beta ApiVersion = "v1beta"

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

type Controller interface {
	RegisterRoutes(g huma.API)
}

var opsWithoutBodies = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodDelete,

	// Probably don't need this one, but leaving for good measure
	http.MethodTrace,
}

func setupHumaMiddlewares(api huma.API) {
	api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		ctx.SetHeader("X-Operation-Id", ctx.Operation().OperationID)
		next(ctx)
	})
}

func setupEchoMiddlewares(e *echo.Echo, logger *slog.Logger, accessLogger *slog.Logger) {
	e.Pre(middleware.RemoveTrailingSlash())
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

	if accessLogger != nil {
		e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogMethod:    true,
			LogURI:       true,
			LogStatus:    true,
			LogLatency:   true,
			LogRemoteIP:  true,
			LogRequestID: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				level := slog.LevelInfo
				switch {
				case v.Status >= 500:
					level = slog.LevelError
				case v.Status >= 400:
					level = slog.LevelWarn
				}

				accessLogger.LogAttrs(c.Request().Context(), level, "access log",
					slog.String("operation_id", c.Response().Header().Get("X-Operation-Id")),
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
	}
}

func addValidationErrorResponse(op *huma.Operation) {
	validationStatus := strconv.Itoa(http.StatusUnprocessableEntity)

	if slices.Contains(opsWithoutBodies, op.Method) {
		validationStatus = strconv.Itoa(http.StatusBadRequest)
	}

	if _, ok := op.Responses[validationStatus]; ok {
		return
	}

	op.Responses[validationStatus] = &huma.Response{
		Description: "validation error",
		Content: map[string]*huma.MediaType{
			"application/problem+json": {
				Schema: &huma.Schema{
					Ref: "#/components/schemas/ErrorModel",
				},
			},
		},
	}
}
func addInternalErrorResponse(op *huma.Operation) {
	internalError := strconv.Itoa(http.StatusInternalServerError)

	if _, ok := op.Responses[internalError]; ok {
		return
	}

	op.Responses[internalError] = &huma.Response{
		Description: "Internal server error",
		Content: map[string]*huma.MediaType{
			"application/problem+json": {
				Schema: &huma.Schema{
					Ref: "#/components/schemas/ErrorModel",
				},
			},
		},
	}
}

func addNotFoundResponse(op *huma.Operation) {
	var notFoundEnabled = true

	if v, ok := op.Metadata[OptDisableNotFound]; ok {
		if optAsBool, ok := v.(bool); ok {
			notFoundEnabled = !optAsBool
		}
	}

	if !notFoundEnabled {
		return
	}

	notFound := strconv.Itoa(http.StatusNotFound)

	if _, ok := op.Responses[notFound]; ok {
		return
	}

	op.Responses[notFound] = &huma.Response{
		Description: "Resource Not Found",
		Content: map[string]*huma.MediaType{
			"application/problem+json": {
				Schema: &huma.Schema{
					Ref: "#/components/schemas/ErrorModel",
				},
			},
		},
	}
}

func configureDefaultResponses(api *huma.OpenAPI, op *huma.Operation) {

	if _, ok := op.Responses["default"]; ok {
		// Remove the default as it's an error, but has no status code
		// Maybe there's another way to turn it off
		op.Responses["default"] = nil
	}

	if v, ok := op.Metadata[OptDisableAllDefaultResponses]; ok {
		if optAsBool, ok := v.(bool); ok && optAsBool {
			return
		}
	}

	addValidationErrorResponse(op)
	addInternalErrorResponse(op)
	addNotFoundResponse(op)
}

func setupHumaHooks(api huma.API) {
	api.OpenAPI().OnAddOperation = append(
		api.OpenAPI().OnAddOperation,
		// Note, this should come before default responses, as we may want to use
		// the security requirements to configure extra responses based on whether
		// authentication is required.
		configureDefaultResponses,
	)
}

// TODO: controllers isn't quite right, as they're all attached to the v1beta api
// Fine for now, and CBA to try and sort, but will need sorting when a new api version
// exists.
func New(port int, logger *slog.Logger, accessLogger *slog.Logger, controllers ...Controller) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	apiBase := fmt.Sprintf("/api/%s", ApiVersionV1Beta)
	api := e.Group(apiBase)
	apiCfg := huma.DefaultConfig("Bifrost API", "v1beta1")

	hg := humaecho.NewWithGroup(e, api, apiCfg)

	setupHumaMiddlewares(hg)
	setupEchoMiddlewares(e, logger, accessLogger)
	setupHumaHooks(hg)

	s := &Server{
		port: port,
		echo: e,
		api:  hg,
	}

	for _, c := range controllers {
		c.RegisterRoutes(hg)
	}

	return s
}

func (s *Server) Start() error {
	return s.echo.Start(fmt.Sprintf(":%d", s.port))
}
