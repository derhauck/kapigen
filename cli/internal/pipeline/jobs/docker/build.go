package docker

import (
	"fmt"
	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/tags"
	"kapigen.kateops.com/internal/pipeline/types"
)

func NewBuildkitBuild(path string, context string, dockerfile string, destination []string) *types.Job {
	return types.NewJob("Build", docker.BUILDKIT.String(), func(ciJob *job.CiJob) {
		ciJob.Image.Entrypoint.
			Add("sh").
			Add("-c")
		ciJob.AddVariable("BUILDKIT_HOST", "tcp://buildkitd:1234")

		daemon := job.NewService(docker.BUILDKIT_DAEMON, "buildkitd", 1234)
		daemon.Command().
			Add("--addr").
			Add("unix:///run/user/1000/buildkit/buildkitd.sock").
			Add("--addr").
			Add("tcp://0.0.0.0:1234").
			Add("--oci-worker-no-process-sandbox")
		ciJob.Services.Add(daemon)

		timeout := job.NewService(docker.Alpine_3_18, "failover", 5000)
		timeout.Entrypoint().
			Add("sh").
			Add("-c")
		timeout.Command().
			Add("sleep 300; touch $CI_PROJECT_DIR/.status.auth")
		ciJob.Services.Add(timeout)

		auth := job.NewService(docker.CRANE_DEBUG, "crane", 5000)
		auth.Entrypoint().
			Add("sh").
			Add("-c")
		auth.Command().
			Add("while [ ! -f $CI_PROJECT_DIR/.status.init ]; do echo 'wait for init'; sleep 1; done; " +
				"export $(cat $CI_PROJECT_DIR/.env); " +
				"crane auth login -u ${REGISTRY_PUSH_USER} -p ${REGISTRY_PUSH_TOKEN} ${CI_REGISTRY}; " +
				"crane auth login -u ${REGISTRY_PUSH_USER} -p ${REGISTRY_PUSH_TOKEN} gitlab.kateops.com; " +
				"touch $CI_PROJECT_DIR/.status.auth")
		auth.AddVariable("DOCKER_CONFIG", "$CI_PROJECT_DIR")
		ciJob.Services.Add(auth)

		cmd := fmt.Sprintf(`buildctl build --frontend dockerfile.v0 --local context="%s" --local dockerfile="%s" `, context, path)
		parameters := fmt.Sprintf(`--progress plain --opt filename="%s" --export-cache type=inline `, dockerfile)
		cache := fmt.Sprintf(`--import-cache type=registry,ref="%s" `, destination)
		push := ""
		for _, d := range destination {
			push = fmt.Sprintf(`%s--output type=image,name="%s",push=true `, push, d)
		}
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
		ciJob.BeforeScript.Value.
			Add(`echo "REGISTRY_PUSH_USER=$REGISTRY_PUSH_USER" > .env`).
			Add(`echo "REGISTRY_PUSH_TOKEN=$REGISTRY_PUSH_TOKEN" >> .env`).
			Add("touch .status.init").
			Add("while [ ! -f $CI_PROJECT_DIR/.status.auth ]; do echo 'wait for auth'; sleep 1; done")
		ciJob.Script.Value.
			Add(command)
		ciJob.Rules = *job.DefaultPipelineRules()
		ciJob.AddVariable("KTC_PATH", path)
		ciJob.AddVariable("DOCKER_CONFIG", "$CI_PROJECT_DIR")
		ciJob.Tags.Add(tags.PRESSURE_BUILDKIT)
	})
}
