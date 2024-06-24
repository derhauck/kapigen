package types

import (
	"reflect"

	"kapigen.kateops.com/internal/types"
)

type Enumeration interface {
	comparable
}
type Enum[K Enumeration, V any] struct {
	Values map[K]V
}

func NewEnum[K Enumeration, V string](value interface{}) (*Enum[K, V], error) {

	final, ok := value.(map[K]V)
	if !ok {
		return nil, types.ErrorHandler("value not of type byte(iota)", 3)
	}

	return &Enum[K, V]{
		Values: final,
	}, nil
}
func (e *Enum[K, V]) Validate() error {
	if e.Values == nil || len(e.Values) == 0 {
		return types.ErrorHandler("should have values", 3)
	}

	return nil
}

func (e *Enum[K, V]) KeyFromValue(value V) (K, error) {
	for k, v := range e.Values {
		if reflect.DeepEqual(v, value) {
			return k, nil
		}
	}
	var inf interface{}
	empty, _ := inf.(K)
	return empty, types.ErrorHandler("value not found", 3)
}

func (e *Enum[K, V]) Key(key K) (V, error) {
	if v, ok := e.Values[key]; ok {
		return v, nil
	}
	var inf interface{}
	empty, _ := inf.(V)
	return empty, types.ErrorHandler("value not found", 3)
}

func (e *Enum[K, V]) KeySafe(key K) V {
	if key, err := e.Key(key); err == nil {
		return key
	}

	var inf interface{}
	empty, _ := inf.(V)
	return empty

}
