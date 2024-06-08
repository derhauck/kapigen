package pipeline

import (
	"reflect"
	"strings"
	"testing"

	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/cli"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/pipeline/config"
	"kapigen.kateops.com/internal/pipeline/jobs/docker"
	"kapigen.kateops.com/internal/pipeline/jobs/php"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/version"
)

func TestLoadJobsFromPipelineConfig(t *testing.T) {
	type args struct {
		factory        *factory.MainFactory
		pipelineConfig types.PipelineConfig
		configTypes    map[types.PipelineType]types.PipelineConfigInterface
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
		want    *types.Jobs
		wantErr bool
	}{
		{
			name: "Can create php pipeline",
			args: args{
				factory:     mainFactory,
				configTypes: configTypes,
				pipelineConfig: types.PipelineConfig{
					Pipelines: []types.PipelineTypeConfig{
						{
							Type: types.PipelineType("php"),
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
			want: &types.Jobs{
				func() *types.Job {
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
				pipelineConfig: types.PipelineConfig{
					Pipelines: []types.PipelineTypeConfig{
						{
							Type: types.PipelineType("docker"),
							Config: config.Docker{
								Path:      ".",
								ImageName: "testImage",
							},
							PipelineId: "golang",
						},
					},
				},
			},
			want: &types.Jobs{
				func() *types.Job {
					job := docker.NewDaemonlessBuildkitBuild("testImage", ".", ".", "Dockerfile", []string{"${CI_REGISTRY_IMAGE}:1.0.0", "${CI_REGISTRY_IMAGE}:latest"}, []string{})
					return job
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := types.LoadJobsFromPipelineConfig(tt.args.factory, tt.args.pipelineConfig, tt.args.configTypes)
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
