package version

import (
	"github.com/xanzy/go-gitlab"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/los"
)

type Mode int

const (
	Gitlab Mode = iota
	Los
	None
)

var values = map[Mode]string{
	Gitlab: "Gitlab",
	Los:    "Logic Operator Server",
	None:   "No versioning",
}

func (m *Mode) Name() string {
	if value, ok := values[*m]; ok {
		return value
	}
	logger.Error("no mode found, use default")
	*m = Gitlab
	return values[*m]
}

type Controller struct {
	current      string
	intermediate string
	mode         Mode
	gitlabClient *gitlab.Client
	losClient    *los.Client
}

func NewController(current string, intermediate string, mode Mode, gitlabClient *gitlab.Client, losClient *los.Client) *Controller {
	return &Controller{
		current,
		intermediate,
		mode,
		gitlabClient,
		losClient,
	}
}
