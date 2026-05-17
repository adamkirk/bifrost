package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adamkirk/bifrost/api/internal/app"
	"github.com/danielgtaylor/huma/v2"
)

type V1BetaEnvironmentComponentsController struct {
	environmentComponentsHandler environmentComponentsHandler
}

func (c *V1BetaEnvironmentComponentsController) RegisterRoutes(v ApiVersion, api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   fmt.Sprintf("%s.environments.components.create", string(v)),
		Method:        http.MethodPost,
		Path:          "/environments/{environment_name}/components",
		Summary:       "Create a new environment component",
		DefaultStatus: http.StatusCreated,
		Tags: []string{
			"Environment Components",
		},
	}, ErrorHandler(c.Create, http.MethodPost))
}

func NewV1BetaEnvironmentComponentsController(environmentComponentsHandler environmentComponentsHandler) *V1BetaEnvironmentComponentsController {
	return &V1BetaEnvironmentComponentsController{
		environmentComponentsHandler: environmentComponentsHandler,
	}
}

type V1BetaCreateEnvironmentComponentRequestBody struct {
	Name          string `json:"name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Unique name for the component within the environment."`
	ChartName     string `json:"chart_name" minLength:"1" doc:"Name of the Helm chart."`
	ChartVersion  string `json:"chart_version" minLength:"1" doc:"Version of the Helm chart."`
	ChartRegistry string `json:"chart_registry" minLength:"1" doc:"Registry where the Helm chart is hosted."`
}

type V1BetaCreateEnvironmentComponentRequest struct {
	EnvironmentName string `path:"environment_name" minLength:"1" pattern:"^[a-zA-Z0-9-]+$" doc:"Name of the environment."`
	Body            V1BetaCreateEnvironmentComponentRequestBody
}

func (req *V1BetaCreateEnvironmentComponentRequest) MapErrorKey(targetField string) string {
	switch targetField {
	case "EnvironmentName":
		return "path.environment_name"
	case "Name":
		return "body.name"
	case "ChartName":
		return "body.chart_name"
	case "ChartVersion":
		return "body.chart_version"
	case "ChartRegistry":
		return "body.chart_registry"
	default:
		return targetField
	}
}

type V1BetaEnvironmentComponent struct {
	ID            string `json:"id"`
	EnvironmentID string `json:"environment_id"`
	Name          string `json:"name"`
	ChartName     string `json:"chart_name"`
	ChartVersion  string `json:"chart_version"`
	ChartRegistry string `json:"chart_registry"`
}

type V1BetaEnvironmentComponentResponseBody struct {
	Data V1BetaEnvironmentComponent `json:"data"`
}

type V1BetaEnvironmentComponentResponse struct {
	Status int
	Body   V1BetaEnvironmentComponentResponseBody
}

func (c *V1BetaEnvironmentComponentsController) Create(ctx context.Context, req *V1BetaCreateEnvironmentComponentRequest) (*V1BetaEnvironmentComponentResponse, error) {
	component, err := c.environmentComponentsHandler.Create(app.CreateEnvironmentComponentDTO{
		EnvironmentName: req.EnvironmentName,
		Name:            req.Body.Name,
		ChartName:       req.Body.ChartName,
		ChartVersion:    req.Body.ChartVersion,
		ChartRegistry:   req.Body.ChartRegistry,
	})

	if err != nil {
		return nil, err
	}

	return &V1BetaEnvironmentComponentResponse{
		Status: http.StatusCreated,
		Body: V1BetaEnvironmentComponentResponseBody{
			Data: V1BetaEnvironmentComponent{
				ID:            component.ID.String(),
				EnvironmentID: component.EnvironmentID.String(),
				Name:          component.Name,
				ChartName:     component.ChartName,
				ChartVersion:  component.ChartVersion,
				ChartRegistry: component.ChartRegistry,
			},
		},
	}, nil
}
