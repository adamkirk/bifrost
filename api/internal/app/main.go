package app

import "github.com/adamkirk/bifrost/api/internal/common"

type environmentsRepository interface {
	ByName(name string) (*common.Environment, error)
	List(limit, offset int) ([]*common.Environment, error)
	Count() (int, error)
	Save(env *common.Environment) error
	Delete(env *common.Environment) error
}
