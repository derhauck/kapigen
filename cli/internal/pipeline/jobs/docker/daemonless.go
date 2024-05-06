package docker

import (
	"fmt"
	"strings"

	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/tags"
	"kapigen.kateops.com/internal/pipeline/types"
)

func NewDaemonlessBuildkitBuild(path string, context string, dockerfile string, destination []string) *types.Job {
	return types.NewJob("Daemonless Build", docker.BUILDKIT_ROTLESS.String(), func(ciJob *job.CiJob) {
		ciJob.Image.Entrypoint.
			Add("sh").
			Add("-c")

		timeout := job.NewService(docker.Alpine_3_18, "failover", 5000)
		timeout.Entrypoint().
			Add("sh").
			Add("-c")
		timeout.Command().
			Add("sleep 30; touch ${CI_BUILDS_DIR}/.status.auth")
		ciJob.Services.Add(timeout)

		auth := job.NewService(docker.CRANE_DEBUG, "crane", 5000)
		auth.Entrypoint().
			Add("sh").
			Add("-c")
		auth.Command().
			Add("while [ ! -f ${CI_BUILDS_DIR}/.status.init ]; do echo 'wait for init'; sleep 1; done; " +
				"export $(cat $CI_PROJECT_DIR/.env); " +
				"crane auth login -u ${REGISTRY_PUSH_USER} -p ${REGISTRY_PUSH_TOKEN} ${GITLAB_REGISTRY}; " +
				"crane auth login -u ${CI_DEPENDENCY_PROXY_USER} -p ${CI_DEPENDENCY_PROXY_PASSWORD} ${CI_DEPENDENCY_PROXY_SERVER}; " +
				"touch ${CI_BUILDS_DIR}/.status.auth")
		auth.AddVariable("DOCKER_CONFIG", "${CI_BUILDS_DIR}")
		ciJob.Services.Add(auth)

		cmd := fmt.Sprintf(`buildctl-daemonless.sh build --frontend dockerfile.v0 --local context="%s" --local dockerfile="%s" `, context, path)
		parameters := fmt.Sprintf(`--progress plain --opt filename="%s" --export-cache type=inline `, dockerfile)
		cache := fmt.Sprintf(`--import-cache type=registry,ref="%s" `, destination[0])
		push := fmt.Sprintf(`--output type=image,\"name=%s\",push=true `, strings.Join(destination, ","))
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
			Add(`echo "GITLAB_REGISTRY=$CI_REGISTRY" >> .env`).
			Add(`echo "CI_DEPENDENCY_PROXY_USER=$CI_DEPENDENCY_PROXY_USER" >> .env`).
			Add(`echo "CI_DEPENDENCY_PROXY_PASSWORD=$CI_DEPENDENCY_PROXY_PASSWORD" >> .env`).
			Add(`echo "CI_DEPENDENCY_PROXY_SERVER=$CI_DEPENDENCY_PROXY_SERVER" >> .env`).
			Add("touch .status.init").
			Add("while [ ! -f ${CI_BUILDS_DIR}/.status.auth ]; do echo 'wait for auth'; sleep 1; done")
		ciJob.Script.Value.
			Add(command)
		ciJob.Rules = *job.DefaultPipelineRules()
		ciJob.AddVariable("KTC_PATH", path).
			AddVariable("BUILDKITD_FLAGS", "--oci-worker-no-process-sandbox").
			AddVariable("DOCKER_CONFIG", "${CI_BUILDS_DIR}").
			AddVariable("BUILDCTL_CONNECT_RETRIES_MAX", "52")
		ciJob.Tags.Add(tags.PRESSURE_BUILDKIT)
	})
}
