package factory

import (
	"github.com/xanzy/go-gitlab"
	"kapigen.kateops.com/internal/cli"
	"kapigen.kateops.com/internal/version"
)

type MainFactory struct {
	Settings *cli.Settings
	clients  *Clients
	version  *version.Controller
}

func New(settings *cli.Settings) *MainFactory {
	return &MainFactory{
		Settings: settings,
	}
}

func (m *MainFactory) GetGitlabClient() *gitlab.Client {
	if m.Settings.PrivateToken == "" {
		return m.GetClients().GetGitlabJobClient()
	} else {
		return m.GetClients().GetGitlabClient(m.Settings.PrivateToken)
	}
}
func (m *MainFactory) GetVersionController() *version.Controller {
	if m.version == nil {
		switch m.Settings.Mode {
		case version.Gitlab:
			m.version = version.NewController(
				m.Settings.Mode,
				m.GetGitlabClient(),
				version.NewFileReader(),
			)
		case version.FILE:
			m.version = version.NewController(
				m.Settings.Mode,
				nil,
				version.NewFileReader(),
			)
		case version.None:
			m.version = version.NewController(
				m.Settings.Mode,
				nil,
				nil,
			)
		default:
			m.version = version.NewController(
				m.Settings.Mode,
				nil,
				nil,
			)
		}

	}
	return m.version
}

func (m *MainFactory) GetClients() *Clients {
	if m.clients == nil {
		m.clients = &Clients{}
	}
	return m.clients
}
