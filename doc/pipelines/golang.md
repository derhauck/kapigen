# Golang Pipeline
This pipeline allows you to run golang tests

### Parameters:
```yaml
type: golang
config:
  path: string
  imageName: string
  coverage:
    packages: Array<string>
    source: string
  lint:
    imageName: string
    mode: string
  services:
    - name: string
      port: number
      imageName: string
      docker:
        path: string
        dockerfile: string
        context: string
        buildArgs: Record<string, string>
  docker:
    path: string
    dockerfile: string
    context: string
    buildArgs: Record<string, string>
```

### Description:
The pipeline will run golang tests. In order to execute those it will run the jobs inside the specified image. 
Either by using `imageName` or the `docker` configuration. If both are set the `docker` config takes precedence.
* `path: [optional | default: '.']` The path to the golang code for executing the tests.
* `imageName: [optional]` The image name to use for running the tests.

**coverage**
* `coverage: [optional]` The coverage options.
* `coverage.packages: [optional | ]` The packages to consider for the coverage calculations.
* `coverage.source: [optional | default: './...' ]` The source code to consider for the coverage calculations.

**lint**
* `lint: [optional]` The linting options.
* `lint.mode: [optional | default: 'enabled']` Supports the following linter modes `enabled`, `permissive`, `disabled`.

**services** (optional)
* `services: [optional]` Array of services which will run as sidecar containers during the job.

**service** (optional)
* `service.name: [required]` The dns name under which the service is available.
* `service.port: [required]` The TCP port under which it listens when it is ready to accept connections.
* `service.imageName: [optional]` The image name to use as a sidecar container.
* `service.docker [optional]` Can be used to run the sidecar container in a custom image.

**docker**
* `docker: [optional]` Can be used to run the tests in a custom image
* `docker.path: [required]` The path to the Dockerfile.
* `docker.dockerfile: [optional | default: 'Dockerfile']` The name of the Dockerfile inside the `<path>`.
* `docker.context: [optional | default: <path>]` The context for the docker build.
* `docker.buildArgs: [optional]` Additional build arguments for the docker build.

### Rules:
Pipeline will execute for the following types:
* `Merge Request` Uses `path` and `docker.context` to watch for changes.
* `Main`

### Tags:
* `pressure:medium`
* `pressure:exclusive` (docker)

### Example:
**Image only**
```yaml
type: golang
id: example
config:
  path: cli
  coverage:
    packages:
      - kapigen.kateops.com/internal/...
  imageName: '${CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX}/golang:1.21.3-alpine3.18'
```
**Docker only**
```yaml
type: golang
id: example
config:
  path: cli
  coverage:
    packages:
      - kapigen.kateops.com/internal/...
    source: ./...
  lint:
    mode: permissive
    imageName: golangci/golangci-lint:v1.59.1
  services:
    - name: postgres
      port: 5432
      docker:
        path: cli/database
        context: cli
        dockerfile: Dockerfile
        buildArgs:
          FOO: bar
  docker:
    path: cli
    context: cli
    dockerfile: Dockerfile
    buildArgs:
      FOO: bar
```
