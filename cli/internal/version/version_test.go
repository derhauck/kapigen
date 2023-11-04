package version

import (
	"github.com/xanzy/go-gitlab"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/los"
	"testing"
)

const GitlabUrl = "https://gitlab.kateops.com"

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
	environment.SetLocalEnv()
	var defaultTag = "0.0.0"
	t.Run("can get empty tag", func(t *testing.T) {
		controller := NewController(Gitlab, nil, nil)
		tag := controller.GetCurrentTag("")
		if tag != defaultTag {
			t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
		}

		controller = NewController(Los, nil, nil)
		tag = controller.GetCurrentTag("")
		if tag != defaultTag {
			t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
		}
	})

	t.Run("can get default tag", func(t *testing.T) {
		gitlabClient, _ := gitlab.NewClient("", gitlab.WithBaseURL(GitlabUrl))
		controller := NewController(Gitlab, gitlabClient, los.NewClient(los.LosHostName, ""))
		tag := controller.GetCurrentTag("")
		if tag != defaultTag {
			t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
		}

		controller = NewController(Los, gitlabClient, los.NewClient(los.LosHostName, ""))
		tag = controller.GetCurrentTag("")
		if tag != defaultTag {
			t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
		}
	})

	t.Run("can refresh tag", func(t *testing.T) {
		gitlabClient, _ := gitlab.NewClient("", gitlab.WithBaseURL(GitlabUrl))
		losClient := los.NewClient(los.LosHostName, "")
		expectedTag := "1.0.0"
		controller := &Controller{
			expectedTag,
			"",
			"",
			Gitlab,
			gitlabClient,
			losClient,
			false,
		}
		tag := controller.GetCurrentTag("")
		if tag != expectedTag {
			t.Errorf("is set initially and should be %s, received %s", expectedTag, tag)
		}

		tag = controller.Refresh().GetCurrentTag("")
		if tag != defaultTag {
			t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
		}

		controller = &Controller{
			expectedTag,
			"",
			"",
			Los,
			gitlabClient,
			losClient,
			false,
		}
		tag = controller.GetCurrentTag("")
		if tag != expectedTag {
			t.Errorf("is set initially and should be %s, received %s", expectedTag, tag)
		}

		tag = controller.Refresh().GetCurrentTag("")
		if tag != defaultTag {
			t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
		}
	})
}
