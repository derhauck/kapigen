package cmd

import (
	"github.com/spf13/cobra"
	cmd "kapigen.kateops.com/cmd/version"
	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/cli"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/version"
)

var versionCmd = &cobra.Command{
	Use:              "version",
	Short:            "Will modify or create a version",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli.NewPersistentConfig(cmd)
		mode, err := cmd.Flags().GetString("mode")
		if err != nil {
			return err
		}
		logger.Info("will create settings")
		settings := cli.NewSettings(
			cli.SetMode(version.GetModeFromString(mode)),
		)
		factory := factory.New(settings)

		controller := factory.GetVersionController()
		logger.Info(controller.GetNewTag(""))
		return nil
	},
}

func init() {
	versionCmd.Flags().String("mode", version.Gitlab.Name(), "mode used for versioning: los,gitlab")
	versionCmd.AddCommand(cmd.NewCmd)
}
