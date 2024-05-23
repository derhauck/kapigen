package job

import (
	"reflect"
	"strings"
	"testing"

	"kapigen.kateops.com/internal/when"
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
func TestDefaultPipelineRules(t *testing.T) {
	t.Parallel()
	t.Run("Can get different paths in changes", func(t *testing.T) {
		expectations := []string{"test/", "test2", "test3"}
		for _, expectation := range expectations {
			rules := DefaultPipelineRules([]string{expectation})
			for _, rule := range rules.Get() {
				for _, change := range rule.Changes.Get() {
					t.Run(expectation, func(t *testing.T) {
						if !strings.Contains(change, expectation) {
							t.Errorf("Expected '%s' to contain '%s'", change, expectation)
						}
						if strings.Contains(change, "//") {
							t.Errorf("Expected '%s' to not contain '%s' for value: '%s'", change, "//", expectation)
						}
						if i := strings.Index(change, "."); i == 0 {
							t.Errorf("Expected '%s' to not begin with '%s' for value: '%s'", change, ".", expectation)
						}
					})

				}
			}
		}

	})
	t.Run("Can get proper path for '.'", func(t *testing.T) {
		expectation := "."
		rules := DefaultPipelineRules([]string{expectation})
		for _, rule := range rules.Get() {
			for _, change := range rule.Changes.Get() {
				t.Run(expectation, func(t *testing.T) {
					if strings.Contains(change, expectation) {
						t.Errorf("Expected '%s' not to contain '%s'", change, expectation)
					}
					if strings.Contains(change, "//") {
						t.Errorf("Expected '%s' to not contain '%s' for value: '%s'", change, "//", expectation)
					}
					if i := strings.Index(change, "."); i == 0 {
						t.Errorf("Expected '%s' to not begin with '%s' for value: '%s'", change, ".", expectation)
					}
				})

			}
		}

	})
}
