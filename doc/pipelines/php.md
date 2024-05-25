# PHP Pipeline
This pipeline allows you to run php tests

### Parameters:
```yaml
type: php
config:
  composerPath: string
  composerArgs: string
  phpunitXmlPath: string
  phpunitArgs: string
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
* `composerPath: [optional | default: '.']` The path to the composer.json file for installing the dependencies.
* `composerArgs: [optional]` Additional arguments for the composer install command.
* `phpunitXmlPath: [optional | default: '.']` The path to the `phpunit.xml`, relative to the composerPath.
* `phpunitArgs: [optional]` Additional arguments for the phpunit command.
* `imageName: [optional]` The image name to use for running the tests.


**docker**
* `docker: [optional]` Can be used to run the tests in a custom image
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
* `pressure:exclusive` (docker)

### Example:
**Image only**
```yaml
id: php
config:
  composerPath: code
  composerArgs: --no-dev
  phpunitXmlPath: tests
  phpunitArgs: --testsuite unit
  imageName: '${CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX}/php:8.1-cli-alpine3.15'
```
**Docker only**
```yaml
type: php
id: php
config:
  composerPath: code
  composerArgs: --no-dev
  phpunitXmlPath: tests
  phpunitArgs: --testsuite unit
  docker:
    path: cli
    context: .
    dockerfile: Dockerfile
    buildArgs:
      FOO: bar
```
