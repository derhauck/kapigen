package gitlab

import "kapigen.kateops.com/internal/pipeline/wrapper"

type AfterScript struct {
	Value wrapper.StringSlice
}

func (a *AfterScript) getRenderedValue() []string {
	return a.Value.Get()
}
