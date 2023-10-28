package factory

import "kapigen.kateops.com/internal/version"

type MainFactory struct {
	clients *Clients
	version *version.Controller
}

func New() *MainFactory {
	return &MainFactory{}
}
func (m *MainFactory) GetVersion() *version.Controller {
	if m.version == nil {
		m.version = m.newVersionController(version.Gitlab)
	}
	return m.version
}

func (m *MainFactory) GetClients() *Clients {
	if m.clients == nil {
		m.clients = &Clients{}
	}
	return m.clients
}
