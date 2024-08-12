package stages

import (
	"gitlab.com/kateops/kapigen/dsl/logger"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
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

func Enum() *wrapper.Enum[Stage, string] {
	enum, err := wrapper.NewEnum[Stage](map[Stage]string{
		LINT:    "lint",
		INIT:    "init",
		BUILD:   "build",
		TEST:    "test",
		RELEASE: "release",
		DYNAMIC: "dynamic",
		TRIGGER: "trigger",
		FINAL:   "final",
	})
	if err != nil {
		logger.ErrorE(err)
	}
	return enum
}

func NewStage() Stage {
	return DYNAMIC
}
