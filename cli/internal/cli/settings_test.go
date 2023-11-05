package cli

import (
	"kapigen.kateops.com/internal/version"
	"testing"
)

func TestSetMode(t *testing.T) {
	t.Run("Can set mode in", func(t *testing.T) {
		setting := NewSettings(SetMode(version.Gitlab))
		if setting.Mode != version.Gitlab {
			t.Errorf("should be gitlab mode, is %s", setting.Mode.Name())
		}

		setting = NewSettings(SetMode(version.Los))
		if setting.Mode != version.Los {
			t.Errorf("should be los mode, is %s", setting.Mode.Name())
		}
	})
}
