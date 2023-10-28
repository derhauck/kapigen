package environment

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/xanzy/go-gitlab"
	"kapigen.kateops.com/internal/logger"
	"regexp"
	"strings"
)

func getVersionIncrease() string {
	mergeLabels := CI_MERGE_REQUEST_LABELS.Get()
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

func GetTagFromGitlab(client *gitlab.Client) string {
	oderBy := "updated"
	tags, _, err := client.Tags.ListTags(CI_PROJECT_ID.Get(), &gitlab.ListTagsOptions{OrderBy: &oderBy})
	if err != nil {
		logger.ErrorE(err)
		return "0.0.0"
	}
	logger.DebugAny(tags)
	return tags[0].Name
}

func GetNewVersion(version string) string {
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
	branch := CI_MERGE_REQUEST_SOURCE_BRANCH_NAME.Get()
	reg := regexp.MustCompile("[/ ]")
	noEmptyOrSlash := reg.ReplaceAllString(branch, "-")
	reg = regexp.MustCompile("[!@#$%^&*()_+\\\\[\\]<>|.,;:'\"]")
	finalBranch := reg.ReplaceAllString(noEmptyOrSlash, "")
	return fmt.Sprintf("%s-%s", tag, finalBranch)
}
