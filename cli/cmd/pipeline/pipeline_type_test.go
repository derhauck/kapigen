package pipeline

import (
	"reflect"
	"strings"
	"testing"

	"gitlab.com/kateops/kapigen/cli/factory"
	"gitlab.com/kateops/kapigen/cli/internal/cli"
	"gitlab.com/kateops/kapigen/cli/internal/pipeline/config"
	"gitlab.com/kateops/kapigen/cli/internal/pipeline/jobs/docker"
	"gitlab.com/kateops/kapigen/cli/internal/pipeline/jobs/php"
	"gitlab.com/kateops/kapigen/cli/internal/version"
	types2 "gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/environment"
)

func TestLoadJobsFromPipelineConfig(t *testing.T) {
	type args struct {
		factory        *factory.MainFactory
		pipelineConfig types2.PipelineConfig
		configTypes    map[types2.PipelineType]types2.PipelineConfigInterface
	}

	environment.CI_COMMIT_BRANCH.Set("feature/test")
	environment.CI_DEFAULT_BRANCH.Set("main")
	environment.CI_MERGE_REQUEST_LABELS.Set("version::patch")
	environment.CI_MERGE_REQUEST_SOURCE_BRANCH_NAME.Set("feature/test")
	environment.CI_SERVER_URL.Set("https://gitlab.com")
	environment.CI_PROJECT_DIR.Set("/app")
	environment.CI_COMMIT_TAG.Set("1.0.0")
	mainFactory := factory.New(cli.NewSettings(
		cli.SetMode(version.Gitlab),
	))
	configTypes := config.PipelineConfigTypes
	tests := []struct {
		name    string
		args    args
		want    *types2.Jobs
		wantErr bool
	}{
		{
			name: "Can create php pipeline",
			args: args{
				factory:     mainFactory,
				configTypes: configTypes,
				pipelineConfig: types2.PipelineConfig{
					Pipelines: []types2.PipelineTypeConfig{
						{
							Type: types2.PipelineType("php"),
							Config: config.Php{
								Composer: config.PhpComposer{
									Path: ".",
								},
								ImageName: "testImage",
							},
							PipelineId: "php",
						},
					},
				},
			},
			want: &types2.Jobs{
				func() *types2.Job {
					job, _ := php.NewPhpUnit("testImage", ".", "", ".", "", "./vendor/bin/phpunit", map[string]int32{})
					return job
				}(),
			},
			wantErr: false,
		}, {
			name: "Can create docker pipeline",
			args: args{
				factory:     mainFactory,
				configTypes: configTypes,
				pipelineConfig: types2.PipelineConfig{
					Pipelines: []types2.PipelineTypeConfig{
						{
							Type: types2.PipelineType("docker"),
							Config: config.Docker{
								Path:      ".",
								ImageName: "testImage",
							},
							PipelineId: "golang",
						},
					},
				},
			},
			want: &types2.Jobs{
				func() *types2.Job {
					job, _ := docker.NewDaemonlessBuildkitBuild("testImage", ".", ".", "Dockerfile", []string{"${CI_REGISTRY_IMAGE}:1.0.0", "${CI_REGISTRY_IMAGE}:latest"}, []string{})
					return job
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := types2.LoadJobsFromPipelineConfig(tt.args.factory, tt.args.pipelineConfig, tt.args.configTypes)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadJobsFromPipelineConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.GetJobs()) != len(tt.want.GetJobs()) {
				t.Errorf("LoadFromPipelineConfig() got %d jobs, wanted %d", len(got.GetJobs()), len(tt.want.GetJobs()))
			}

			for i, job := range got.GetJobs() {
				if strings.Contains(tt.want.GetJobs()[i].GetName(), job.GetName()) {
					t.Errorf("LoadFromPipelineConfig() got job name %s, wanted %s", job.GetName(), tt.want.GetJobs()[i].GetName())
				}

				if !reflect.DeepEqual(job.CiJob.Script, tt.want.GetJobs()[i].CiJob.Script) {
					t.Errorf("LoadFromPipelineConfig() got job script %v, wanted %v", job.CiJob.Script.GetRenderedValue(), tt.want.GetJobs()[i].CiJob.Script.GetRenderedValue())
				}
			}
		})
	}
}
