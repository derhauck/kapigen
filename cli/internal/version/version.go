package version

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/xanzy/go-gitlab"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/los"
	"regexp"
	"strings"
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
	new          string
	mode         Mode
	gitlabClient *gitlab.Client
	losClient    *los.Client
	refresh      bool
}

func NewController(mode Mode, gitlabClient *gitlab.Client, losClient *los.Client) *Controller {
	return &Controller{
		"",
		"",
		"",
		mode,
		gitlabClient,
		losClient,
		false,
	}
}

func (c *Controller) getTagFromGitlab() string {
	if c.refresh == false && c.current != "" {
		return c.current
	}
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
	if c.refresh == false && c.current != "" {
		return c.current
	}
	if c.losClient == nil {
		return emptyTag
	}
	return c.losClient.GetLatestVersion(environment.CI_PROJECT_ID.Get(), path)
}
func (c *Controller) Refresh() *Controller {
	c.refresh = true
	return c
}
func (c *Controller) GetCurrentTag(path string) string {
	if c.mode == Gitlab {
		c.current = c.getTagFromGitlab()
	}
	if c.mode == Los {
		c.current = c.getTagFromLos(path)
	}
	c.refresh = false
	return c.current
}

func (c *Controller) GetNewTag(path string) string {
	if c.new == "" || c.refresh {
		c.new = getNewVersion(c.GetCurrentTag(path))
	}
	c.refresh = false
	return c.new
}

func (c *Controller) GetIntermediateTag(path string) string {
	if c.intermediate == "" || c.refresh {
		c.intermediate = GetFeatureBranchVersion(c.GetCurrentTag(path))
	}
	c.refresh = false
	return c.intermediate
}
func getNewVersion(version string) string {
	tag, err := semver.NewVersion(version)
	if err != nil {
		logger.ErrorE(err)
		return "0.0.0"
	}
	switch getVersionIncrease() {
	case "major":
		return tag.IncMajor().String()
	case "minor":
		return tag.IncMinor().String()
	case "patch":
		return tag.IncPatch().String()
	default:
		logger.Error("no version increase found")
		return "0.0.0"

	}
}

func GetFeatureBranchVersion(tag string) string {
	branch := environment.CI_MERGE_REQUEST_SOURCE_BRANCH_NAME.Get()
	reg := regexp.MustCompile("[/ ]")
	noEmptyOrSlash := reg.ReplaceAllString(branch, "-")
	reg = regexp.MustCompile("[!@#$%^&*()_+\\\\[\\]<>|.,;:'\"]")
	finalBranch := reg.ReplaceAllString(noEmptyOrSlash, "")
	return fmt.Sprintf("%s-%s", tag, finalBranch)
}

func getVersionIncrease() string {
	mergeLabels := environment.CI_MERGE_REQUEST_LABELS.Get()
	labels := strings.Split(mergeLabels, ",")
	for _, label := range labels {
		if strings.Contains(label, "version::") {
			versionLabel := strings.Split(label, "::")
			if len(versionLabel) == 2 {
				return versionLabel[1]
			}
		}
	}
	logger.Error("no version increase found")
	return "none"
}
