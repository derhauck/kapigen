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
)

var values = map[Stage]string{
	LINT:    "lint",
	INIT:    "init",
	BUILD:   "build",
	TEST:    "test",
	RELEASE: "release",
	DYNAMIC: "dynamic",
	TRIGGER: "trigger",
}

func NewStage() Stage {
	return DYNAMIC
}

func (s Stage) Name() string {
	if value, ok := values[s]; ok {
		return value
	}
	logger.Error(fmt.Sprintf("Stage not found for id: '%v'", s))
	return values[DYNAMIC]
}

func GetStages() []string {
	var stages []string
	for _, value := range values {
		stages = append(stages, value)
	}
	return stages
}
