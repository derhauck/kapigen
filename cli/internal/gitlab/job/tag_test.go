package job

import (
	"kapigen.kateops.com/internal/gitlab/tags"
	"testing"
)

func TestNewTag(t *testing.T) {
	t.Run("Can create new tag", func(t *testing.T) {
		medium := tags.PRESSURE_MEDIUM
		tag := NewTag(&medium)
		if tag == nil {
			t.Error("should not be nil")
		}
		if tag.Value != &medium {
			t.Error("should be equal")
		}
	})
	t.Run("Can create new tag with nil", func(t *testing.T) {
		tag := NewTag(nil)
		if tag == nil {
			t.Error("should not be nil")
		}
		if tag.Value == nil {
			t.Error("should not be nil")
		}

		if *tag.Value != tags.PRESSURE_MEDIUM {
			t.Error("should be the default")
		}
	})
}
