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

	fldErrors := []common.FieldError{}

	if found != nil {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "Name",
			Error: "an environment with this name already exists.",
			Value: dto.Name,
		})
	}

	if !common.IsValidEnvironmentName(dto.Name) {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "Name",
			Error: "must contain alphanumeric or hyphen characters only",
			Value: dto.Name,
		})

	}
	if len(fldErrors) > 0 {
		return common.ValidationError{
			FieldErrors: fldErrors,
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

type GetEnvironmentDTO struct {
	Name string
}

func (dto GetEnvironmentDTO) Validate() error {
	fldErrors := []common.FieldError{}

	if !common.IsValidEnvironmentName(dto.Name) {
		fldErrors = append(fldErrors, common.FieldError{
			Key:   "Name",
			Error: "must contain alphanumeric or hyphen characters only",
			Value: dto.Name,
		})

	}
	if len(fldErrors) > 0 {
		return common.ValidationError{
			FieldErrors: fldErrors,
		}
	}

	return nil
}

func (h *EnvironmentsHandler) Get(dto GetEnvironmentDTO) (*common.Environment, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	return h.environmentsRepository.ByName(dto.Name)
}

func NewEnvironmentsHandler(l *slog.Logger, environmentsRepository environmentsRepository) *EnvironmentsHandler {
	return &EnvironmentsHandler{
		l:                      l,
		environmentsRepository: environmentsRepository,
	}
}
