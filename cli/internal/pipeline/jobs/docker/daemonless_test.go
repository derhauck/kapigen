package docker

import (
	"testing"

	job2 "gitlab.com/kateops/kapigen/dsl/gitlab/job"
)

func TestNewDaemonlessBuildkitBuild(t *testing.T) {
	t.Run("can create new job", func(t *testing.T) {
		job, err := NewDaemonlessBuildkitBuild("testImage", ".", ".", "Dockerfile", []string{"${CI_REGISTRY_IMAGE}:1.0.0", "${CI_REGISTRY_IMAGE}:latest"}, []string{})
		if err != nil {
			t.Error(err)
		}
		if job == nil {
			t.Error("should be able to create new job")
			t.FailNow()
		}
		job.CiJob.Rules.AddRules(*job2.DefaultMainBranchRules())
		if err := job.Render(); err != nil {
			t.Error(err)
		}
	})
	t.Run("can not create new job", func(t *testing.T) {
		job, err := NewDaemonlessBuildkitBuild("testImage", ".", ".", "Dockerfile", []string{}, []string{})
		if err == nil || job != nil {
			t.Error("should not be able to create new job")
		}
	})
}
