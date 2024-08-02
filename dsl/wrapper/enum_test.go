package wrapper

import (
	"slices"
	"testing"
)

func TestEnum(t *testing.T) {
	type test byte
	const (
		keyA test = iota
		keyB
		keyC
	)

	values := map[test]string{
		keyA: "a",
		keyB: "b",
	}
	t.Run("can create and use enum", func(t *testing.T) {

		enum, err := NewEnum[test](values)
		if err != nil {
			t.Errorf("expected nil, received %s", err)
		}
		if enum.Validate() != nil {
			t.Errorf("expected nil, received %s", err)
		}

		if result, err := enum.Value(keyA); err != nil {
			t.Errorf("expected nil, received %s", err)
		} else if result != "a" {
			t.Errorf("expected a, received %s", result)
		}

		result, err := enum.Value(keyC)
		if err == nil {
			t.Error("expected error, received nil")
		}
		if result != "" {
			t.Errorf("expected empty string, received %s", result)
		}
		result = enum.ValueSafe(keyB)
		if result != "b" {
			t.Errorf("expected 'b', received %s", result)
		}

		result = enum.ValueSafe(keyC)
		if result != "" {
			t.Errorf("expected '', received %v", result)
		}
		value, err := enum.Value(keyA)
		if err != nil {
			t.Errorf("expected nil, received %s", err)
		}
		if value != "a" {
			t.Errorf("expected 'a', received %s", value)
		}
		key, err := enum.KeyFromValue("a")
		if err != nil {
			t.Error(err)
		}

		if key != keyA {
			t.Errorf("expected 0 (KeyA), received %d", key)
		}
		_, err = enum.KeyFromValue("d")
		if err == nil {
			t.Error("expected error, received nil")
		}

	})
	t.Run("can not create empty enum", func(t *testing.T) {
		enum, err := NewEnum[test](map[test]string{})
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if err = enum.Validate(); err == nil {
			t.Error("expected error, received nil")
		}
	})
	t.Run("can not create enum", func(t *testing.T) {
		type invalid string

		invalidValues := map[invalid]string{
			invalid("a"): "a",
			invalid("b"): "b",
		}
		enum, err := NewEnum[test](invalidValues)
		if err == nil {
			t.Error("expected error, received nil")
		}
		if enum != nil {
			t.Errorf("expected nil, received %v", enum)
		}
	})
	t.Run("can get values", func(t *testing.T) {
		enum, err := NewEnum[test](values)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		enumValues := enum.GetValues()
		if len(enumValues) == 0 {
			t.Error("expected values, received empty")
		}
		if !slices.ContainsFunc(enumValues, func(n string) bool {
			for _, v := range values {
				if v == n {
					return true
				}
			}
			return false
		}) {
			t.Error("expected values to be same as input")
		}
	})
}
