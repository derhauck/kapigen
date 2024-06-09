package job

import "kapigen.kateops.com/internal/pipeline/wrapper"

type AfterScript struct {
	Value *wrapper.Array[string]
}

func (a *AfterScript) GetRenderedValue() []string {
	if a.Value != nil {
		return a.Value.Get()
	}
	return nil
}

func NewAfterScript() AfterScript {
	return AfterScript{
		Value: wrapper.NewArray[string](),
	}
}
