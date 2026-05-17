package app

import "github.com/adamkirk/bifrost/api/internal/common"

type environmentsRepository interface {
	Create(env *common.Environment) error
}
