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
