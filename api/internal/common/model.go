package common

import (
	"regexp"

	"github.com/google/uuid"
)

var slugRegex = regexp.MustCompile("^[a-zA-Z0-9-]+$")

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

func IsValidSlug(name string) bool {
	return slugRegex.MatchString(name)
}
