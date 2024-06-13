package types

import "testing"

func TestEnum(t *testing.T) {
	t.Run("can create and use enum", func(t *testing.T) {
		type test int
		const (
			keyA test = iota
			keyB
			keyC
		)

		values := map[string]test{
			"a": keyA,
			"b": keyB,
		}

		enum := NewEnum(values)
		if enum.Validate() != nil {
			t.Error("expected nil")
		}

		if result, err := enum.Key(keyA); err != nil {
			t.Error("expected nil")
		} else if result != "a" {
			t.Errorf("expected a, received %s", result)
		}

		result, err := enum.Key(keyC)
		if err == nil {
			t.Errorf("expected error, received %s", result)
		}
		if result != "" {
			t.Errorf("expected empty string, received %s", result)
		}
		value, err := enum.Value("a")
		if err != nil {
			t.Error("expected nil")
		}
		if *value != keyA {
			t.Errorf("expected %d, received %d", keyA, *value)
		}
	})
}
