package types

import (
	"reflect"

	"kapigen.kateops.com/internal/types"
)

type Enum[T any] struct {
	Values map[string]T
}

func NewEnum[T any](value map[string]T) Enum[T] {
	return Enum[T]{
		Values: value,
	}
}
func (e *Enum[T]) Validate() error {
	if e.Values == nil {
		return types.ErrorHandler("should have values", 3)
	}

	return nil
}

func (e *Enum[T]) Value(key string) (*T, error) {
	if value, ok := e.Values[key]; ok {
		return &value, nil
	}
	return nil, types.ErrorHandler("value not found", 3)
}

func (e *Enum[T]) Key(value T) (string, error) {
	for k, v := range e.Values {
		if reflect.DeepEqual(v, value) {
			return k, nil
		}
	}
	return "", types.ErrorHandler("value not found", 3)
}

func (e *Enum[T]) KeySafe(value T) string {
	if key, err := e.Key(value); err == nil {
		return key
	}
	return ""

}
