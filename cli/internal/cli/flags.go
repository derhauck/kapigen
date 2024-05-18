package cli

import (
	"os"

	"github.com/spf13/cobra"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/logger/level"
)

type PersistentConfig struct {
	Verbose bool
}

func PreparePersistentFlags(cmd *cobra.Command) *PersistentConfig {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		logger.ErrorE(err)
	}
	if verbose {
		err := os.Setenv("LOGGER_LEVEL", level.Debug.String())
		if err != nil {
			logger.ErrorE(err)
		}
	}

	return &PersistentConfig{
		verbose,
	}
}
