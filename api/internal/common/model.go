package common

import (
	"regexp"

	"github.com/google/uuid"
)

var environmentNameRegex = regexp.MustCompile("^[a-zA-Z0-9-]+$")

type Environment struct {
	ID   uuid.UUID
	Name string
}

type EnvironmentComponent struct {
	ID            uuid.UUID
	EnvironmentID uuid.UUID
	Name          string
	ChartName     string
	ChartVersion  string
	ChartRegistry string
}

func IsValidEnvironmentName(name string) bool {
	return environmentNameRegex.MatchString(name)
}
