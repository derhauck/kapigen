package enum

import (
	"gitlab.com/kateops/kapigen/dsl/logger"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

type When int

const (
	WhenOnSuccess When = iota
	WhenOnFailure
	WhenAlways
	WhenNever
)

func (w When) String() string {
	return WhenEnum().ValueSafe(w)
}

func WhenEnum() *wrapper.Enum[When, string] {
	enum, err := wrapper.NewEnum[When](map[When]string{
		WhenOnSuccess: "on_success",
		WhenOnFailure: "on_failure",
		WhenAlways:    "always",
		WhenNever:     "never",
	})
	if err != nil {
		logger.ErrorE(err)
	}
	return enum
}
