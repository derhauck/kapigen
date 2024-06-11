package config

import (
	"errors"
	"strings"
	"testing"

	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/cli"
	errors2 "kapigen.kateops.com/internal/types"
	"kapigen.kateops.com/internal/version"
)

func TestPhpComposerValidate(t *testing.T) {
	t.Run("sets default path and args when empty", func(t *testing.T) {
		composer := &PhpComposer{}
		err := composer.Validate()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if composer.Path != "." {
			t.Errorf("Expected path to be '.', got %s", composer.Path)
		}
		if composer.Args != "--no-progress --no-cache --no-interaction" {
			t.Errorf("Expected args to be '--no-progress --no-cache --no-interaction', got %s", composer.Args)
		}
	})
}

func TestPhpunitValidate(t *testing.T) {
	composer := &PhpComposer{Path: "/path/to/composer"}

	t.Run("sets default path and bin when empty", func(t *testing.T) {
		phpunit := &Phpunit{}
		err := phpunit.Validate(composer)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if phpunit.Path != "/path/to/composer" {
			t.Errorf("Expected path to be '/path/to/composer', got %s", phpunit.Path)
		}
		if phpunit.Bin != "/path/to/composer/vendor/bin/phpunit" {
			t.Errorf("Expected bin to be '/path/to/composer/vendor/bin/phpunit', got %s", phpunit.Bin)
		}
	})

}

func TestPhpValidate(t *testing.T) {

	t.Run("returns error when phpunit job creation fails", func(t *testing.T) {
		php := &Php{
			ImageName: "",
		}
		err := php.Validate()
		if err == nil {
			t.Error("Expected error when phpunit job creation fails, but got nil")
		}
		var re *errors2.DetailedError
		if !errors.As(err, &re) {
			t.Errorf("Expected error to be of type 'DetailedError', got '%s'", err.Error())
		}
		expectedErr := "no imageName set, required"
		if strings.Contains(err.Error(), expectedErr) {
			t.Errorf("Expected error to contain '%s', got '%s'", expectedErr, err.Error())
		}
	})

	// ... (other test cases)
}
func TestPhpRules(t *testing.T) {
	t.Run("returns default pipeline rules based on changes", func(t *testing.T) {
		php := &Php{
			InternalChanges: []string{"/path/to/composer", "/path/to/context"},
		}

		rules := php.Rules()
		if rules == nil {
			t.Error("Expected non-nil rules, got nil")
		}

		rulesSlice := rules.Get()
		if len(rulesSlice) == 0 {
			t.Error("Expected at least one rule, got empty slice")
		}

		for _, rule := range rulesSlice {
			changes := rule.Changes.Get()
			for _, change := range changes {
				found := false
				for _, expectedChange := range php.InternalChanges {
					if strings.Contains(change, expectedChange) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Unexpected change path '%s' in rules", change)
				}
			}
		}
	})
}

func TestPhpBuild(t *testing.T) {
	mainFactory := factory.New(&cli.Settings{Mode: version.Gitlab})
	pipelineType := PHPPipeline

	t.Run("returns error when phpunit job creation fails", func(t *testing.T) {
		php := &Php{
			ImageName: "",
		}
		_, err := php.Build(mainFactory, pipelineType, "test-id")
		if err == nil {
			t.Error("Expected error when phpunit job creation fails, but got nil")
		}
		expectedErr := "no imageName set, required"
		if strings.Contains(err.Error(), expectedErr) {
			t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("creates docker pipeline jobs when Docker is present", func(t *testing.T) {
		php := &Php{
			ImageName: "php:7.4",
			Docker: &SlimDocker{
				Path:    "/path/to/docker",
				Context: "/path/to/context",
			},
		}
		_ = php.Validate()
		jobs, err := php.Build(mainFactory, pipelineType, "test-id")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if jobs == nil {
			t.Error("Expected non-nil jobs, got nil")
		}

		jobSlice := jobs.GetJobs()
		if len(jobSlice) == 0 {
			t.Error("Expected at least one job, got empty slice")
		}

		// Add more assertions to check the created jobs
	})

	t.Run("creates service jobs when Services are present", func(t *testing.T) {
		php := &Php{
			ImageName: "php:7.4",
			Services: Services{
				{
					Name:      "test-service",
					Port:      8080,
					ImageName: "test-image",
				},
			},
		}
		_ = php.Validate()
		jobs, err := php.Build(mainFactory, pipelineType, "test-id")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if jobs == nil {
			t.Error("Expected non-nil jobs, got nil")
		}

		jobSlice := jobs.GetJobs()
		if len(jobSlice) == 0 {
			t.Error("Expected at least one job, got empty slice")
		}

		// Add more assertions to check the created jobs and services
	})
}
