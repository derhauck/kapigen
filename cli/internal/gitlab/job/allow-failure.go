package job

import "kapigen.kateops.com/internal/pipeline/wrapper"

type AllowFailure struct {
	ExitCodes wrapper.IntSlice `yaml:"exit_codes"`
	Failure   bool
}

func (a *AllowFailure) AllowAll() {
	a.Failure = true
}

func (a *AllowFailure) add(code int32) *AllowFailure {
	a.ExitCodes.Add(code)
	return a
}

func (a *AllowFailure) addSeveral(codes ...int32) *AllowFailure {
	a.ExitCodes.AddSeveral(codes)
	return a
}

func (a *AllowFailure) Get() any {
	return NewAllowFailureYaml(a)
}

type AllowFailureYaml struct {
	ExitCodes []int32 `yaml:"exit_codes"`
}

func NewAllowFailureYaml(allowFailure *AllowFailure) any {
	if len(allowFailure.ExitCodes.Get()) == 0 {
		return allowFailure.Failure
	}

	return &AllowFailureYaml{
		ExitCodes: allowFailure.ExitCodes.Get(),
	}
}
