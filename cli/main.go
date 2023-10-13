package main

import (
	"kapigen.kateops.com/cmd"
	"kapigen.kateops.com/internal/environment"
	"os"
)

func main() {
	if os.Getenv("ENVIRONMENT") == "DEV" {
		environment.SetLocalEnv()
	}
	cmd.Execute()
}
