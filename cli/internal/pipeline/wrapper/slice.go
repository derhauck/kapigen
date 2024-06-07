package wrapper

type Slice[T any] interface {
	Add(str T) *Slice[T]
	AddSeveral(strSlice []T) *Slice[T]
	Get() []T
}

type Array[T any] struct {
	slice []T
}

func (s *Array[T]) Push(str T) *Array[T] {
	s.slice = append(s.slice, str)
	return s
}

func (s *Array[T]) ForEach(fn func(element T)) {
	for _, element := range s.slice {
		fn(element)
	}
}
