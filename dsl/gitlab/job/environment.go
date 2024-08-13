package job

import (
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

type EnvironmentAction byte

func (e EnvironmentAction) String() string {
	return EnvironmentActionEnum().ValueSafe(e)
}

const (
	EnvironmentActionStart EnvironmentAction = iota
	EnvironmentActionPrepare
	EnvironmentActionStop
	EnvironmentActionVerify
	EnvironmentActionAccess
)

func EnvironmentActionEnum() *wrapper.Enum[EnvironmentAction, string] {
	enum, _ := wrapper.NewEnum[EnvironmentAction](map[EnvironmentAction]string{
		EnvironmentActionStart:   "start",
		EnvironmentActionPrepare: "prepare",
		EnvironmentActionStop:    "stop",
		EnvironmentActionVerify:  "verify",
		EnvironmentActionAccess:  "access",
	})
	return enum
}

type EnvironmentTier byte

func (e EnvironmentTier) String() string {
	return EnvironmentTierEnum().ValueSafe(e)
}

const (
	EnvironmentTierProduction EnvironmentTier = iota
	EnvironmentTierStaging
	EnvironmentTierTesting
	EnvironmentTierDeployment
	EnvironmentTierOther
)

func EnvironmentTierEnum() *wrapper.Enum[EnvironmentTier, string] {
	enum, _ := wrapper.NewEnum[EnvironmentTier](map[EnvironmentTier]string{
		EnvironmentTierProduction: "production",
		EnvironmentTierStaging:    "staging",
		EnvironmentTierTesting:    "testing",
		EnvironmentTierDeployment: "development",
		EnvironmentTierOther:      "other",
	})
	return enum
}

type Environment struct {
	Name           string
	Url            string
	OnStop         wrapper.GetNamer
	Action         EnvironmentAction
	AutoStopIn     string
	DeploymentTier EnvironmentTier
}

func (e Environment) Render() *EnvironmentYaml {
	var onStop string
	if e.OnStop != nil {
		onStop = e.OnStop.GetName()
	}

	return &EnvironmentYaml{
		Name:           e.Name,
		Url:            e.Url,
		OnStop:         onStop,
		Action:         e.Action.String(),
		AutoStopIn:     e.AutoStopIn,
		DeploymentTier: e.DeploymentTier.String(),
	}
}

type EnvironmentYaml struct {
	Name           string `yaml:"name,omitempty"`
	Url            string `yaml:"url,omitempty"`
	OnStop         string `yaml:"on_stop,omitempty"`
	OnStopJob      *CiJob
	Action         string `yaml:"action,omitempty"`
	AutoStopIn     string `yaml:"auto_stop_in"`
	DeploymentTier string `yaml:"deployment_tier,omitempty"`
}
