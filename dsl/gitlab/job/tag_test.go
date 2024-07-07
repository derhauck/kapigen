package job

import (
	"testing"

	"gitlab.com/kateops/kapigen/dsl/enum"
)

func TestNewTag(t *testing.T) {
	t.Run("Can create new tag", func(t *testing.T) {
		medium := enum.TagPressureMedium
		tag := NewTag(&medium)
		if tag == nil {
			t.Error("should not be nil")
		}
		if tag == nil && tag.Value != nil && tag.Value != &medium {
			t.Error("should be equal")
		}
	})
	t.Run("Can create new tag with nil", func(t *testing.T) {
		tag := NewTag(nil)
		if tag == nil {
			t.Error("should not be nil")
		}
		if tag != nil && tag.Value == nil {
			t.Error("should not be nil")
		}

		if tag != nil && tag.Value != nil && *tag.Value != enum.TagPressureMedium {
			t.Error("should be the default")
		}
	})
}

func TestTags_Add(t *testing.T) {
	tests := []struct {
		name string
		want []enum.Tag
	}{
		{
			name: "can add one",
			want: []enum.Tag{
				enum.TagPressureMedium,
			},
		},
		{
			name: "can add two",
			want: []enum.Tag{
				enum.TagPressureMedium,
				enum.TagPressureExclusive,
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
