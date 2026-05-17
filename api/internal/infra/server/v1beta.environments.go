package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adamkirk/bifrost/api/internal/app"
	"github.com/danielgtaylor/huma/v2"
)

type V1BetaChartReference struct {
	Registry string `json:"registry" minLength:"1"`
	Name     string `json:"name" minLength:"1"`
	Version  string `json:"version" minLength:"1"`
}

type V1BetaCreateEnvironmentRequestBody struct {
	Name  string               `json:"name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Unique name for the environment. Must contain only alphanumeric characters and hyphens."`
	Chart V1BetaChartReference `json:"chart"`
}

type V1BetaCreateEnvironmentRequest struct {
	Body V1BetaCreateEnvironmentRequestBody
}

func (req *V1BetaCreateEnvironmentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "Name":
		return "name"
	default:
		return targetField
	}
}

type V1BetaCreateEnvironmentResponseBody struct {
	V1BetaCreateEnvironmentRequestBody
}

type V1BetaCreateEnvironmentResponse struct {
	Body V1BetaCreateEnvironmentResponseBody
}

type V1BetaEnvironmentsController struct {
	environmentsHandler environmentsHandler
}

func (c *V1BetaEnvironmentsController) RegisterRoutes(v ApiVersion, api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.create", string(v)),
		Method:        http.MethodPost,
		Path:          "/environments",
		Summary:       "Create a new environment",
		DefaultStatus: http.StatusCreated,
		Tags: []string{
			"Environments",
		},
		Metadata: map[string]any{
			OptDisableAllDefaultResponses:   true,
			OptDisableDefaultAuthentication: true,
		},
	}, ErrorHandler(c.Create))
}

func NewV1BetaEnvironmentsController(environmentsHandler environmentsHandler) *V1BetaEnvironmentsController {
	return &V1BetaEnvironmentsController{
		environmentsHandler: environmentsHandler,
	}
}

func (c *V1BetaEnvironmentsController) Create(ctx context.Context, req *V1BetaCreateEnvironmentRequest) (*V1BetaCreateEnvironmentResponse, error) {
	env, err := c.environmentsHandler.Create(app.CreateEnvironmentDTO{
		Name: req.Body.Name,
	})

	if err != nil {
		return nil, err
	}

	return &V1BetaCreateEnvironmentResponse{
		Body: V1BetaCreateEnvironmentResponseBody{
			V1BetaCreateEnvironmentRequestBody: V1BetaCreateEnvironmentRequestBody{
				Name: env.Name,
			},
		},
	}, nil
}
