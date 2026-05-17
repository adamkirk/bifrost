package app

import (
	"log/slog"

	"github.com/adamkirk/bifrost/api/internal/common"
	"github.com/google/uuid"
)

type CreateEnvironmentComponentDTO struct {
	EnvironmentName string
	Name            string
	ChartName       string
	ChartVersion    string
	ChartRegistry   string
}

func (dto CreateEnvironmentComponentDTO) Validate(repo environmentComponentsRepository, environmentID uuid.UUID) error {
	fldErrors := []common.FieldError{}

	if !common.IsValidSlug(dto.Name) {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "Name",
			Error: "must contain alphanumeric or hyphen characters only",
			Value: dto.Name,
		})
	} else {
		existing, err := repo.ByEnvironmentAndName(environmentID, dto.Name)
		if err != nil {
			return err
		}

		if existing != nil {
			fldErrors = append(fldErrors, common.FieldError{
				Key:   "Name",
				Error: "a component with this name already exists in this environment",
				Value: dto.Name,
			})
		}
	}

	if len(fldErrors) > 0 {
		return common.ValidationError{FieldErrors: fldErrors}
	}

	return nil
}

type EnvironmentComponentsHandler struct {
	l                               *slog.Logger
	environmentsRepository          environmentsRepository
	environmentComponentsRepository environmentComponentsRepository
}

func (h *EnvironmentComponentsHandler) Create(dto CreateEnvironmentComponentDTO) (*common.EnvironmentComponent, error) {
	env, err := h.environmentsRepository.ByName(dto.EnvironmentName)
	if err != nil {
		return nil, err
	}

	if env == nil {
		return nil, common.ErrNotFound{}
	}

	if err := dto.Validate(h.environmentComponentsRepository, env.ID); err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	c := &common.EnvironmentComponent{
		ID:            id,
		EnvironmentID: env.ID,
		Name:          dto.Name,
		ChartName:     dto.ChartName,
		ChartVersion:  dto.ChartVersion,
		ChartRegistry: dto.ChartRegistry,
	}

	return c, h.environmentComponentsRepository.Save(c)
}

func NewEnvironmentComponentsHandler(
	l *slog.Logger,
	environmentsRepository environmentsRepository,
	environmentComponentsRepository environmentComponentsRepository,
) *EnvironmentComponentsHandler {
	return &EnvironmentComponentsHandler{
		l:                               l,
		environmentsRepository:          environmentsRepository,
		environmentComponentsRepository: environmentComponentsRepository,
	}
}
