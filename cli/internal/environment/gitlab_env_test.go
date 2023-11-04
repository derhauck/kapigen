package environment

import (
	"os"
	"testing"
)

const UNREACHABLE = 50000000

func TestGet(t *testing.T) {
	t.Run("Can get project id", func(t *testing.T) {
		t.Parallel()
		expectation := "PROJECT_ID_17"
		err := os.Setenv("CI_PROJECT_ID", expectation)
		if err != nil {
			t.Error(err)
		}
		result := CI_PROJECT_ID.Get()
		if result != expectation {
			t.Errorf("Can not get Project ID, expected: %s, received: %s", expectation, result)
		}
	})
	t.Run("Can get merge request id", func(t *testing.T) {
		t.Parallel()
		expectation := "MERGE_REQUEST_ID_17"
		err := os.Setenv("CI_MERGE_REQUEST_ID", expectation)
		if err != nil {
			t.Error(err)
		}
		result := CI_MERGE_REQUEST_ID.Get()
		if result != expectation {
			t.Error()
		}
	})
}

func TestLookup(t *testing.T) {

	t.Run("Can lookup project id", func(t *testing.T) {
		t.Parallel()
		expectation := "PROJECT_ID_17"
		err := os.Setenv("CI_PROJECT_ID", expectation)
		if err != nil {
			t.Error(err)
		}
		result, err := CI_PROJECT_ID.Lookup()
		if err != nil {
			t.Errorf("Can not get Project ID, err: %s", err.Error())
		}
		if result != expectation {
			t.Errorf("Can not get Project ID, expected: %s, received: %s", expectation, result)
		}
	})
	t.Run("Can lookup merge request id", func(t *testing.T) {
		t.Parallel()
		expectation := "MERGE_REQUEST_ID_17"
		err := os.Setenv("CI_MERGE_REQUEST_ID", expectation)
		if err != nil {
			t.Error(err)
		}
		result, err := CI_MERGE_REQUEST_ID.Lookup()
		if err != nil {
			t.Errorf("Can not get Merge Request ID, err: %s", err.Error())
		}
		if result != expectation {
			t.Error()
		}
	})
	t.Run("Can not lookup test var", func(t *testing.T) {
		t.Parallel()
		var test Variable = UNREACHABLE
		result, err := test.Lookup()
		if err != nil && err.Error() != "env var '' is not set" {
			t.Errorf("Unexpected error looking up test var: %s", err.Error())
		}
		if result != "" {
			t.Errorf("Should not be able to lookup test value: '%s'", result)
		}
	})
}

func TestValue(t *testing.T) {
	t.Run("Can not get test var", func(t *testing.T) {
		t.Parallel()
		var test Variable = UNREACHABLE
		if test.Key() == "test" {
			t.Error("No Variable for 'test' should exist")
		}
	})
	t.Run("Can get CI vars", func(t *testing.T) {
		t.Parallel()
		var varsToCheck = []Variable{
			CI_MERGE_REQUEST_ID,
			CI_PROJECT_ID,
		}
		for index, value := range varsToCheck {
			if value.Key() == "" {
				t.Errorf("CI variable for number '%v' should exist varsToCheck number '%v'", value, index)
			}
		}
	})

}

func TestVariable_Unset(t *testing.T) {
	t.Run("Can unset CI vars", func(t *testing.T) {
		t.Parallel()
		var varsToCheck = []Variable{
			CI_MERGE_REQUEST_ID,
			CI_PROJECT_ID,
		}
		for _, value := range varsToCheck {
			value.Set("123")
			value.Unset()
			if value.Get() != "" {
				t.Errorf("CI variable for number '%v' should be unset", value)
			}
		}
	})
}
