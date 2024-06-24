package types

import "testing"

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

		if result, err := enum.Key(keyA); err != nil {
			t.Errorf("expected nil, received %s", err)
		} else if result != "a" {
			t.Errorf("expected a, received %s", result)
		}

		result, err := enum.Key(keyC)
		if err == nil {
			t.Error("expected error, received nil")
		}
		if result != "" {
			t.Errorf("expected empty string, received %s", result)
		}
		value, err := enum.Key(keyA)
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
	t.Run("can not create enum", func(t *testing.T) {
		enum, err := NewEnum[test](map[test]string{})
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if err = enum.Validate(); err == nil {
			t.Error("expected error, received nil")
		}
	})
}
