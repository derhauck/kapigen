package version

import (
	"testing"

	"gitlab.com/kateops/kapigen/dsl/environment"
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
		controller := NewController(Gitlab, nil, nil)
		tag := controller.GetCurrentTag("")
		if tag != defaultTag {
			t.Errorf("is unset and should be %s, received %s", defaultTag, tag)
		}

		controller = NewController(FILE, nil, NewFileReader())
		tag = controller.GetCurrentTag("")
		if tag != "" {
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

func TestController_GetCurrentPipelineTag(t *testing.T) {

	type args struct {
		settings func()
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Cat get tag from git tag",
			args{
				func() {
					environment.CI_COMMIT_TAG.Set("1.0.0")
				},
			},
			"1.0.0",
		},
		{
			"Cat use new tag - main",
			args{
				func() {
					environment.CI_COMMIT_TAG.Unset()
					environment.CI_COMMIT_BRANCH.Set("main")
					environment.CI_MERGE_REQUEST_ID.Unset()
					environment.CI_MERGE_REQUEST_SOURCE_BRANCH_NAME.Unset()
				},
			},
			"0.0.0-main",
		},
		{
			"Cat use new tag - feature",
			args{
				func() {
					environment.CI_COMMIT_TAG.Unset()
					environment.CI_COMMIT_BRANCH.Set("feature")
					environment.CI_MERGE_REQUEST_ID.Set("100")
					environment.CI_MERGE_REQUEST_SOURCE_BRANCH_NAME.Set("feature")
				},
			},
			"0.0.0-feature",
		},
		{
			"Cat use tag",
			args{
				func() {
					environment.CI_COMMIT_TAG.Set("1.0.0")
				},
			},
			"1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewController(
				Gitlab,
				nil,
				nil,
			)
			tt.args.settings()
			if got := c.GetCurrentPipelineTag(""); got != tt.want {
				t.Errorf("GetCurrentPipelineTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetModeFromString(t *testing.T) {
	tests := []struct {
		name   string
		args   string
		expect Mode
	}{
		{
			"Get File",
			"file",
			FILE,
		},
		{
			"Get Gitlab",
			"gitlab",
			Gitlab,
		},
		{
			"Get None",
			"none",
			None,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetModeFromString(tt.args); got != tt.expect {
				t.Errorf("GetModeFromString() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestController_GetNewTag(t *testing.T) {

	type args struct {
		settings func()
	}
	tests := []struct {
		name   string
		args   args
		expect string
	}{
		{
			"Get new Tag from 0.0.0 - major",
			args{
				settings: func() {
					environment.CI_PROJECT_ID.Set("10")
					environment.CI_MERGE_REQUEST_ID.Set("10")
					environment.CI_MERGE_REQUEST_LABELS.Set("version::major")
					environment.CI_MERGE_REQUEST_SOURCE_BRANCH_NAME.Set("feature")
					environment.CI_COMMIT_TAG.Unset()
				},
			},
			"1.0.0",
		},
		{
			"Get new Tag from 0.0.0 - minor",
			args{
				settings: func() {
					environment.CI_PROJECT_ID.Set("10")
					environment.CI_MERGE_REQUEST_ID.Set("10")
					environment.CI_MERGE_REQUEST_LABELS.Set("version::minor")
					environment.CI_MERGE_REQUEST_SOURCE_BRANCH_NAME.Set("feature")
					environment.CI_COMMIT_TAG.Unset()
				},
			},
			"0.1.0",
		},
		{
			"Get new Tag from 0.0.0 - patch",
			args{
				settings: func() {
					environment.CI_PROJECT_ID.Set("10")
					environment.CI_MERGE_REQUEST_ID.Set("10")
					environment.CI_MERGE_REQUEST_LABELS.Set("version::patch")
					environment.CI_MERGE_REQUEST_SOURCE_BRANCH_NAME.Set("feature")
					environment.CI_COMMIT_TAG.Unset()
				},
			},
			"0.0.1",
		},
	}
	for _, tt := range tests {
		tt.args.settings()
		c := NewController(
			Gitlab,
			nil,
			nil,
		)
		t.Run(tt.name, func(t *testing.T) {
			if got := c.GetNewTag(""); got != tt.expect {
				t.Errorf("GetNewTag() = %v, want %v", got, tt.expect)
			}
		})
	}
}
