package gitlab

import "kapigen.kateops.com/internal/pipeline/wrapper"

type Script struct {
	Value wrapper.StringSlice
}

func (s *Script) getRenderedValue() []string {
	return s.Value.Get()
}

func NewScript() Script {
	return Script{
		Value: *wrapper.NewStringSlice(),
	}
}
