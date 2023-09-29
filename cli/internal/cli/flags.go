package cli

import (
	"github.com/spf13/cobra"
	"kapigen.kateops.com/internal/logger"
	"os"
)

type CmdConfig[T any] struct {
	Persistent *PersistentConfig
	Local      *T
}
type PersistentConfig struct {
	Verbose bool
}

func NewConfig[T any](cmd *cobra.Command, config T) *CmdConfig[T] {
	return &CmdConfig[T]{
		Persistent: NewPersistentConfig(cmd),
		Local:      &config,
	}
}

func NewPersistentConfig(cmd *cobra.Command) *PersistentConfig {
	return preparePersistentFlags(cmd)
}
func preparePersistentFlags(cmd *cobra.Command) *PersistentConfig {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		logger.ErrorE(err)
	}
	if verbose {
		err := os.Setenv("LOGGER_LEVEL", logger.Level.Debug)
		if err != nil {
			logger.ErrorE(err)
		}
	}

	return &PersistentConfig{
		verbose,
	}
}
