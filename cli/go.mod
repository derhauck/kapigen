module gitlab.com/kateops/kapigen/cli

go 1.21

require (
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/kylelemons/godebug v1.1.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/spf13/cobra v1.8.1
	github.com/xanzy/go-gitlab v0.106.0
	gitlab.com/kateops/kapigen/dsl v0.0.0-20240707081511-1478eb656d50
	gopkg.in/yaml.v3 v3.0.1
)

replace gitlab.com/kateops/kapigen/dsl => ../dsl

require (
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/oauth2 v0.21.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)
