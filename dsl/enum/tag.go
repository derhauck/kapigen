package enum

import (
	"gitlab.com/kateops/kapigen/dsl/logger"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

type Tag int

const (
	TagPressureMedium Tag = iota
	TagPressureExclusive
)

func (w Tag) String() string {
	return TagEnum().ValueSafe(w)
}

func TagEnum() *wrapper.Enum[Tag, string] {
	enum, err := wrapper.NewEnum[Tag](map[Tag]string{
		TagPressureMedium:    "pressure:medium",
		TagPressureExclusive: "pressure:exclusive",
	})
	if err != nil {
		logger.Error(err.Error())
	}
	return enum
}
