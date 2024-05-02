package version

import (
	"testing"

	"kapigen.kateops.com/internal/environment"
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
	environment.SetLocalEnv()
	var defaultTag = "0.0.0"
	t.Run("can get empty tag", func(t *testing.T) {
		controller := NewController(Gitlab, nil)
		tag := controller.GetCurrentTag("")
		if tag != defaultTag {
			t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
		}

		controller = NewController(FILE, nil)
		tag = controller.GetCurrentTag("")
		if tag != defaultTag {
			t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
		}
	})

	//t.Run("can get default tag", func(t *testing.T) {
	//	gitlabClient, _ := gitlab.NewClient("", gitlab.WithBaseURL(GitlabUrl))
	//	controller := NewController(Gitlab, gitlabClient)
	//	tag := controller.GetCurrentTag("")
	//	if tag != defaultTag {
	//		t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
	//	}
	//
	//	controller = NewController(Los, gitlabClient)
	//	tag = controller.GetCurrentTag("")
	//	if tag != defaultTag {
	//		t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
	//	}
	//})
	//
	//t.Run("can refresh tag", func(t *testing.T) {
	//	gitlabClient, _ := gitlab.NewClient("", gitlab.WithBaseURL(GitlabUrl))
	//	expectedTag := "1.0.0"
	//	controller := &Controller{
	//		expectedTag,
	//		"",
	//		"",
	//		Gitlab,
	//		gitlabClient,
	//		false,
	//	}
	//	tag := controller.GetCurrentTag("")
	//	if tag != expectedTag {
	//		t.Errorf("is set initially and should be %s, received %s", expectedTag, tag)
	//	}
	//
	//	tag = controller.Refresh().GetCurrentTag("")
	//	if tag != defaultTag {
	//		t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
	//	}
	//
	//	controller = &Controller{
	//		expectedTag,
	//		"",
	//		"",
	//		Los,
	//		gitlabClient,
	//		false,
	//	}
	//	tag = controller.GetCurrentTag("")
	//	if tag != expectedTag {
	//		t.Errorf("is set initially and should be %s, received %s", expectedTag, tag)
	//	}
	//
	//	tag = controller.Refresh().GetCurrentTag("")
	//	if tag != defaultTag {
	//		t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
	//	}
	//})
}
