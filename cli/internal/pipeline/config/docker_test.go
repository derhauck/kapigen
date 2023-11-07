package config

import (
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

//func TestDocker_Build(t *testing.T) {
//	t.Run("Can build Docker", func(t *testing.T) {
//		docker := Docker{
//			Path:    "test",
//			Context: "not set",
//		}
//		jobs, err := docker.Build(factory.New(cli.NewSettings(cli.SetMode(version.Gitlab))), DockerPipeline, "testId")
//		if err != nil {
//			t.Errorf("Build failed with: %s", err.Error())
//		}
//		if len(jobs.GetJobs()) == 0 {
//			t.Error("Docker pipeline is empty")
//		}
//
//	})
//}

func TestDocker_Validate(t *testing.T) {
	t.Run("Can validate valid Docker config", func(t *testing.T) {
		docker := Docker{
			Path:    "test",
			Context: "context",
		}
		err := docker.Validate()
		if err != nil {
			t.Error(err)
		}
		if docker.Path != "test" {
			t.Error("path should be test")
		}
		if docker.Context != "context" {
			t.Error("context should be context")
		}
		if docker.Dockerfile != "Dockerfile" {
			t.Error("dockerfile should be Dockerfile")
		}
	})
	t.Run("Can validate valid Docker config", func(t *testing.T) {
		docker := Docker{
			Path: "test",
		}
		err := docker.Validate()
		if err != nil {
			t.Error(err)
		}
		if docker.Path != "test" {
			t.Error("path should be test")
		}
		if docker.Context != docker.Path {
			t.Error("context should be same as path")
		}
		if docker.Dockerfile != "Dockerfile" {
			t.Error("dockerfile should be Dockerfile")
		}
	})
}
