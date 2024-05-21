package factory

import (
	"os"

	"github.com/xanzy/go-gitlab"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/logger"
)

type Clients struct {
	Gitlab *gitlab.Client
}

func (c *Clients) GetGitlabJobClient() *gitlab.Client {
	if c.Gitlab == nil {
		client, err := gitlab.NewJobClient(environment.CI_JOB_TOKEN.Get(), gitlab.WithBaseURL(environment.CI_SERVER_URL.Get()))
		if err != nil {
			logger.Error("could not create gitlab client")
			return nil
		}
		c.Gitlab = client
	}
	return c.Gitlab
}

func (c *Clients) GetGitlabClient(privateToken string) *gitlab.Client {
	if c.Gitlab == nil {
		client, err := gitlab.NewClient(os.Getenv(privateToken), gitlab.WithBaseURL(environment.CI_SERVER_URL.Get()))
		if err != nil {
			logger.Error("could not create gitlab client")
			return nil
		}
		c.Gitlab = client
	}
	return c.Gitlab
}
