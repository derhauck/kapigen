package reports

type CoverageFormat int

const (
	Cobertura CoverageFormat = iota
)

var values = map[CoverageFormat]string{
	Cobertura: "cobertura",
}

func (v *CoverageFormat) String() string {
	if s, ok := values[*v]; ok {
		return s
	}

	return values[Cobertura]
}
