# Kapigen
#### Kateops Pipeline Generator

----
# CLI
## Commands
### Global Parameters
```shell
-v verbose log output
```
### Generate
`generate pipeline`
```shell
kapigen generate pipeline
```
parameters
```shell
--file 'output file' | default:  pipeline.yaml
--config 'pipeline config' | default: config.kapigen.yaml
```

## Pipeline Generation
### Versions
For the pipeline to release any version it will look up the current latest version for the repository and path 
in the LOS (Logic Operator Server).

It will also look up the version increase on the Merge Request and use it to increase the version in the LOS.
The version increase will be set once the Merge Request is merged. Until then you can only see the would be new version.

```
Merge Request   ->  current version + feature branch name
Main            ->  will release on success with new version and set LOS version to latest
```
