# Kapigen
#### Kateops Pipeline Generator

----
### Core
* [versioning](versioning.md)

### Pipelines
All pipelines follow the same format:
```yaml
id: string
type: string
config: object
need: Array<id>
```
#### Description
* `id: [required]` Unique identifier inside this pipeline.
* `type: [required]` The type of pipeline that will run.
* `config: [required]` The configuration for the specific type, differs and each type will have its own configuration.
* `need: [optional | default: []]` References pipeline ids from pipelines declared above the current pipeline. The current pipeline will then wait until said pipelines are finished before starting.

### Available Pipelines
  * [docker](pipelines/docker.md)