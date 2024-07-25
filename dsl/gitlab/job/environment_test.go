package job

import "testing"

func TestEnvironmentActionEnum(t *testing.T) {
	t.Run("can run enum", func(t *testing.T) {
		values := EnvironmentActionEnum().GetValues()
		if len(values) == 0 {
			t.Error("should not be empty")
			t.FailNow()
		}
	})
}
