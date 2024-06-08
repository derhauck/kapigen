
## CLI
Kapigen CLI is a tool to generate Gitlab pipelines based on configuration files. 
It is an abstraction of the `.gitlab-ci.yml` syntax Gitlab uses. The goal is to enable developers to
leverage GitLab's CI/CD capabilities without needing extensive expertise in pipeline configuration.

### Commands

#### Global Parameters

```shell
-v  # Verbose log output
--private-token  # ENV var for private token to use in Gitlab API calls
```
---
## Pipeline
### Generate
Will generate a pipeline file based on the configuration defined in the `config.kapigen.yaml` file.

```shell
kapigen pipeline generate
```

#### Parameters:
- `--file 'output file' | default: 'pipeline.yml'`: Specifies the output file.
- `--config 'pipeline config | default: 'config.kapigen.yaml'`: Specifies the pipeline configuration file. 
- `--no-merge`:  Disables automatic merging of duplicate jobs in the pipeline. Use this to avoid potential issues with automatic job merging.

### Reports
Downloads reports for the current pipeline, looking inside downstream pipelines for jobs with JUnit artifacts.

```shell
kapigen pipeline reports
```
#### Parameters:
- `--config 'pipeline config | default: 'config.kapigen.yaml'`: Specifies the pipeline configuration file. 

---

## Notes
* Ensure the configuration file is correctly set up before running the CLI commands. (can even be empty)
* Use the -v flag for detailed log output during execution.