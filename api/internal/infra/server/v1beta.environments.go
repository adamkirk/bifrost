package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adamkirk/bifrost/api/internal/app"
	"github.com/danielgtaylor/huma/v2"
)

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
		// Metadata: map[string]any{
		// 	OptDisableAllDefaultResponses:   false,
		// 	OptDisableDefaultAuthentication: false,
		// },
	}, ErrorHandler(c.Create, http.MethodPost))

	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.get", string(v)),
		Method:        http.MethodGet,
		Path:          "/environments/{name}",
		Summary:       "Get an environment",
		DefaultStatus: http.StatusOK,
		Tags: []string{
			"Environments",
		},
		// Metadata: map[string]any{
		// 	OptDisableAllDefaultResponses:   false,
		// 	OptDisableDefaultAuthentication: false,
		// },
	}, ErrorHandler(c.Get, http.MethodGet))
}

func NewV1BetaEnvironmentsController(environmentsHandler environmentsHandler) *V1BetaEnvironmentsController {
	return &V1BetaEnvironmentsController{
		environmentsHandler: environmentsHandler,
	}
}

type V1BetaCreateEnvironmentRequestBody struct {
	Name string `json:"name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Unique name for the environment. Must contain only alphanumeric characters and hyphens."`
}

type V1BetaCreateEnvironmentRequest struct {
	Body V1BetaCreateEnvironmentRequestBody
}

func (req *V1BetaCreateEnvironmentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "Name":
		return "body.name"
	default:
		return targetField
	}
}

type V1BetaEnvironment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type V1BetaEnvironmentResponseBody struct {
	Data V1BetaEnvironment `json:"data"`
}

type V1BetaEnvironmentResponse struct {
	Status int
	Body   V1BetaEnvironmentResponseBody
}

func (c *V1BetaEnvironmentsController) Create(ctx context.Context, req *V1BetaCreateEnvironmentRequest) (*V1BetaEnvironmentResponse, error) {
	env, err := c.environmentsHandler.Create(app.CreateEnvironmentDTO{
		Name: req.Body.Name,
	})

	if err != nil {
		return nil, err
	}

	return &V1BetaEnvironmentResponse{
		Status: http.StatusCreated,
		Body: V1BetaEnvironmentResponseBody{
			Data: V1BetaEnvironment{
				ID:   env.ID.String(),
				Name: env.Name,
			},
		},
	}, nil
}

type V1BetaGetEnvironmentRequest struct {
	Name string `path:"name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Unique name for the environment. Must contain only alphanumeric characters and hyphens."`
}

func (req *V1BetaGetEnvironmentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "Name":
		return "path.name"
	default:
		return targetField
	}
}

func (c *V1BetaEnvironmentsController) Get(ctx context.Context, req *V1BetaGetEnvironmentRequest) (*V1BetaEnvironmentResponse, error) {
	env, err := c.environmentsHandler.Get(app.GetEnvironmentDTO{
		Name: req.Name,
	})

	if err != nil {
		return nil, err
	}

	if env == nil {
		return nil, huma.Error404NotFound("no environment with this name exists")
	}

	return &V1BetaEnvironmentResponse{
		Status: http.StatusOK,
		Body: V1BetaEnvironmentResponseBody{
			Data: V1BetaEnvironment{
				ID:   env.ID.String(),
				Name: env.Name,
			},
		},
	}, nil
}
