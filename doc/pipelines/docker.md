# Docker Pipeline
This pipeline allows you to build a docker image via buildkit

### Parameters:
```yaml
type: docker
config:
  path: string
  context: string
  name: string
  dockerfile: string
  release: bool
  buildArgs: Record<string, string>
  imageName: string
```

### Description:

* `path: [optional | default: '.']` The path to the Dockerfile.
* `context: [optional | default: <path>]` The context for the docker build.
* `name: [optional | default: '']` The name of the image inside the project registry (will use the root as default).
* `dockerfile: [optional | default: 'Dockerfile']` The name of the Dockerfile inside the `<path>`.
* `release: [optional | default: true]` Whether the build will run on a release pipeline (tag pipeline) or only on feature branches.
* `buildArgs: [optional]` Additional build arguments for the docker build.
* `imageName: [optional | default: 'moby/buildkit:v0.12.3-rootless']` The image name to use for running the build job.

### Rules:
Pipeline will execute for the following types:
* `Merge Request` Uses `context` to watch for changes.
* `Release`

### Tags:
* `pressure:exclusive`

### Example
```yaml
id: example
type: docker
config:
  path: cli
  context: cli
  dockerfile: Dockerfile
  name: cli
  release: true
  buildArgs:
    FOO: bar 
```
