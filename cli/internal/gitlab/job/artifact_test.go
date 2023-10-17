package job

import (
	"kapigen.kateops.com/internal/pipeline/wrapper"
	"kapigen.kateops.com/internal/when"
	"testing"
)

func TestArtifactsSerialization(t *testing.T) {
	// Create an instance of the Ci struct
	artifacts := &Artifact{
		Paths:     *(wrapper.NewStringSlice().AddSeveral([]string{"path1", "path2"})),
		ExpireIn:  "7 days",
		ExposeAs:  "public",
		Name:      "my-artifacts",
		Untracked: true,
		When:      NewWhen(when.OnSuccess),
	}

	artifactsYaml := NewArtifactsYaml(artifacts)

	if artifactsYaml.Paths == nil {
		t.Error("Expected non-nil Paths in Yaml")
	}
	if len(artifactsYaml.Paths) != 2 {
		t.Error("Expected 2 paths in Yaml")
	}
	if artifactsYaml.ExpireIn != "7 days" {
		t.Errorf("Expected ExpireIn to be '7 days', got %s", artifactsYaml.ExpireIn)
	}
	if artifactsYaml.ExposeAs != "public" {
		t.Errorf("Expected ExposeAs to be 'public', got %s", artifactsYaml.ExposeAs)
	}
	if artifactsYaml.Name != "my-artifacts" {
		t.Errorf("Expected Name to be 'my-artifacts', got %s", artifactsYaml.Name)
	}
	if artifactsYaml.Untracked != true {
		t.Error("Expected Untracked to be true")
	}
	if artifactsYaml.When != "on_success" {
		t.Errorf("Expected When to be 'on_success', got %s", artifactsYaml.When)
	}
}
