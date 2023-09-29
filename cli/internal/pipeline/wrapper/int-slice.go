package wrapper

type IntSlice struct {
	Value []int32
}

func (s *IntSlice) Add(script int32) *IntSlice {
	s.Value = append(s.Value, script)
	return s
}
func (s *IntSlice) AddSeveral(script []int32) *IntSlice {
	s.Value = append(s.Value, script...)
	return s
}

func (s *IntSlice) Get() []int32 {
	return s.Value
}
