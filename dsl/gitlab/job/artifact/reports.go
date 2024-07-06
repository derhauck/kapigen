package artifact

type Reports struct {
	CoverageReport CoverageReport `yaml:"coverage_report,omitempty"`
	Junit          JunitReport    `yaml:"junit,omitempty"`
}

func NewReports() Reports {
	return Reports{}
}

func (r Reports) SetCoverageReport(coverageReport CoverageReport) Reports {
	r.CoverageReport = coverageReport
	return r
}

func (r Reports) SetJunit(junit JunitReport) Reports {
	r.Junit = junit
	return r
}

func (r Reports) isValid() bool {
	return true
}

type ReportsYaml struct {
	CoverageReport *CoverageReportYaml `yaml:"coverage_report,omitempty"`
	Junit          string              `yaml:"junit,omitempty"`
}

func (r Reports) Render() *ReportsYaml {
	if !r.isValid() {
		return nil
	}
	return &ReportsYaml{
		CoverageReport: r.CoverageReport.Render(),
		Junit:          r.Junit.Render(),
	}
}
