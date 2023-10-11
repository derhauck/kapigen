package job

import (
	"kapigen.kateops.com/internal/when"
	"reflect"
	"testing"
)

func TestNewRulesYaml(t *testing.T) {
	t.Run("Can create empty Rules", func(t *testing.T) {
		t.Parallel()
		result := NewRulesYaml(Rules{})
		if result == nil {
			t.Error("Can not create RulesYaml from Empty Rules")
		}
	})
	t.Run("Can create Rules", func(t *testing.T) {
		t.Parallel()
		result := NewRulesYaml(Rules{&Rule{}})
		if result == nil {
			t.Error("Can not create RulesYaml")
		}
	})
}

func TestWhenEnumType_When(t *testing.T) {
	t.Run("Can create When from type", func(t *testing.T) {
		t.Parallel()
		result := NewWhen(when.Always)
		if reflect.ValueOf(result.Get()).Kind() != reflect.String {
			t.Error("NewWhen does not create a string")
		}
	})
	t.Run("Can get Default from Empty", func(t *testing.T) {
		t.Parallel()
		expectation := NewWhen(when.OnSuccess)
		expectation.Get()
		test := WhenWrapper{}
		if test.Get() != expectation.Get() {
			t.Errorf("Expected on_success, got: '%s'", test.Get())
		}

	})
}
