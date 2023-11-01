package environment

import (
	"os"
	"testing"
)

func TestIsRelease(t *testing.T) {
	t.Run("will be wrong without env", func(t *testing.T) {
		SetLocalEnv()
		err := os.Unsetenv(CI_MERGE_REQUEST_ID.Key())
		if err != nil {
			t.Error(err.Error())
		}
		if IsRelease() {
			t.Errorf("should be false as no env was prepared, CI_MR = %s", CI_MERGE_REQUEST_ID.Get())
		}
	})
	t.Run("will work with ci vars", func(t *testing.T) {
		CI_MERGE_REQUEST_ID.Set("123")
		CI_COMMIT_BRANCH.Set("master")
		CI_DEFAULT_BRANCH.Set("master")
		if IsRelease() {
			t.Error("should be true as env was prepared")
		}
	})
	t.Run("will not work when not on default branch", func(t *testing.T) {
		CI_MERGE_REQUEST_ID.Set("123")
		CI_COMMIT_BRANCH.Set("feature")
		CI_DEFAULT_BRANCH.Set("master")
		if IsRelease() {
			t.Errorf("should be false as commit is not default branch, %s", CI_COMMIT_BRANCH.Get())
		}
	})
}
