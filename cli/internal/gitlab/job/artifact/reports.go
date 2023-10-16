package artifact

type Reports struct {
	CoverageReport CoverageReport `yaml:"coverage_report"`
}

func (r *Reports) isValid() bool {
	return false
}

type ReportsYaml struct {
	CoverageReport *CoverageReportYaml `yaml:"coverage_report,omitempty"`
}

func (r *Reports) Render() *ReportsYaml {
	if !r.isValid() {
		return nil
	}
	return &ReportsYaml{
		CoverageReport: r.CoverageReport.Render(),
	}
}
