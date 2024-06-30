# PHP Pipeline
This pipeline allows you to run php tests

### Parameters:
```yaml
type: generic
config:
  scripts: Array<string>
  variables: Record<string, string>
  changes: Array<string>
  artifacts: object
  imageName: string
  docker:
    path: string
    dockerfile: string
    context: string
    buildArgs: Record<string, string>
```

### Description:
The pipeline will run php tests. In order to execute those it will run the jobs inside the specified image. 
Either by using `imageName` or the `docker` configuration. If both are set the `docker` config takes precedence.
* `imageName: [optional]` The image name to use for running the tests.

* `variables: [optional]` Additional variables to set for the pipeline.
* `changes: [optional | default: '.']` The changes to watch for in the pipeline
* `scripts: [required]` The instructions to run this job.

* `artifacts: [optional]` The artifacts to upload to the pipeline. Same as native gitlab
* `artifacts.paths: [required]` The file paths to upload.
* `artifacts.name: [optional | default: 'artifacts']` The name of the artifact.

 
* **docker** (optional)
* `docker: [optional]` Can be used to run the tests in a custom image.
* `docker.path: [required]` The path to the Dockerfile.
* `docker.dockerfile: [optional | default: 'Dockerfile']` The name of the Dockerfile inside the `<path>`.
* `docker.context: [optional | default: <path>]` The context for the docker build.
* `docker.buildArgs: [optional]` Additional build arguments for the docker build.

### Rules:
Pipeline will execute for the following types:
* `Merge Request` Uses `changes` to watch for changes.
* `Main`

### Tags:
* `pressure:medium`
* `pressure:exclusive` (docker)

### Example:
**Image only**
```yaml
id: php
config:
  scripts:
     - 'echo " \"hello ${HOST}!" > test.txt'
  variables:
    HOST: "vm-private-ci-1"
  artifacts:
    name: test
    paths:
      - test.txt
  imageName: '${CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX}/php:8.1-cli-alpine3.15'
```
**Docker only**
```yaml
type: php
id: php
config:
  scripts:
    - 'echo " \"hello ${HOST}!" > test.txt'
  variables:
    HOST: "vm-private-ci-1"
  artifacts:
    name: test
    paths:
      - test.txt
  docker:
    path: cli
    context: .
    dockerfile: Dockerfile
    buildArgs:
      FOO: bar
```

For a real example look at the inside this repository.

The pipeline configuration for running those tests: [Test Pipeline Config](../../cli/test.kapigen.yaml)
```yaml
noop: true
versioning: false
dependencyProxy: ${CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX}
#tags:
#  - saas-linux-medium-amd64
pipelines:
# ...
  - type: generic
    id: generic-job
    config:
      scripts:
        - 'echo " \"hello ${HOST}!" > test.txt'
      variables:
        HOST: "vm-private-ci-1"
      artifacts:
        name: test
        paths:
          - test.txt
```
