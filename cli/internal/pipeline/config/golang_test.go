package config

import (
	"errors"
	"reflect"
	"testing"

	"gitlab.com/kateops/kapigen/dsl/environment"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

func TestGolangAutoConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Golang
	}{
		{
			name: "Can get golang auto config",
			want: &Golang{
				ImageName: "golang:1.21",
				Path:      "dsl",
			},
		},
	}
	environment.SetLocalEnv()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GolangAutoConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GolangAutoConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGolang_Validate(t *testing.T) {
	t.Run("can validate golang", func(t *testing.T) {
		golang := &Golang{
			ImageName: "golang:1.21",
			Path:      "cli",
		}
		if err := golang.Validate(); err != nil {
			t.Errorf("should not have error: %s", err)
		}
		if golang.ImageName != "golang:1.21" {
			t.Errorf("should be equal")
		}
		if golang.Path != "cli" {
			t.Errorf("should be equal")
		}
	})
	t.Run("can not validate golang", func(t *testing.T) {
		golang := &Golang{
			ImageName: "",
		}
		if err := golang.Validate(); err == nil {
			t.Error("should have error")
		}

	})
}

func TestGolangLint_Validate(t *testing.T) {
	t.Run("can validate golang lint", func(t *testing.T) {
		config := &GolangLint{}
		if err := config.Validate(); err != nil {
			t.Errorf("should not have error: %s", err)
		}
	})
	t.Run("can not validate golang lint", func(t *testing.T) {
		config := &GolangLint{
			Mode: "test",
		}

		if err := config.Validate(); err == nil {
			var detailed *wrapper.DetailedError
			if errors.As(err, &detailed) {
				if detailed.Filename != "golang.go" {
					t.Errorf("should have error in 'golang.go' file, received: %s", detailed.Filename)
				}
			}
			t.Error("should have error")
		}
	})
}
