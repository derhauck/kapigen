package environment

import (
	"testing"
)

func TestIsRelease(t *testing.T) {
	t.Run("will be wrong without env", func(t *testing.T) {
		//SetLocalEnv()
		CI_MERGE_REQUEST_ID.Unset()
		CI_COMMIT_BRANCH.Unset()
		CI_DEFAULT_BRANCH.Unset()
		if IsRelease() {
			t.Errorf("should be false as no env was prepared, CI_MR = %s, COMMIT_BRANCH = %s, DEFAULT_BRANCH = %s", CI_MERGE_REQUEST_ID.Get(), CI_COMMIT_BRANCH.Get(), CI_DEFAULT_BRANCH.Get())
		}
	})
	t.Run("will work with ci vars", func(t *testing.T) {
		CI_MERGE_REQUEST_ID.Unset()
		CI_COMMIT_BRANCH.Set("master")
		CI_DEFAULT_BRANCH.Set("master")
		if !IsRelease() {
			t.Errorf("should be true as env was prepared id: %s, commit: %s, default :%s", CI_MERGE_REQUEST_ID.Get(), CI_COMMIT_BRANCH.Get(), CI_DEFAULT_BRANCH.Get())
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

func TestGetMergeRequestId(t *testing.T) {
	t.Run("Can get ID from commit message", func(t *testing.T) {
		CI_COMMIT_MESSAGE.Set("add missing workflow rules\n\nSee merge request !3")
		result := getMergeRequestIdFromCommit(CI_COMMIT_MESSAGE.Get())
		if result != 3 {
			t.Errorf("Can not get ID from commit message, expected: %s, received: %d", "3", result)
		}
	})
	t.Run("Can parse ID from commit message correctly", func(t *testing.T) {
		CI_COMMIT_MESSAGE.Set("add missing workflow rules\n\nSee merge request !03")
		result := getMergeRequestIdFromCommit(CI_COMMIT_MESSAGE.Get())
		if result != 3 {
			t.Errorf("Can not get ID from commit message, expected: %s, received: %d", "3", result)
		}

		CI_COMMIT_MESSAGE.Set("add missing workflow rules\n\nSee merge request !1.3")
		result = getMergeRequestIdFromCommit(CI_COMMIT_MESSAGE.Get())
		if result != 1 {
			t.Errorf("Can not get ID from commit message, expected: %s, received: %d", "1", result)
		}
	})
	t.Run("Can get ID from env", func(t *testing.T) {
		CI_MERGE_REQUEST_ID.Set("123")
		result := getMergeRequestIdFromEnv()
		if result != 123 {
			t.Errorf("Can not get ID from env, expected: %s, received: %d", "123", result)
		}
	})
}
