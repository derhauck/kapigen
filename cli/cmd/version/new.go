package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/kateops/kapigen/cli/factory"
	"gitlab.com/kateops/kapigen/cli/internal/cli"
	"gitlab.com/kateops/kapigen/cli/internal/version"
	"gitlab.com/kateops/kapigen/dsl/logger"
)

var NewCmd = &cobra.Command{
	Use:              "new",
	Short:            "Will create a new version",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli.PreparePersistentFlags(cmd)
		mode, err := cmd.Flags().GetString("mode")
		if err != nil {
			return err
		}
		privateTokenName, err := cmd.Flags().GetString("private-token")
		if err != nil {
			return err
		}
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}
		logger.Info("will create settings")
		settings := cli.NewSettings(
			cli.SetMode(version.GetModeFromString(mode)),
			cli.SetPrivateToken(privateTokenName),
		)
		factory := factory.New(settings)

		controller := factory.GetVersionController()
		newVersion := controller.SetNewVersion(path)
		if newVersion == version.EmptyTag {
			logger.Info("Empty version ref, did not create new tag")
			return nil
		}
		logger.Info(fmt.Sprintf("Created new version %s", newVersion))
		return nil
	},
}

func init() {
	NewCmd.Flags().String("mode", version.FILE.Name(), "mode used for versioning: file,gitlab")
	NewCmd.Flags().String("path", "", "path for version (file mode only)")
}
