package job

import (
	"gitlab.com/kateops/kapigen/dsl/logger"
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
	enum, err := wrapper.NewEnum[EnvironmentAction](map[EnvironmentAction]string{
		EnvironmentActionStart:   "start",
		EnvironmentActionPrepare: "prepare",
		EnvironmentActionStop:    "stop",
		EnvironmentActionVerify:  "verify",
		EnvironmentActionAccess:  "access",
	})
	if err != nil {
		logger.ErrorE(wrapper.DetailedErrorE(err))
		return nil
	}
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
	enum, err := wrapper.NewEnum[EnvironmentTier](map[EnvironmentTier]string{
		EnvironmentTierProduction: "production",
		EnvironmentTierStaging:    "staging",
		EnvironmentTierTesting:    "testing",
		EnvironmentTierDeployment: "development",
		EnvironmentTierOther:      "other",
	})
	if err != nil {
		logger.ErrorE(wrapper.DetailedErrorE(err))
		return nil
	}
	return enum
}

type Environment struct {
	Name           string
	Url            string
	OnStop         string
	Action         EnvironmentAction
	AutoStopIn     string
	DeploymentTier string
}

func (e Environment) Render() *EnvironmentYaml {
	return &EnvironmentYaml{
		Name:           e.Name,
		Url:            e.Url,
		OnStop:         e.OnStop,
		Action:         e.Action.String(),
		AutoStopIn:     e.AutoStopIn,
		DeploymentTier: e.DeploymentTier,
	}
}

type EnvironmentYaml struct {
	Name           string `yaml:"name,omitempty"`
	Url            string `yaml:"url,omitempty"`
	OnStop         string `yaml:"on_stop,omitempty"`
	Action         string `yaml:"action,omitempty"`
	AutoStopIn     string `yaml:"auto_stop_in"`
	DeploymentTier string `yaml:"deployment_tier,omitempty"`
}
