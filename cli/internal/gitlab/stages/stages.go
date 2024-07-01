package stages

import (
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/types"
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

func Enum() *types.Enum[Stage, string] {
	enum, err := types.NewEnum[Stage](values)
	if err != nil {
		logger.Error(err.Error())
	}
	return enum
}

func NewStage() Stage {
	return DYNAMIC
}
