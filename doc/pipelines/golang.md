# Golang Pipeline
This pipeline allows you to run golang tests

### Parameters:
```yaml
type: golang
config:
  path: string
  coverage:
    packages: Array<string>
  imageName: string
  docker:
    path: string
    dockerfile: string
    context: string
```

### Description:
The pipeline will run golang tests. In order to execute those it will run the jobs inside the specified image. 
Either by using `imageName` or the `docker` configuration. If both are set the `docker` config takes precedence.
* `path: [optional | default: '.']` The path to the golang code for executing the tests.
* `imageName: [optional]` The image name to use for running the tests

**coverage**
* `coverage: [optional]` The coverage options
* `coverage.packages: [optional | default: []]` The packages to consider for the coverage calculations

**docker**
* `docker: [optional]` Can be used to run the tests in a custom image
* `docker.path: [optional | default: '.']` The path to the Dockerfile.
* `docker.dockerfile: [optional | default: 'Dockerfile']` The name of the Dockerfile inside the `<path>`.
* `docker.context: [optional | default: <path>]` The context for the docker build.

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
  docker:
    path: cli
    context: cli
    dockerfile: Dockerfile
```
