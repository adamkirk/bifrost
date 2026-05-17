package app

import (
	"log/slog"

	"github.com/adamkirk/bifrost/api/internal/common"
	"github.com/google/uuid"
)

type CreateEnvironmentDTO struct {
	Name string
}

func (dto CreateEnvironmentDTO) Validate(repo environmentsRepository) error {
	found, err := repo.ByName(dto.Name)

	if err != nil {
		return err
	}

	if found != nil {
		return common.ValidationError{
			FieldErrors: []common.FieldError{
				{
					Key: "Name",
					Errors: common.Violations{
						&common.ConflictViolation{
							BaseViolation: common.BaseViolation{
								Error: "an environment with this name already exists.",
							},
						},
					},
				},
			},
		}
	}

	return nil
}

type EnvironmentsHandler struct {
	l                      *slog.Logger
	environmentsRepository environmentsRepository
}

func (h *EnvironmentsHandler) Create(dto CreateEnvironmentDTO) (*common.Environment, error) {
	if err := dto.Validate(h.environmentsRepository); err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()

	if err != nil {
		return nil, err
	}

	env := &common.Environment{
		ID:   id,
		Name: dto.Name,
	}

	return env, h.environmentsRepository.Create(env)
}

func NewEnvironmentsHandler(l *slog.Logger, environmentsRepository environmentsRepository) *EnvironmentsHandler {
	return &EnvironmentsHandler{
		l:                      l,
		environmentsRepository: environmentsRepository,
	}
}
