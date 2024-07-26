package wrapper

import (
	"errors"
	"strings"
	"testing"
)

func TestNewMissingArgsError(t *testing.T) {
	err := NewMissingArgsError("a", "b", "c")
	if err == nil {
		t.Error("expected error, received nil")
	}
	if !strings.Contains(err.Error(), "'a'") {
		t.Error("expected error to contain 'a', received", err.Error())
	}
	if !strings.Contains(err.Error(), "'b'") {
		t.Error("expected error to contain 'b', received", err.Error())
	}
	if !strings.Contains(err.Error(), "'c'") {
		t.Error("expected error to contain 'c', received", err.Error())
	}
}

func TestNewMissingArgError(t *testing.T) {
	err := NewMissingArgError("missingArg")
	if err == nil {
		t.Error("expected error, received nil")
	}
	if !strings.Contains(err.Error(), "'missingArg'") {
		t.Error("expected error to contain 'missingArg', received", err.Error())
	}
}

func TestDetailedErrorf(t *testing.T) {
	err := DetailedErrorf("test %s", "arg")
	if err == nil {
		t.Error("expected error, received nil")
	}
	if !strings.Contains(err.Error(), "test arg") {
		t.Error("expected error to contain 'test arg', received", err.Error())
	}
}

func TestDetailedErrorE(t *testing.T) {
	testErr := errors.New("test error")
	err := DetailedErrorE(testErr)
	if err == nil {
		t.Error("expected error, received nil")
	}
	if !strings.Contains(err.Error(), "errors_test.go") {
		t.Error("expected error to contain filename 'errors_test.go', received:", err.Error())
	}
}
