package cli

import (
	"gitlab.com/kateops/kapigen/cli/internal/version"
)

type Settings struct {
	Mode         version.Mode
	PrivateToken string
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

func SetPrivateToken(privateToken string) SettingsFn {
	return func(s *Settings) {
		s.PrivateToken = privateToken
	}
}
