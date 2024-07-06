package cli

import (
	"testing"

	"gitlab.com/kateops/kapigen/cli/internal/version"
)

func TestSetMode(t *testing.T) {
	t.Run("Can set mode in", func(t *testing.T) {
		setting := NewSettings(SetMode(version.Gitlab))
		if setting.Mode != version.Gitlab {
			t.Errorf("should be gitlab mode, is %s", setting.Mode.Name())
		}

		setting = NewSettings(SetMode(version.FILE))
		if setting.Mode != version.FILE {
			t.Errorf("should be file mode, is %s", setting.Mode.Name())
		}

		setting = NewSettings(SetMode(version.None))
		if setting.Mode != version.None {
			t.Errorf("should be file mode, is %s", setting.Mode.Name())
		}
	})
}
