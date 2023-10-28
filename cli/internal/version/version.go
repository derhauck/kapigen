package version

import (
	"github.com/xanzy/go-gitlab"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/los"
)

const emptyTag = "0.0.0"

type Mode int

const (
	Gitlab Mode = iota
	Los
	None
)

func (m Mode) getTag() string {
	switch m {
	case Gitlab:
		return Gitlab.Name()
	case Los:
		return Los.Name()
	default:
		return None.Name()

	}
}

var values = map[Mode]string{
	Gitlab: "Gitlab",
	Los:    "Logic Operator Server",
	None:   "No versioning",
}

func (m Mode) Name() string {
	if value, ok := values[m]; ok {
		return value
	}
	logger.Error("no mode found, use default")
	m = Gitlab
	return values[m]
}

type Controller struct {
	current      string
	intermediate string
	mode         Mode
	gitlabClient *gitlab.Client
	losClient    *los.Client
	refresh      bool
}

func NewController(mode Mode, gitlabClient *gitlab.Client, losClient *los.Client) *Controller {
	return &Controller{
		"",
		"",
		mode,
		gitlabClient,
		losClient,
		false,
	}
}

func (c *Controller) getTagFromGitlab() string {
	oderBy := "updated"
	if c.gitlabClient == nil {
		logger.Error("no gitlab client configured")
		return emptyTag
	}
	tags, _, err := c.gitlabClient.Tags.ListTags(environment.CI_PROJECT_ID.Get(), &gitlab.ListTagsOptions{OrderBy: &oderBy})
	if err != nil {
		logger.ErrorE(err)
		return emptyTag
	}
	logger.DebugAny(tags)
	return tags[0].Name
}

func (c *Controller) getTagFromLos(path string) string {
	if c.losClient == nil {
		return emptyTag
	}
	return c.losClient.GetLatestVersion(environment.CI_PROJECT_ID.Get(), path)
}
func (c *Controller) Refresh() *Controller {
	c.refresh = true
	return c
}
func (c *Controller) getCurrentTag(path string) string {
	if c.current == "" || c.refresh {
		c.refresh = false
		if c.mode == Gitlab {
			c.current = c.getTagFromGitlab()
		}
		if c.mode == Los {
			c.current = c.getTagFromLos(path)
		}
	}
	return c.current
}
