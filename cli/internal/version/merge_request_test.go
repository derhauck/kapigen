package version

import (
	"kapigen.kateops.com/internal/environment"
	"testing"
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
