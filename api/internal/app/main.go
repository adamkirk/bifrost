package app

import (
	"github.com/adamkirk/bifrost/api/internal/common"
	"github.com/google/uuid"
)

type environmentsRepository interface {
	ByName(name string) (*common.Environment, error)
	List(limit, offset int) ([]*common.Environment, error)
	Count() (int, error)
	Save(env *common.Environment) error
	Delete(env *common.Environment) error
}

type environmentComponentsRepository interface {
	ByEnvironmentAndName(environmentID uuid.UUID, name string) (*common.EnvironmentComponent, error)
	CountByEnvironment(environmentID uuid.UUID) (int, error)
	ListByEnvironment(environmentID uuid.UUID, limit, offset int) ([]*common.EnvironmentComponent, error)
	Save(c *common.EnvironmentComponent) error
	Delete(c *common.EnvironmentComponent) error
}

type deploymentsRepository interface {
	Save(d *common.Deployment) error
}
