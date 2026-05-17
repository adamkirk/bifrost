package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adamkirk/bifrost/api/internal/app"
	"github.com/danielgtaylor/huma/v2"
)

type V1BetaEnvironmentComponentDeploymentsController struct {
	deploymentsHandler deploymentsHandler
}

func (c *V1BetaEnvironmentComponentDeploymentsController) RegisterRoutes(v ApiVersion, api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.components.deployments.create", string(v)),
		Method:        http.MethodPost,
		Path:          "/environments/{environment_name}/components/{component_name}/deployments",
		Summary:       "Create a new deployment for an environment component",
		DefaultStatus: http.StatusCreated,
		Tags: []string{
			"Environment Component Deployments",
		},
	}, ErrorHandler(c.Create, http.MethodPost))
}

func NewV1BetaEnvironmentComponentDeploymentsController(deploymentsHandler deploymentsHandler) *V1BetaEnvironmentComponentDeploymentsController {
	return &V1BetaEnvironmentComponentDeploymentsController{
		deploymentsHandler: deploymentsHandler,
	}
}

type V1BetaCreateEnvironmentComponentDeploymentRequest struct {
	EnvironmentName string `path:"environment_name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the environment."`
	ComponentName   string `path:"component_name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the component."`
}

func (req *V1BetaCreateEnvironmentComponentDeploymentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "EnvironmentName":
		return "path.environment_name"
	case "ComponentName":
		return "path.component_name"
	default:
		return targetField
	}
}

type V1BetaEnvironmentComponentDeployment struct {
	ID                       string    `json:"id"`
	Status                   string    `json:"status"`
	CreatedAt                time.Time `json:"created_at"`
	EnvironmentName          string    `json:"environment_name"`
	EnvironmentComponentName string    `json:"environment_component_name"`
}

type V1BetaCreateEnvironmentComponentDeploymentResponseBody struct {
	Data V1BetaEnvironmentComponentDeployment `json:"data"`
}

type V1BetaCreateEnvironmentComponentDeploymentResponse struct {
	Status int
	Body   V1BetaCreateEnvironmentComponentDeploymentResponseBody
}

func (c *V1BetaEnvironmentComponentDeploymentsController) Create(ctx context.Context, req *V1BetaCreateEnvironmentComponentDeploymentRequest) (*V1BetaCreateEnvironmentComponentDeploymentResponse, error) {
	deployment, err := c.deploymentsHandler.Create(app.CreateDeploymentDTO{
		EnvironmentName: req.EnvironmentName,
		ComponentName:   req.ComponentName,
	})

	if err != nil {
		return nil, err
	}

	return &V1BetaCreateEnvironmentComponentDeploymentResponse{
		Status: http.StatusCreated,
		Body: V1BetaCreateEnvironmentComponentDeploymentResponseBody{
			Data: V1BetaEnvironmentComponentDeployment{
				ID:                       deployment.ID.String(),
				Status:                   string(deployment.Status),
				CreatedAt:                deployment.CreatedAt,
				EnvironmentName:          req.EnvironmentName,
				EnvironmentComponentName: req.ComponentName,
			},
		},
	}, nil
}
