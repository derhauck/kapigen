# Kapigen
#### Kateops Pipeline Generator

----
### Core
* [versioning](versioning.md)

---
## Pipelines
### Get Started
If you want to use the pipeline generation inside your own project add the following to your `.gitlab-ci.yml`
file.

**Public Gitlab**
```yaml
include:
  - project: 'kateops/kapigen'
    ref: main
    file: 'default.gitlab-ci.yml'
```

**Self Hosted**
```yaml
include:
  - remote: 'https://gitlab.com/kateops/kapigen/-/raw/main/default.gitlab-ci.yml'
```

In case you are not using the default kapigen tags you can overwrite them to your liking e.g. for the pipeline generation itself.
```yaml
default:
  tags:
    - saas-linux-medium-amd64
```
Inside your repository add the `.kapigen` folder and the `config.kapigen.yaml` file.
This is where you will define your pipelines.

---
### Pipeline Configuration
```yaml
versioning: bool
privateTokenName: string
noop: bool
tags: Array<string>
pipelines: Array<Pipeline>
```
* `versioning: [optional | default: false]` Will automatically create a tag with the version increase from the merge request. See [versioning](versioning.md)
  for more information.
* `privateTokenName [optional]` Will use the ENV variable for authenticating with Gitlab API instead of CI_JOB_TOKEN
* `noop: [optional | default: true]` Add a dummy job to every pipeline, so even if there is nothing to do the pipeline will be successful. Sometimes you just want to update a path of the repo without actually triggering a pipeline but the commit does anyway. So in order not to fail since the pipelines run only when changes where detected you need to have a dummy job.
* `tags: [optinal]` Allows to overwrite the tags for the pipeline so the jobs will start on a runner of your choosing.


**Pipeline**

All pipelines follow the same format:
```yaml
id: string
type: string
config: object
needs: Array<id>
tags: Array<string>
```
#### Description
* `id: [required]` Unique identifier inside this pipeline.
* `type: [required]` The type of pipeline that will run.
* `config: [required]` The configuration for the specific type, differs and each type will have its own configuration.
* `needs: [optional]` References pipeline ids from pipelines declared above the current pipeline. The current pipeline will then wait until said pipelines are finished before starting.
* `tags: [optinal]` Allows to overwrite the tags for the pipeline so the jobs will start on a runner of your choosing.
### Available Pipeline Configurations
  * [docker](pipelines/docker.md)
  * [golang](pipelines/golang.md)