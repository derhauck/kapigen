package version

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"gitlab.com/kateops/kapigen/dsl/logger"
)

func GetFeatureBranchVersion(tag string, branch string) string {
	reg := regexp.MustCompile("[/ ]")
	noEmptyOrSlash := reg.ReplaceAllString(branch, "-")
	reg = regexp.MustCompile("[!@#$%^&*()_+\\\\[\\]<>|.,;:'\"]")
	finalBranch := reg.ReplaceAllString(noEmptyOrSlash, "")
	return fmt.Sprintf("%s-%s", tag, finalBranch)
}

func getVersionIncreaseFromLabels(mergeLabels string) string {
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

func getNewVersion(version string, increase string) string {
	tag, err := semver.NewVersion(version)
	if err != nil {
		logger.ErrorE(err)
		return "0.0.0"
	}
	switch increase {
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
