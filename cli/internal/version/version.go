package version

import (
	"strings"

	"github.com/xanzy/go-gitlab"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/logger"
)

const EmptyTag = "0.0.0"

const NoTag = "latest"

type Mode int

const (
	Gitlab Mode = iota
	FILE
	None
)

func (m Mode) getTag() string {
	switch m {
	case Gitlab:
		return Gitlab.Name()
	default:
		return None.Name()

	}
}

var values = map[Mode]string{
	Gitlab: "gitlab",
	FILE:   "file",
	None:   "none",
}

func (m Mode) Name() string {
	if value, ok := values[m]; ok {
		return value
	}
	logger.Error("no mode found, use default")
	return values[Gitlab]
}

type Controller struct {
	current      string
	intermediate string
	new          string
	mode         Mode
	gitlabClient *gitlab.Client
	FileReader   *FileReader
	refresh      bool
}

func NewController(mode Mode, gitlabClient *gitlab.Client, reader *FileReader) *Controller {
	return &Controller{
		"",
		"",
		"",
		mode,
		gitlabClient,
		reader,
		false,
	}
}

func (c *Controller) getTagFromGitlab() string {
	if c.refresh == false && c.current != "" {
		return c.current
	}
	oderBy := "updated"
	sort := "desc"
	if c.gitlabClient == nil {
		logger.Error("no gitlab client configured")
		return EmptyTag
	}
	tags, _, err := c.gitlabClient.Tags.ListTags(environment.CI_PROJECT_ID.Get(), &gitlab.ListTagsOptions{OrderBy: &oderBy, Sort: &sort})
	if err != nil {
		logger.ErrorE(err)
		return EmptyTag
	}

	if len(tags) == 0 {
		return EmptyTag
	}

	return tags[0].Name
}

func (c *Controller) createTagFromGitlab(version string) string {
	ref := environment.CI_DEFAULT_BRANCH.Get()
	msg := "Kapigen auto increment new tag"
	if c.gitlabClient == nil {
		logger.Error("no gitlab client configured")
		return EmptyTag
	}

	release, _, err := c.gitlabClient.Releases.CreateRelease(environment.CI_PROJECT_ID.Get(), &gitlab.CreateReleaseOptions{
		TagName:    &version,
		Ref:        &ref,
		TagMessage: &msg,
	})
	if err != nil {
		return ""
	}
	if err != nil {
		logger.ErrorE(err)
		return EmptyTag
	}
	logger.DebugAny(release)
	return version
}
func (c *Controller) Refresh() *Controller {
	c.refresh = true
	return c
}
func (c *Controller) GetCurrentTag(path string) string {
	if c.current == "" || c.refresh {
		switch c.mode {
		case Gitlab:
			c.current = c.getTagFromGitlab()
		case None:
			c.current = EmptyTag
		case FILE:
			c.current = c.FileReader.Get(path)
		default:
			c.current = EmptyTag
		}
	}

	c.refresh = false
	return c.current
}

func (c *Controller) GetCurrentPipelineTag(path string) string {
	if environment.IsRelease() {
		return environment.CI_COMMIT_TAG.Get()
	}

	return c.GetIntermediateTag(path)
}

func (c *Controller) GetNewTag(path string) string {
	if c.new == "" || c.refresh {
		switch c.mode {
		case Gitlab:
			c.new = getNewVersion(
				c.GetCurrentTag(path),
				c.getVersionIncrease(environment.CI_PROJECT_ID.Get(), environment.GetMergeRequestId()),
			)
		case FILE:
			c.new = NoTag
		case None:
			c.new = NoTag

		}
	}
	c.refresh = false
	return c.new
}

func (c *Controller) GetIntermediateTag(path string) string {
	if c.intermediate == "" || c.refresh {
		switch c.mode {
		case Gitlab:
			c.intermediate = GetFeatureBranchVersion(c.GetCurrentTag(path), environment.GetBranchName())
		case FILE:
			c.intermediate = EmptyTag
		case None:
			c.intermediate = NoTag

		}
	}
	c.refresh = false
	return c.intermediate
}

func (c *Controller) SetNewVersion(path string) string {
	if c.new == "" || c.refresh {
		c.GetNewTag(path)
	}
	switch c.mode {
	case Gitlab:
		return c.createTagFromGitlab(c.new)
	case FILE:
		return c.new
	case None:
		return c.new
	default:
		return c.new
	}
}

func (c *Controller) getVersionIncrease(projectId string, mrId int) string {
	if _, err := environment.CI_MERGE_REQUEST_ID.Lookup(); err != nil {
		return getVersionIncreaseFromLabels(
			c.getMrLabelsFromApi(projectId, mrId),
		)
	}
	return getVersionIncreaseFromLabels(environment.CI_MERGE_REQUEST_LABELS.Get())
}

func (c *Controller) getMrLabelsFromApi(projectId string, mrId int) string {
	mr, _, err := c.gitlabClient.MergeRequests.GetMergeRequest(projectId, mrId, nil)
	if err != nil {
		logger.ErrorE(err)
		return "none"
	}
	return strings.Join(mr.Labels, ",")
}

func GetModeFromString(mode string) Mode {
	switch mode {
	case Gitlab.Name():
		return Gitlab
	default:
		return None
	}
}
