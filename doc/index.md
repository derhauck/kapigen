# Kapigen
#### Kateops Pipeline Generator

----



### Core
* [versioning](versioning.md)
* [cli](cli.md)
### Available Pipeline Configurations
* [docker](pipelines/docker.md)
* [golang](pipelines/golang.md)
* [php](pipelines/php.md)
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

In case you are not using the default Kapigen tags you can overwrite them to your liking e.g. for the pipeline generation itself.
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
---
## Rules
Different pipeline rules sets are available. More complicated rules will be mentioned inside the pipeline configuration itself.
But the most basic ones are the following:

### Merge Request
Will execute when a merge request is opened, regardless of source or target branch. 
This one is a little special as it will allow only run pipelines with changes to the respective pipeline path. 
Which path is used for the change detection depends on the pipeline type and will be documented for each pipeline.
Some pipelines may even run regardless of changes as they might not be path specific.

**Example:**
```yaml
pipelines:
  - type: golang
    id: example-1
    config:
      path: first
      imageName: '${CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX}/golang:1.21.3-alpine3.18'
  - type: golang
    id: example-2
    config:
      path: second
      imageName: '${CI_DEPENDENCY_PROXY_GROUP_IMAGE_PREFIX}/golang:1.21.3-alpine3.18'
```
Now we make changes inside the `first` dir for our `example-1` go application. Then create a new merge request and a pipeline will run.
But instead of both pipelines it will only run the `example-1` pipeline as we only made changes to the `first` dir.

### Main
Will execute when any commit is made to the default branch.
### Release
Will execute when a new Tag is added to the repository.

---

