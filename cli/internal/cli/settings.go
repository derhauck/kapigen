package cli

import (
	"kapigen.kateops.com/internal/version"
)

type Settings struct {
	Mode version.Mode
}

type SettingsFn func(s *Settings)

func NewSettings(fns ...SettingsFn) *Settings {
	s := &Settings{}
	for _, fn := range fns {
		fn(s)
	}
	return s
}

func SetMode(mode version.Mode) SettingsFn {
	return func(s *Settings) {
		s.Mode = mode
	}
}
