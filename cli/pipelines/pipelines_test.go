package pipelines

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
	"gitlab.com/kateops/kapigen/dsl/gitlab/pipeline"
	"gitlab.com/kateops/kapigen/dsl/gitlab/stages"
	"gopkg.in/yaml.v3"
)

func TestCreatePipeline(t *testing.T) {
	t.Run("can create pipeline", func(t *testing.T) {
		file := "pipeline.yaml"
		CreatePipeline(func(jobs *types.Jobs, ciPipeline *pipeline.CiPipeline) {
			ciPipeline.DefaultCiPipeline()
			jobs.AddJob(types.NewJob("generic", "alpine", func(ciJob *job.CiJob) {
				ciJob.TagMediumPressure().
					AddScript("echo hello world").
					SetStage(stages.TEST).
					Rules.AddRules(*job.DefaultMainBranchRules())
				ciJob.Secrets = job.Secrets{
					"TEST": &job.VaultSecret{
						Vault: job.VaultSecretConfig{
							Engine: job.VaultSecretEngine{
								Name: job.EnumVaultSecretEngineKv2,
								Path: "mount",
							},
							Path:  "path",
							Field: "field",
						},
						Token: "token",
					},
				}
			}))
		})
		readFile, err := os.ReadFile(file)
		if err != nil {
			t.Error(err)
		}
		pipelineConfig := map[string]any{}
		err = yaml.NewDecoder(bytes.NewReader(readFile)).Decode(&pipelineConfig)
		if err != nil {
			t.Error(err)
		}
		snaps.MatchSnapshot(t, pipelineConfig["generic"], pipelineConfig["variables"], pipelineConfig["workflow"], pipelineConfig["default"])
		snaps.MatchSnapshot(t, string(readFile))
		err = os.Remove(file)
		if err != nil {
			t.Error(err)
		}

	})
	t.Run("can not create pipeline", func(t *testing.T) {
		file := "pipeline.yaml"
		_ = os.Remove(file)
		CreatePipeline(func(jobs *types.Jobs, ciPipeline *pipeline.CiPipeline) {
			jobs.AddJob(types.NewJob("invalid", "alpine", func(ciJob *job.CiJob) {
			}))
		})
		_, err := os.ReadFile(file)
		if err == nil {
			t.Errorf("should not be able to open %s", file)
			_ = os.Remove(file)
		}

		if err.Error() != "open pipeline.yaml: no such file or directory" {
			t.Errorf("expectec: %s, received: %s", err, "open pipeline.yaml: no such file or directory")
		}

	})
}

func TestReadPipelineConfig(t *testing.T) {
	t.Run("can create pipelineConfig", func(t *testing.T) {
		configPath := "../config.kapigen.yaml"
		config, err := ReadPipelineConfig(configPath)
		if err != nil {
			t.Error(err)
		}
		if config.Noop == false {
			t.Error("should be true")
		}

		if config.Versioning == false {
			t.Error("should be true")
		}

		if config.PrivateTokenName != "CI_PIPELINE_TOKEN" {
			t.Error("should be CI_PIPELINE_TOKEN")
		}

	})
	t.Run("can not create pipelineConfig", func(t *testing.T) {
		configPath := "config.kapigen.yaml-missing"
		config, err := ReadPipelineConfig(configPath)
		expectedError := "open config.kapigen.yaml-missing: no such file or directory"
		if config != nil {
			t.Error("config should be nil")
		}
		if err == nil {
			t.Error("should contain error")
			t.FailNow()
		}
		if !strings.Contains(err.Error(), expectedError) {
			t.Errorf("expected: '%s', received: '%s'", expectedError, err)
		}

	})
}
