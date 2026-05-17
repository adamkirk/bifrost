package server

import (
	"github.com/adamkirk/bifrost/api/internal/app"
	"github.com/adamkirk/bifrost/api/internal/common"
)

type environmentsHandler interface {
	Create(dto app.CreateEnvironmentDTO) (*common.Environment, error)
}
