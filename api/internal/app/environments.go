package app

import (
	"log/slog"

	"github.com/adamkirk/bifrost/api/internal/common"
	"github.com/google/uuid"
)

type CreateEnvironmentDTO struct {
	Name string
}

type EnvironmentsHandler struct {
	l                      *slog.Logger
	environmentsRepository environmentsRepository
}

func (h *EnvironmentsHandler) Create(dto CreateEnvironmentDTO) (*common.Environment, error) {
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
