package version

import (
	"testing"

	"gitlab.com/kateops/kapigen/dsl/environment"
)

func TestController_GetMergeRequestId(t *testing.T) {
	t.Run("can get version increase from env", func(t *testing.T) {
		environment.CI_MERGE_REQUEST_ID.Set("13")
		environment.CI_MERGE_REQUEST_LABELS.Set("none")
		controller := NewController(Gitlab, nil, nil)
		result := controller.getVersionIncrease("13", 13)
		if result != "none" {
			t.Errorf("should be none, got %s", result)
		}
	})

	t.Run("can get version increase from env", func(t *testing.T) {
		environment.CI_MERGE_REQUEST_ID.Set("13")
		environment.CI_MERGE_REQUEST_LABELS.Set("version::minor,text")
		controller := NewController(Gitlab, nil, nil)
		result := controller.getVersionIncrease("13", 13)
		if result != "minor" {
			t.Errorf("should be minor, got %s", result)
		}
	})

}

func TestController_GetMrLabelsFromEnv(t *testing.T) {
	t.Run("can get labels from env", func(t *testing.T) {
		environment.CI_MERGE_REQUEST_ID.Set("13")
		environment.CI_MERGE_REQUEST_LABELS.Set("none")
		result := getVersionIncreaseFromLabels(environment.CI_MERGE_REQUEST_LABELS.Get())
		if result != "none" {
			t.Errorf("should be none, got %s", result)
		}
	})

	t.Run("can get labels from env", func(t *testing.T) {
		environment.CI_MERGE_REQUEST_ID.Set("13")
		environment.CI_MERGE_REQUEST_LABELS.Set("version::minor,text")
		result := getVersionIncreaseFromLabels(environment.CI_MERGE_REQUEST_LABELS.Get())
		if result != "minor" {
			t.Errorf("should be minor, got %s", result)
		}
	})
}

func TestGetNewVersion(t *testing.T) {
	t.Run("can get new version", func(t *testing.T) {

		result := getNewVersion("1.0.0", "major")
		if result != "2.0.0" {
			t.Errorf("should be 2.0.0, got %s", result)
		}

		result = getNewVersion("1.0.0", "minor")
		if result != "1.1.0" {
			t.Errorf("should be 1.1.0, got %s", result)
		}

		result = getNewVersion("1.0.0", "patch")
		if result != "1.0.1" {
			t.Errorf("should be 1.0.1, got %s", result)
		}
	})
}

func TestGetFeatureBranchVersion(t *testing.T) {
	t.Run("can get feature branch version", func(t *testing.T) {
		result := GetFeatureBranchVersion("1.0.0", "feature/test")
		if result != "1.0.0-feature-test" {
			t.Errorf("should be 1.0.0-feature-test, got %s", result)
		}
	})
}
