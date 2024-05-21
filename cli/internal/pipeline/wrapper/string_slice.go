package wrapper

import "fmt"

type StringSlice struct {
	Value []string
}

func (s *StringSlice) Add(script string) *StringSlice {
	s.Value = append(s.Value, script)
	return s
}
func (s *StringSlice) Addf(script string, a ...any) *StringSlice {
	s.Value = append(s.Value, fmt.Sprintf(script, a...))
	return s
}
func (s *StringSlice) AddSeveral(script ...string) *StringSlice {
	s.Value = append(s.Value, script...)
	return s
}

func (s *StringSlice) AddSlice(script []string) *StringSlice {
	return s.AddSeveral(script...)
}

func (s *StringSlice) Get() []string {
	return s.Value
}

func (s *StringSlice) Has(script string) bool {
	for _, value := range s.Value {
		if value == script {
			return true
		}
	}
	return false
}

func NewStringSlice() *StringSlice {
	return &StringSlice{
		Value: []string{},
	}
}
