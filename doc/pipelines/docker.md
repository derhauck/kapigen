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
```

### Description:

* `path: [optional | default: '.']` the path to the Dockerfile.
* `context: [optional | default: <path>]` the context for the docker build
* `name: [optional | default: '']` the name of the image inside the project registry (will use the root as default)
* `dockerfile: [optional | default: 'Dockerfile']` the name of the Dockerfile inside the `<path>`
* `release: [optional | default: true]` whether the build will run on a release pipeline (tag pipeline) or only on feature branches

### Example:
```yaml
id: example
type: docker
config:
  path: cli
  context: cli
  dockerfile: Dockerfile
  name: cli
  release: true
```
