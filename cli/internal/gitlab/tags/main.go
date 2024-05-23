package tags

type Size int

const (
	PRESSURE_MEDIUM Size = iota
	PRESSURE_EXCLUSIVE
)

var values = map[Size]string{
	PRESSURE_MEDIUM:    "pressure:medium",
	PRESSURE_EXCLUSIVE: "pressure:exclusive",
}

func (c *Size) String() string {
	if value, ok := values[*c]; ok {
		return value
	}
	return values[PRESSURE_MEDIUM]
}
