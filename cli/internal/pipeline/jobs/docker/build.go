package docker

import (
	"fmt"
	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/gitlab"
	"kapigen.kateops.com/internal/gitlab/rules"
	"kapigen.kateops.com/internal/gitlab/services"
	"kapigen.kateops.com/internal/gitlab/tags"
	"kapigen.kateops.com/internal/pipeline/types"
)

func NewBuildkitBuild(path string, context string, dockerfile string, destination string) *types.Job {
	return types.NewJob("Build", docker.BUILDKIT, func(job *gitlab.CiJob) {
		job.Image.Entrypoint.
			Add("sh").
			Add("-c")
		job.AddVariable("BUILDKIT_HOST", "tcp://buildkitd:1234")

		daemon := services.New(docker.BUILDKIT_DAEMON, "buildkitd", 1234)
		daemon.Command().
			Add("--addr").
			Add("unix:///run/user/1000/buildkit/buildkitd.sock").
			Add("--addr").
			Add("tcp://0.0.0.0:1234").
			Add("--oci-worker-no-process-sandbox")
		job.Services.Add(daemon)

		auth := services.New(docker.CRANE_DEBUG, "crane", 5000)
		auth.Entrypoint().
			Add("sh").
			Add("-c")
		auth.Command().
			Add("while [ ! -f $CI_PROJECT_DIR/.status.init ]; do echo 'wait for init'; sleep 1; done").
			Add("crane auth login -u ${REGISTRY_PUSH_USER} -p ${REGISTRY_PUSH_TOKEN} ${CI_REGISTRY}").
			Add("crane auth login -u ${REGISTRY_PUSH_USER} -p ${REGISTRY_PUSH_TOKEN} gitlab.kateops.com").
			Add("touch $CI_PROJECT_DIR/.status.auth")
		auth.AddVariable("DOCKER_CONFIG", "$CI_PROJECT_DIR")
		job.Services.Add(auth)

		cmd := fmt.Sprintf(`buildctl build --frontend dockerfile.v0 --local context="%s" --local dockerfile="%s" `, context, path)
		parameters := fmt.Sprintf(`--progress plain --opt filename="%s" --export-cache type=inline `, dockerfile)
		cache := fmt.Sprintf(`--import-cache type=registry,ref="%s" `, destination)
		push := fmt.Sprintf(`--output type=image,name="%s",push=true `, destination)
		command := fmt.Sprintf(
			"%s \\\n"+
				"%s \\\n"+
				"%s \\\n"+
				"%s \\\n",
			cmd,
			parameters,
			cache,
			push,
		)
		job.BeforeScript.Value.
			Add("touch .status.init").
			Add("while [ ! -f $CI_PROJECT_DIR/.status.auth ]; do echo 'wait for auth'; sleep 1; done")
		job.Script.Value.
			Add("sleep 300").
			Add(command)
		job.Rules = *rules.DefaultPipelineRules()
		job.Variables["KTC_PATH"] = path
		job.Variables["DOCKER_CONFIG"] = "$CI_PROJECT_DIR"
		job.Tags.Add(tags.PRESSURE_BUILDKIT)
	})
}
