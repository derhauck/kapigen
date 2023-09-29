package gitlab

import "kapigen.kateops.com/internal/pipeline/wrapper"

type AllowFailure struct {
	ExitCodes wrapper.IntSlice `yaml:"exit_codes"`
}

func (a *AllowFailure) add(code int32) *AllowFailure {
	a.ExitCodes.Add(code)
	return a
}

func (a *AllowFailure) addSeveral(codes ...int32) *AllowFailure {
	a.ExitCodes.AddSeveral(codes)
	return a
}

func (a *AllowFailure) Get() []int32 {
	return a.ExitCodes.Get()
}
