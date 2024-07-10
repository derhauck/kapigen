module gitlab.com/kateops/kapigen/cli

go 1.21

require (
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/gkampitakis/go-snaps v0.5.4
	github.com/kylelemons/godebug v1.1.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/spf13/cobra v1.8.1
	github.com/xanzy/go-gitlab v0.106.0
	gitlab.com/kateops/kapigen/dsl v0.0.0-20240707081511-1478eb656d50
	gopkg.in/yaml.v3 v3.0.1
)

replace gitlab.com/kateops/kapigen/dsl => ../dsl

require (
	github.com/gkampitakis/ciinfo v0.3.0 // indirect
	github.com/gkampitakis/go-diff v1.3.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/maruel/natural v1.1.1 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tidwall/gjson v1.17.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	golang.org/x/oauth2 v0.21.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)
