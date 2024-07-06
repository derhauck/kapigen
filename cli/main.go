package main

import (
	"os"

	"gitlab.com/kateops/kapigen/cli/cmd"
	"gitlab.com/kateops/kapigen/dsl/environment"
)

func main() {
	if os.Getenv("ENVIRONMENT") == "DEV" {
		environment.SetLocalEnv()
	}
	cmd.Execute()
}
