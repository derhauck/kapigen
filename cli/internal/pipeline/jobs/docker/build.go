package docker

import (
	"fmt"
	"kapigen.kateops.com/docker"
	"kapigen.kateops.com/internal/gitlab"
	"kapigen.kateops.com/internal/gitlab/services"
	"kapigen.kateops.com/internal/pipeline/types"
)

func NewBuildkitBuild(path string, context string, dockerfile string, destination string) *types.Job {
	return types.NewJob("Build", docker.Buildkit, func(job *gitlab.CiJob) {
		job.Image.Entrypoint.
			Add("sh").
			Add("-c")
		job.AddVariable("BUILDKIT_HOST", "tcp://buildkitd:1234")
		daemon := services.New(docker.BUILDKITD, "buildkitd", 1234)
		daemon.Command().
			Add("--addr").
			Add("unix:///run/user/1000/buildkit/buildkitd.sock").
			Add("--addr").
			Add("tcp://0.0.0.0:1234").
			Add("--oci-worker-no-process-sandbox")
		job.Services.Add(daemon)
		command := fmt.Sprintf(`buildctl build --frontend dockerfile.v0 --local context="%s" --local dockerfile="%s" `, context, path)
		command = command + fmt.Sprintf(`--progress plain --opt filename="%s" --export-cache type=inline `, dockerfile)
		command = command + fmt.Sprintf(`--import-cache type=registry,ref="%s" `, destination)
		command = command + fmt.Sprintf(`--output type=image,name="%s",push=true `, destination)

		job.Script.Value.Add(command)
	})
}
