package job

import "gitlab.com/kateops/kapigen/dsl/wrapper"

type AllowFailure struct {
	ExitCodes wrapper.Array[int32] `yaml:"exit_codes"`
	Failure   bool
}

func (a *AllowFailure) AllowAll() {
	a.Failure = true
}

//func (a *AllowFailure) add(codes ...int32) *AllowFailure {
//	a.ExitCodes.Push(codes...)
//	return a
//}

func (a *AllowFailure) Get() any {
	return NewAllowFailureYaml(a)
}

type AllowFailureYaml struct {
	ExitCodes []int32 `yaml:"exit_codes"`
}

func NewAllowFailureYaml(allowFailure *AllowFailure) any {
	if allowFailure.ExitCodes.Length() == 0 {
		return allowFailure.Failure
	}

	return &AllowFailureYaml{
		ExitCodes: allowFailure.ExitCodes.Get(),
	}
}
