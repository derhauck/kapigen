package artifact

import "gitlab.com/kateops/kapigen/dsl/gitlab/job/artifact/reports"

type CoverageReport struct {
	CoverageFormat reports.CoverageFormat
	Path           string
}

func NewCoverageReport(coverageFormat reports.CoverageFormat, path string) CoverageReport {
	return CoverageReport{
		CoverageFormat: coverageFormat,
		Path:           path,
	}
}

func (c *CoverageReport) Render() *CoverageReportYaml {
	return NewCoverageReportYaml(c)
}

type CoverageReportYaml struct {
	CoverageFormat string `yaml:"coverage_format"`
	Path           string `yaml:"path"`
}

func NewCoverageReportYaml(coverageReport *CoverageReport) *CoverageReportYaml {
	if coverageReport.Path == "" {
		return nil
	}

	return &CoverageReportYaml{
		CoverageFormat: coverageReport.CoverageFormat.String(),
		Path:           coverageReport.Path,
	}
}
