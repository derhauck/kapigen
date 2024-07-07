package job

import "gitlab.com/kateops/kapigen/dsl/wrapper"

type BeforeScript struct {
	Value *wrapper.Array[string]
}

func (a *BeforeScript) GetRenderedValue() []string {
	if a.Value != nil {
		return a.Value.Get()
	}
	return nil
}

func NewBeforeScript() BeforeScript {
	return BeforeScript{
		Value: wrapper.NewArray[string](),
	}
}
