package tags

type Size int

const (
	PRESSURE_MEDIUM Size = iota
	PRESSURE_BUILDKIT
)

var values = map[Size]string{
	PRESSURE_MEDIUM:   "pressure::medium",
	PRESSURE_BUILDKIT: "pressure::buildkit",
}

func (c *Size) String() string {
	if value, ok := values[*c]; ok {
		return value
	}
	return values[PRESSURE_MEDIUM]
}
