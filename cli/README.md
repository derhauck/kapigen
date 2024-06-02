# Kapigen CLI Development

## Setup

```shell
go mod download
```

The [docker-compose](docker-compose.yml) file can be used to run tests via IDE with docker container. The [makefile](Makefile) 
uses the same configuration but set via `docker run` instead of `docker-compose`.


## Build

```shell
make build
```

Will create a `kapigen` bin in this directory.


## Tests
**Unit Tests**
```shell
make test
```
Will run all unit tests and show total coverage at the bottom.

**Pipeline Generation**
```shell
make pipeline
```
Will generate a pipeline for the current kapigen pipeline [configuration](config.kapigen.yaml) to `pipeline.yaml`.

**Test Pipeline Generation**
```shell
make test-pipeline
```
Will generate a pipeline for the current kapigen test pipeline [configuration](test.kapigen.yaml) to `pipeline.yaml`.
