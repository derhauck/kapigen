package main

import (
	"os"

	"kapigen.kateops.com/cmd"
	"kapigen.kateops.com/internal/environment"
)

func main() {
	if os.Getenv("ENVIRONMENT") == "DEV" {
		environment.SetLocalEnv()
	}
	cmd.Execute()
}
