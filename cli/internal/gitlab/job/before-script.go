package job

import "kapigen.kateops.com/internal/pipeline/wrapper"

type BeforeScript struct {
	Value *wrapper.StringSlice
}

func (a *BeforeScript) GetRenderedValue() []string {
	if a.Value != nil {
		return a.Value.Get()
	}
	return nil
}

func NewBeforeScript() BeforeScript {
	return BeforeScript{
		Value: wrapper.NewStringSlice(),
	}
}
