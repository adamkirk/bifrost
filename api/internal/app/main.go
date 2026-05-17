package app

import "github.com/adamkirk/bifrost/api/internal/common"

type environmentsRepository interface {
	Create(env *common.Environment) error
	ByName(name string) (*common.Environment, error)
	List(limit, offset int) ([]*common.Environment, error)
	Count() (int, error)
}
