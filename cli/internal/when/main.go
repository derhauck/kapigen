package when

import (
	"fmt"
	"kapigen.kateops.com/internal/logger"
)

type When int

const (
	OnSuccess When = iota
	OnFailure
	Always
	Never
)

var values = map[When]string{
	OnSuccess: "on_success",
	OnFailure: "on_failure",
	Always:    "always",
	Never:     "never",
}

func (c When) String() string {
	if v, ok := values[c]; ok {
		return v
	}
	logger.Error(fmt.Sprintf("when id: %d not found", c))
	return values[OnSuccess]
}
