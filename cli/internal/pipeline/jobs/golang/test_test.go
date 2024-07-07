package golang

import (
	"slices"
	"strings"
	"testing"

	"gitlab.com/kateops/kapigen/dsl/enum"
)

var defaultCoveragePackages = []string{"./..."}
var defaultSource = "./..."

func TestCreate(t *testing.T) {
	t.Parallel()
	t.Run("can create job with correct parameters", func(t *testing.T) {
		job, err := NewUnitTest("golang:1.16", "test", defaultCoveragePackages, defaultSource)
		if job == nil && err != nil {
			t.Error("can not create job")
			t.Error(err)
		}
	})

	t.Run("can not create test with wrong parameters", func(t *testing.T) {
		job, err := NewUnitTest("", "", defaultCoveragePackages, defaultSource)
		if job != nil || err == nil {
			t.Error("created job succeeded without image and path but should not")
		}
		if err != nil && !strings.Contains(err.Error(), "no imageName set") {
			t.Error(err)
		}

		job, err = NewUnitTest("golang:1.16", "", defaultCoveragePackages, defaultSource)
		if job != nil {
			t.Error("created job succeeded without path but should not")
		}

		if err != nil && !strings.Contains(err.Error(), "no path set") {
			t.Error(err)
		}
	})

	t.Run("has correct parameters", func(t *testing.T) {
		expectedImageName := "golang:1.16"
		expectedPath := "test"
		job, _ := NewUnitTest(expectedImageName, expectedPath, defaultCoveragePackages, defaultSource)

		if !slices.Contains(job.Names, "Unit Test") {
			t.Error("job has wrong name")
		}
		for _, tag := range job.CiJob.Tags {
			medium := enum.TagPressureMedium
			if tag.Get() != medium.String() {
				t.Error("job has wrong tags")
			}
		}

		if job.CiJob.Image.Name != expectedImageName {
			t.Errorf("expected image name %s, got %s", expectedImageName, job.CiJob.Image.Name)
		}

		//if job.CiJob.Artifacts.Reports.CoverageReport.Path == "" {
		//	t.Error("missing coverage report path")
		//	t.Error(job.Render())
		//}
		//expectedReportPath := expectedPath + "/report.xml"
		//actualReportPath := job.CiJob.Artifacts.Reports.CoverageReport.Path
		//if actualReportPath != "coverage.html" {
		//	t.Errorf("expexted coverage report path %s, got %s", expectedReportPath, actualReportPath)
		//}

	})
}
