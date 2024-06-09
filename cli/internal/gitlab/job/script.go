package job

import "kapigen.kateops.com/internal/pipeline/wrapper"

type Script struct {
	Value *wrapper.Array[string]
}

func (s *Script) GetRenderedValue() []string {
	if s.Value != nil {
		return s.Value.Get()
	}
	return nil
}

func NewScript() Script {
	return Script{
		Value: wrapper.NewArray[string](),
	}
}
