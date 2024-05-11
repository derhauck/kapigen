package job

import (
	"testing"

	"kapigen.kateops.com/internal/gitlab/tags"
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

func TestTags_Add(t *testing.T) {
	tests := []struct {
		name string
		want []tags.Size
	}{
		{
			name: "can add one",
			want: []tags.Size{
				tags.PRESSURE_MEDIUM,
			},
		},
		{
			name: "can add two",
			want: []tags.Size{
				tags.PRESSURE_MEDIUM,
				tags.PRESSURE_BUILDKIT,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actualTags = &Tags{}
			for _, tag := range tt.want {
				actualTags.Add(tag)
			}
			var hasTag bool
			for _, tag := range tt.want {
				for _, actualTag := range actualTags.Get() {
					if actualTag.Get() == tag.String() {
						hasTag = true
					}
				}

				if hasTag != true {
					t.Errorf("Add() = %v, want %v", actualTags, tag)
				}
				hasTag = false
			}

		})
	}
}
