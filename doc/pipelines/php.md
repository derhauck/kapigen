# PHP Pipeline
This pipeline allows you to run php tests

### Parameters:
```yaml
type: php
config:
  composer:
    path: string
    args: string
  phpunit:
    path: string
    args: string
  services:
    - name: string
      port: number
      imageName: string
      docker:
        path: string
        dockerfile: string
        context: string
        buildArgs: Record<string, string>
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

**composer** (optional)
* `composer.path: [optional | default: '.']` The path to the composer.json file for installing the dependencies.
* `composer.args: [optional | default: '--no-progress --no-cache --no-interaction']` Additional arguments for the composer install command.

**phpunit** (optional)
* `phpunit.path: [optional | default: '.']` The path to the `phpunit.xml`.
* `phpunit.args: [optional]` Additional arguments for the phpunit command.

**services** (optional)
* `services: [optional]` Array of services which will run as sidecar containers during the job.

**service** (optional)
* `service.name: [required]` The dns name under which the service is available.
* `service.port: [required]` The TCP port under which it listens when it is ready to accept connections.
* `service.imageName: [optional]` The image name to use as a sidecar container.
* `service.docker [optional]` Can be used to run the sidecar container in a custom image. 

**docker** (optional)
* `docker: [optional]` Can be used to run the tests in a custom image.
* `docker.path: [required]` The path to the Dockerfile.
* `docker.dockerfile: [optional | default: 'Dockerfile']` The name of the Dockerfile inside the `<path>`.
* `docker.context: [optional | default: <path>]` The context for the docker build.
* `docker.buildArgs: [optional]` Additional build arguments for the docker build.

### Rules:
Pipeline will execute for the following types:
* `Merge Request` Uses `composerPath` and `docker.context` to watch for changes.
* `Main`

### Tags:
* `pressure:medium`
* `pressure:exclusive` (docker,service.docker)

### Example:
**Image only**
```yaml
id: php
config:
  composer:
    path: code
    args: --no-dev
  phpunit:
    path: tests
    args: --testsuite unit
  services:
    - name: mysql
      port: 3306
      imageName: mysql:8
  imageName: '${CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX}/php:8.1-cli-alpine3.15'
```
**Docker only**
```yaml
type: php
id: php
config:
  composer: 
    path: code
    args: --no-dev
  phpunit:
    path: tests
    args: --testsuite unit
  services:
    - name: mysql
      port: 3306
      docker:
        path: cli/database
        context: cli
        dockerfile: Dockerfile
        buildArgs:
          FOO: bar
  docker:
    path: cli
    context: .
    dockerfile: Dockerfile
    buildArgs:
      FOO: bar
```

For a real example look at the [php tests](../../cli/tests/php) inside this repository.

The pipeline configuration for running those tests: [Test Pipeline Config](../../cli/test.kapigen.yaml)
```yaml
noop: true
versioning: false
dependencyProxy: ${CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX}
#tags:
#  - saas-linux-medium-amd64
pipelines:
# ...
  - type: php
    id: php
    config:
      composer:
        path: cli/tests/php
      docker:
        path: cli/tests/php
```
