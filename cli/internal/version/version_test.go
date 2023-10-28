package version

import (
	"testing"
)

func TestTag(t *testing.T) {
	t.Parallel()
	t.Run("can get different tag", func(t *testing.T) {
		mode := Gitlab
		if mode.getTag() != Gitlab.Name() {
			t.Error("should be gitlab")
		}
	})
}

func TestGetTag(t *testing.T) {
	t.Parallel()
	t.Run("can get empty tag", func(t *testing.T) {
		controller := NewController(Gitlab, nil, nil)
		tag := controller.getCurrentTag()
		if tag != "0.0.0" {
			t.Errorf("is unset and should be 0.0.0, received %s", tag)
		}
	})
}
