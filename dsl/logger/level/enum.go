package level

import "fmt"

type Level int

const (
	Error Level = iota
	Warning
	Info
	Debug
)

var values = map[Level]string{
	Error:   "ERROR",
	Warning: "WARNING",
	Info:    "INFO",
	Debug:   "DEBUG",
}

func (l Level) String() string {
	if v, ok := values[l]; ok {
		return v
	}
	return values[Error]
}

func FromString(s string) (Level, error) {
	for k, v := range values {
		if v == s {
			return k, nil
		}
	}
	return Error, fmt.Errorf("invalid level string: %s", s)
}

func (l Level) IsActive() bool {
	return l <= GetCurrentLevel()
}
