package factory

import (
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
func (m *MainFactory) GetVersionController() *version.Controller {
	if m.version == nil {
		switch m.Settings.Mode {
		case version.Gitlab:
			m.version = version.NewController(
				m.Settings.Mode,
				m.GetClients().GetGitlabClient(),
				nil,
			)
		case version.Los:
			m.version = version.NewController(
				m.Settings.Mode,
				nil,
				m.GetClients().GetLosClient(),
			)
		case version.None:
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
