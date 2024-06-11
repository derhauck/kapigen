package stages

import (
	"fmt"

	"kapigen.kateops.com/internal/logger"
)

type Stage int

const (
	LINT Stage = iota
	INIT
	BUILD
	TEST
	RELEASE
	DYNAMIC
	TRIGGER
	FINAL
)

var values = map[Stage]string{
	LINT:    "lint",
	INIT:    "init",
	BUILD:   "build",
	TEST:    "test",
	RELEASE: "release",
	DYNAMIC: "dynamic",
	TRIGGER: "trigger",
	FINAL:   "final",
}

func NewStage() Stage {
	return DYNAMIC
}

func (s Stage) String() string {
	if value, ok := values[s]; ok {
		return value
	}
	logger.Error(fmt.Sprintf("Stage not found for id: '%d'", s))
	return values[DYNAMIC]
}

func FromString(value string) (Stage, error) {
	for k, v := range values {
		if v == value {
			return k, nil
		}
	}
	return DYNAMIC, fmt.Errorf("stage not found for value: '%s'", value)
}

func GetAllStages() []string {
	var stages []string
	for _, value := range values {
		stages = append(stages, value)
	}
	return stages
}
