package pipelines

import (
	"bytes"
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
	"gitlab.com/kateops/kapigen/dsl/gitlab/stages"
	"gopkg.in/yaml.v3"
)

func TestCreatePipeline(t *testing.T) {
	t.Run("can create pipeline", func(t *testing.T) {
		file := "pipeline.yaml"
		CreatePipeline(func(jobs *types.Jobs) {
			jobs.AddJob(types.NewJob("generic", "alpine", func(ciJob *job.CiJob) {
				ciJob.TagMediumPressure().
					AddScript("echo hello world").
					SetStage(stages.TEST).
					Rules.AddRules(*job.DefaultMainBranchRules())
			}))
		})
		readFile, err := os.ReadFile(file)
		if err != nil {
			t.Error(err)
		}
		pipeline := map[string]any{}
		err = yaml.NewDecoder(bytes.NewReader(readFile)).Decode(&pipeline)
		if err != nil {
			t.Error(err)
		}
		snaps.MatchSnapshot(t, pipeline["generic"], pipeline["variables"])

		err = os.Remove(file)
		if err != nil {
			t.Error(err)
		}

	})
	t.Run("can not create pipeline", func(t *testing.T) {
		file := "pipeline.yaml"
		_ = os.Remove(file)
		CreatePipeline(func(jobs *types.Jobs) {
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
