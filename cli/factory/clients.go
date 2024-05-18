package factory

import (
	"github.com/xanzy/go-gitlab"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/logger"
)

type Clients struct {
	Gitlab *gitlab.Client
}

func (c *Clients) GetGitlabClient() *gitlab.Client {
	if c.Gitlab == nil {
		client, err := gitlab.NewClient(environment.CI_JOB_TOKEN.Get(), gitlab.WithBaseURL(environment.CI_SERVER_URL.Get()))
		if err != nil {
			logger.Error("could not create gitlab client")
			return nil
		}
		c.Gitlab = client
	}
	return c.Gitlab
}
