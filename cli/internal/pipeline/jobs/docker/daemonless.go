package docker

import (
	"fmt"
	"strings"

	"gitlab.com/kateops/kapigen/cli/internal/docker"
	"gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/enum"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

func NewDaemonlessBuildkitBuild(imageName string, path string, context string, dockerfile string, destination []string, buildArgs []string) (*types.Job, error) {
	if imageName == "" {
		imageName = docker.BUILDKIT_ROTLESS.String()
	}
	if len(destination) == 0 {
		return nil, wrapper.DetailedErrorf("destination must be set")
	}
	return types.NewJob("Daemonless Build", imageName, func(ciJob *job.CiJob) {
		ciJob.Image.Entrypoint.
			Push("sh").
			Push("-c")

		timeout := job.NewService(docker.Alpine_3_18.String(), "failover", 5000)
		timeout.Entrypoint().
			Push("sh").
			Push("-c")
		timeout.Command().
			Push("sleep 30; touch ${CI_PROJECT_DIR}/.status.auth")
		ciJob.Services.Add(timeout)

		auth := job.NewService(docker.CRANE_DEBUG.String(), "crane", 5000)
		auth.Entrypoint().
			Push("sh").
			Push("-c")
		auth.Command().
			Push("while [ ! -f ${CI_PROJECT_DIR}/.status.init ]; do echo 'wait for init'; sleep 1; done; " +
				"export $(cat ${CI_PROJECT_DIR}/.env); " +
				"crane auth login -u \"${CI_REGISTRY_USER}\" -p \"${CI_JOB_TOKEN}\" \"${CI_REGISTRY}\"; " +
				"crane auth login -u \"${CI_DEPENDENCY_PROXY_USER}\" -p \"${CI_DEPENDENCY_PROXY_PASSWORD}\" \"${CI_DEPENDENCY_PROXY_SERVER}\"; " +
				"crane auth login -u \"${CI_DEPENDENCY_PROXY_USER}\" -p \"${CI_DEPENDENCY_PROXY_PASSWORD}\" \"${CI_SERVER_HOST}\"; " +
				"touch ${CI_PROJECT_DIR}/.status.auth; " +
				"chmod 666 ${CI_PROJECT_DIR}/config.json")
		auth.AddVariable("DOCKER_CONFIG", "${CI_PROJECT_DIR}")
		ciJob.Services.Add(auth)
		args := ""
		for _, buildArg := range buildArgs {
			args += fmt.Sprintf("--opt build-arg:%s", buildArg)
		}

		cmd := fmt.Sprintf(`buildctl-daemonless.sh build --frontend dockerfile.v0 --local context="%s" --local dockerfile="%s" %s`, context, path, args)
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
			Push(`echo "CI_JOB_TOKEN=$CI_JOB_TOKEN" > .env`).
			Push(`echo "CI_DEPENDENCY_PROXY_PASSWORD=$CI_DEPENDENCY_PROXY_PASSWORD" >> .env`).
			Push("touch .status.init").
			Push("while [ ! -f ${CI_PROJECT_DIR}/.status.auth ]; do echo 'wait for auth'; sleep 1; done")
		ciJob.Script.Value.
			Push(command)
		ciJob.AddVariable("KTC_PATH", path).
			AddVariable("BUILDKITD_FLAGS", "--oci-worker-no-process-sandbox").
			AddVariable("DOCKER_CONFIG", "${CI_PROJECT_DIR}").
			AddVariable("BUILDCTL_CONNECT_RETRIES_MAX", "52")
		ciJob.Tags.Add(enum.TagPressureExclusive)
	}), nil
}
