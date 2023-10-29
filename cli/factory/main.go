package factory

import "kapigen.kateops.com/internal/version"

type MainFactory struct {
	clients *Clients
	version *version.Controller
}

func New() *MainFactory {
	return &MainFactory{}
}
func (m *MainFactory) GetVersionController() *version.Controller {
	if m.version == nil {
		m.version = version.NewController(
			version.Gitlab,
			m.GetClients().GetGitlabClient(),
			m.GetClients().GetLosClient(),
		)
	}
	return m.version
}

func (m *MainFactory) GetClients() *Clients {
	if m.clients == nil {
		m.clients = &Clients{}
	}
	return m.clients
}
