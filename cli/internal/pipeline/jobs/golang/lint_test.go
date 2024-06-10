package golang

import (
	"testing"
)

func TestLint(t *testing.T) {
	t.Run("can create new job", func(t *testing.T) {
		job := Lint("golang:latest", ".")
		if job == nil {
			t.Error("should not be nil")
		}
		if job == nil && job.CiJob.Image.Name != "golang:latest" {
			t.Error("should be equal")
		}
		if job.CiJob.Artifact.Name == "" {
			t.Error("should not be empty")
		}
	})
}
