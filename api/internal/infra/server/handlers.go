package server

import (
	"github.com/adamkirk/bifrost/api/internal/app"
	"github.com/adamkirk/bifrost/api/internal/common"
)

type environmentComponentsHandler interface {
	Create(dto app.CreateEnvironmentComponentDTO) (*common.EnvironmentComponent, error)
}

type environmentsHandler interface {
	Create(dto app.CreateEnvironmentDTO) (*common.Environment, error)
	Get(dto app.GetEnvironmentDTO) (*common.Environment, error)
	List(dto app.ListEnvironmentsDTO) (*app.ListEnvironmentsResult, error)
	Update(dto app.UpdateEnvironmentDTO) (*common.Environment, error)
	Delete(dto app.DeleteEnvironmentDTO) error
}
