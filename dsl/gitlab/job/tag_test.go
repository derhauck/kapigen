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

		if tag != nil && tag.Value != nil && tag.Value.String() != enum.TagPressureMedium.String() {
			t.Errorf("should be the default '%s'", enum.TagPressureMedium.String())
		}
	})
}

type TestTag string

func (t *TestTag) String() string {
	value := *t
	return string(value)
}

func TestTags_Add(t *testing.T) {
	testTag := TestTag("test")
	test2Tag := TestTag("test2")
	tests := []struct {
		name string
		want []Tagger
	}{
		{
			name: "can add one",
			want: []Tagger{
				enum.TagPressureMedium,
			},
		},
		{
			name: "can add two",
			want: []Tagger{
				enum.TagPressureMedium,
				enum.TagPressureExclusive,
			},
		},
		{
			name: "can add custom",
			want: []Tagger{
				&testTag,
			},
		},
		{
			name: "can add two custom",
			want: []Tagger{
				&testTag,
				&test2Tag,
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
