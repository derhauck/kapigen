# Kapigen
#### Kateops Pipeline Generator

----

## CLI

### Commands

#### Global Parameters

```shell
-v  # Verbose log output
```

#### Generate

`generate pipeline`

```shell
kapigen generate pipeline
```

### Parameters:
- `--file 'output file'`: Specifies the output file. Default is `pipeline.yaml`.
- `--config 'pipeline config'`: Specifies the pipeline configuration file. Default is `config.kapigen.yaml`.

## Pipeline Generation

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

[Kapigen Config](cli/config.kapigen.yaml)

For more detailed information visit the [docs](doc/index.md) inside the repository.