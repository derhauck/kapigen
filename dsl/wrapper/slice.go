package wrapper

import "reflect"

type Array[T any] struct {
	slice []T
}

func NewArray[T any]() *Array[T] {
	return &Array[T]{
		slice: make([]T, 0),
	}
}

func (s *Array[T]) Get() []T {
	return s.slice
}

func (s *Array[T]) Push(elements ...T) *Array[T] {
	s.slice = append(s.slice, elements...)
	return s
}

func (s *Array[T]) ForEach(fn func(element *T)) {
	for _, element := range s.slice {
		fn(&element)
	}
}

func (s *Array[T]) Find(fn func(element *T) bool) (*T, int) {
	for index, element := range s.slice {
		if fn(&element) {
			return &element, index
		}
	}
	return nil, -1
}

func (s *Array[T]) Has(element T) bool {
	_, index := s.Find(func(e *T) bool {
		return reflect.DeepEqual(*e, element)
	})
	return index != -1
}

func (s *Array[T]) Length() int {
	return len(s.slice)
}

func (s *Array[T]) Map(fn func(element *T) T) *Array[T] {
	newSlice := make([]T, len(s.slice))
	for index, element := range s.slice {
		newSlice[index] = fn(&element)
	}
	return &Array[T]{
		slice: newSlice,
	}
}
