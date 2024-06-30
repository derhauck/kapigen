package config

import (
	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/gitlab/job"
	artifact2 "kapigen.kateops.com/internal/gitlab/job/artifact"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/pipeline/wrapper"
	"kapigen.kateops.com/internal/when"
)

type Generic struct {
	ImageName       string             `yaml:"imageName"`
	Mode            string             `yaml:"mode"`
	Scripts         []string           `yaml:"scripts"`
	Variables       map[string]string  `yaml:"variables"`
	Stage           string             `yaml:"stage"`
	Artifacts       *job.ArtifactsYaml `yaml:"artifacts,omitempty"`
	Changes         []string           `yaml:"changes"`
	RuleSet         job.RulesYaml      `yaml:"rules"`
	InternalStage   stages.Stage
	InternalChanges []string
}

func (g *Generic) New() types.PipelineConfigInterface {
	return &Generic{}
}

func (g *Generic) Validate() error {

	if g.ImageName == "" {
		g.ImageName = docker.Alpine_3_18.String()
	}
	g.InternalStage, _ = stages.FromString(g.Stage)
	if len(g.Changes) == 0 {
		g.Changes = append(g.Changes, ".")
	}
	g.InternalChanges = g.Changes

	return nil
}

func (g *Generic) Build(_ *factory.MainFactory, _ types.PipelineType, _ string) (*types.Jobs, error) {
	var allJobs types.Jobs
	generic := types.NewJob("Generic Job", g.ImageName, func(ciJob *job.CiJob) {
		ciJob.SetStage(g.InternalStage).
			AddScripts(g.Scripts).
			TagMediumPressure()

		if g.Artifacts != nil {
			artifact := job.Artifacts{
				Name:    g.Artifacts.Name,
				Paths:   *wrapper.NewArray[string]().Push(g.Artifacts.Paths...),
				Reports: artifact2.Reports{},
			}
			ciJob.AddArtifact(artifact)
		}

		for key, value := range g.Variables {
			ciJob.AddVariable(key, value)
		}

	})
	return allJobs.AddJob(generic), nil
}

func (g *Generic) Rules() *job.Rules {
	var rules job.Rules
	for _, rule := range g.RuleSet {
		allowFailure := false
		if result := rule.AllowFailure.(bool); result {
			allowFailure = result
		}
		rules = append(rules, &job.Rule{
			If:           rule.If,
			Changes:      *wrapper.NewArray[string]().Push(rule.Changes...),
			AllowFailure: wrapper.Bool{Value: allowFailure},
			Variables:    rule.Variables,
			When:         job.NewWhen(when.OnSuccess),
		})
	}
	if len(g.RuleSet) > 0 {
		return &rules
	}
	return job.DefaultPipelineRules(g.InternalChanges)
}
