package artifact

type JunitReport struct {
	Value string
}

func (j *JunitReport) Render() string {
	return j.Value
}

func NewJunitReport(path string) JunitReport {
	return JunitReport{
		Value: path,
	}
}
