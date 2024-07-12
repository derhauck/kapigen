package pipelines_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"gitlab.com/kateops/kapigen/cli/factory"
	"gitlab.com/kateops/kapigen/cli/internal/cli"
	"gitlab.com/kateops/kapigen/cli/internal/pipeline/config"
	"gitlab.com/kateops/kapigen/cli/internal/version"
	"gitlab.com/kateops/kapigen/cli/pipelines"
	"gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
	"gitlab.com/kateops/kapigen/dsl/gitlab/stages"
)

type TestPipelineConfig struct {
	ConfigId string `yaml:"id"`
}

func (t *TestPipelineConfig) New() types.PipelineConfigInterface {
	return &TestPipelineConfig{}
}
func (t *TestPipelineConfig) Validate() error {
	return nil
}

func (t *TestPipelineConfig) Build(_ *factory.MainFactory, _ types.PipelineType, _ string) (*types.Jobs, error) {
	return &types.Jobs{
		types.NewJob("test", "alpine", func(ciJob *job.CiJob) {
			ciJob.TagMediumPressure().
				AddScript("echo 'hello world'").
				SetStage(stages.TEST)
		}),
	}, nil
}

func (t *TestPipelineConfig) Rules() *job.Rules {
	return job.DefaultMainBranchRules()
}

func TestExtendPipelines(t *testing.T) {
	t.Run("extend pipeline functions", func(t *testing.T) {
		const TestPipeline types.PipelineType = "test"

		pipelines.ExtendPipelines(map[types.PipelineType]types.PipelineConfigInterface{
			TestPipeline: &TestPipelineConfig{},
		})
		settings := cli.NewSettings(
			cli.SetMode(version.GetModeFromString(version.Gitlab.Name())),
		)

		pipelineConfig := &types.PipelineConfig{
			Pipelines: []types.PipelineTypeConfig{
				{
					Type: TestPipeline,
					Config: TestPipelineConfig{
						ConfigId: "configId",
					},
					PipelineId: "testId",
				},
			},
		}

		pipelineJobs, err := types.LoadJobsFromPipelineConfig(factory.New(settings), pipelineConfig, config.PipelineConfigTypes)
		if err != nil {
			t.Error(err)
		}
		if len(pipelineJobs.GetJobs()) == 0 {
			t.Error("no jobs found")
		}

		snaps.MatchSnapshot(t, pipelineJobs)

	})
}
