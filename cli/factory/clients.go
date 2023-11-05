package factory

import (
	gitlab "github.com/xanzy/go-gitlab"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/los"
)

const GitlabUrl = "https://gitlab.kateops.com"

type Clients struct {
	Gitlab *gitlab.Client
	Los    *los.Client
}

func (c *Clients) GetLosClient() *los.Client {
	if c.Los == nil {
		client := los.NewClient(los.LosHostName, environment.LOS_AUTH_TOKEN.Get())
		if client == nil {
			logger.Error("could not create LOS client")
		}
		c.Los = client
	}
	return c.Los
}

func (c *Clients) GetGitlabClient() *gitlab.Client {
	if c.Gitlab == nil {
		client, err := gitlab.NewClient(environment.CI_PIPELINE_TOKEN.Get(), gitlab.WithBaseURL(GitlabUrl))
		if err != nil {
			logger.Error("could not create gitlab client")
			return nil
		}
		c.Gitlab = client
	}
	return c.Gitlab
}
