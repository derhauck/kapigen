package config

import (
	"errors"
	"reflect"
	"slices"
	"testing"

	"gitlab.com/kateops/kapigen/cli/factory"
	"gitlab.com/kateops/kapigen/cli/internal/cli"
	"gitlab.com/kateops/kapigen/cli/internal/docker"
	"gitlab.com/kateops/kapigen/cli/internal/version"
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

func TestGolang_Build(t *testing.T) {
	t.Run("can create valid jobs", func(t *testing.T) {
		golang := &Golang{
			ImageName: "golang:1.21",
			Path:      "cli",
			Coverage:  &GolangCoverage{},
			Lint: &GolangLint{
				imageName: docker.GOLANG_GOLANGCI_LINT.String(),
			},
		}
		main := factory.New(&cli.Settings{
			Mode:         version.Gitlab,
			PrivateToken: "",
		})

		build, err := golang.Build(main, GOLANG, "id")
		if err != nil {
			t.Error(err)
		}
		if len(build.GetJobs()) != 2 {
			t.Errorf("expected 2 jobs, got %d", len(build.GetJobs()))
		}
		for _, job := range build.GetJobs() {
			// general job validation
			if job.CiJob == nil {
				t.Error("job.CiJob is nil")
				t.FailNow()
			}

			if len(job.CiJob.Tags) == 0 {
				t.Error("expected tags to be set")
			}

			// golang lint job validation
			if slices.Contains(job.Names, "lint") {
				if job.CiJob.Image.Name != golang.Lint.imageName {
					t.Errorf("expected image name to be %s, received %s", golang.Lint.imageName, job.CiJob.Image.Name)
				}
			}

			// golang test job validation
			if slices.Contains(job.Names, "test") {
				if job.CiJob.Image.Name != golang.ImageName {
					t.Errorf("expected image name to be %s, received %s", golang.ImageName, job.CiJob.Image.Name)
				}
				if job.CiJob.Artifact.Reports.Junit.Value == "" {
					t.Error("expected artifact report - junit to be set")
				}
			}
		}

	})
}

func TestGolang_New(t *testing.T) {
	t.Run("creating correct type", func(t *testing.T) {
		old := &Golang{}
		golang := old.New()
		if golang == nil {
			t.Error("should not be nil")
		}
		if old == golang {
			t.Error("Golang was not a new instance")
		}
		if _, ok := golang.(*Golang); !ok {
			t.Errorf("should be of type Golang")
		}
	})
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
