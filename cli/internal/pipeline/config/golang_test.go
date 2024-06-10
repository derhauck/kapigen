package config

import (
	"reflect"
	"testing"

	"kapigen.kateops.com/internal/environment"
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
				Path:      "cli",
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
