package factory

import "kapigen.kateops.com/internal/version"

func (m *MainFactory) newVersionController(mode version.Mode) *version.Controller {
	return version.NewController(
		"",
		"",
		mode,
		m.clients.GetGitlabClient(),
		m.clients.GetLosClient(),
	)
}
