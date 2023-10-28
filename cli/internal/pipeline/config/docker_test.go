package config

import (
	"kapigen.kateops.com/factory"
	"testing"
)

func TestDocker_New(t *testing.T) {
	t.Run("Can create new Docker", func(t *testing.T) {
		t.Parallel()
		oldDocker := Docker{
			Path: "old",
		}
		newDocker := oldDocker.New()
		if result, ok := newDocker.(*Docker); ok {
			if oldDocker.Path == result.Path {
				t.Error("Could not create New struct from type Docker")
			}
		}
	})
}

func TestDocker_Build(t *testing.T) {
	t.Run("Can build Docker", func(t *testing.T) {
		t.Parallel()
		docker := Docker{
			Path:    "test",
			Context: "not set",
		}
		jobs, err := docker.Build(factory.New(), DockerPipeline, "testId")
		if err != nil {
			t.Errorf("Build failed with: %s", err.Error())
		}
		if len(jobs.GetJobs()) == 0 {
			t.Error("Docker pipeline is empty")
		}

	})
}
