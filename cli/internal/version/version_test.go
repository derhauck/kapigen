package version

import (
	"github.com/xanzy/go-gitlab"
	"kapigen.kateops.com/internal/los"
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
	var defaultTag = "0.0.0"
	t.Run("can get empty tag", func(t *testing.T) {
		controller := NewController(Gitlab, nil, nil)
		tag := controller.getCurrentTag("")
		if tag != defaultTag {
			t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
		}

		controller = NewController(Los, nil, nil)
		tag = controller.getCurrentTag("")
		if tag != defaultTag {
			t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
		}
	})

	t.Run("can get default tag", func(t *testing.T) {
		gitlabClient, _ := gitlab.NewClient("")
		controller := NewController(Gitlab, gitlabClient, los.NewClient(los.LosHostName, ""))
		tag := controller.getCurrentTag("")
		if tag != defaultTag {
			t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
		}

		controller = NewController(Los, gitlabClient, los.NewClient(los.LosHostName, ""))
		tag = controller.getCurrentTag("")
		if tag != defaultTag {
			t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
		}
	})
}
