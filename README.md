# Kapigen
#### Kateops Pipeline Generator

Kapigen CLI is a tool to generate Gitlab pipelines based on configuration files.
It is an abstraction of the `.gitlab-ci.yml` syntax Gitlab uses. The goal is to enable developers to
leverage GitLab's CI/CD capabilities without needing extensive expertise in pipeline configuration.

---- 

## [Pipeline Generation Documentation](doc/index.md)

### Quick start

Kapigen allows you to generate pipelines based on the configuration defined in your repository.
To use inside your own project simply use the following inside your `.gitlab-ci.yml` file.

```yaml
include:
  - project: 'kateops/kapigen'
    ref: main
    file: 'default.gitlab-ci.yml'
```
In case you are not using the default kapigen tags you can overwrite them to your liking e.g. for the pipeline generation itself.
```yaml
default:
  tags:
    - saas-linux-medium-amd64
```

### Pipelines Configuration

[Kapigen Config](cli/config.kapigen.yaml):

```yaml
noop: true
versioning: true
dependencyProxy: ${CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX}
privateTokenName: CI_PIPELINE_TOKEN
tags:
  - saas-linux-medium-amd64
pipelines:
  - type: golang
    id: cli-golang
    config:
      path: cli
      coverage:
        packages:
          - kapigen.kateops.com/internal/...
      imageName: '${CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX}/golang:1.21.3-alpine3.18'
  - type: docker
    id: cli-docker
    config:
      path: cli
      name: cli
    needs:
      - cli-golang


```

For more detailed information visit the [docs](doc/index.md) inside the repository.
