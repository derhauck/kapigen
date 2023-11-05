package version

import (
	"fmt"
	"github.com/spf13/cobra"
	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/cli"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/version"
)

var versionCmd = &cobra.Command{
	Use:              "new",
	Short:            "Will create a new version",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli.NewPersistentConfig(cmd)
		mode, err := cmd.Flags().GetString("mode")
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}
		logger.Info("will create settings")
		settings := cli.NewSettings(
			cli.SetMode(version.GetModeFromString(mode)),
		)
		factory := factory.New(settings)

		controller := factory.GetVersionController()
		logger.Info(fmt.Sprintf("Created new version %s", controller.SetNewVersion(path)))
		return nil
	},
}

func init() {
	versionCmd.Flags().String("mode", version.Gitlab.Name(), "mode used for versioning: los,gitlab")
	versionCmd.Flags().String("path", "", "path for version (los mode only)")
}
