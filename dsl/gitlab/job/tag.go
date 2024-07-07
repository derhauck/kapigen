package job

import (
	"gitlab.com/kateops/kapigen/dsl/enum"
)

type Tag struct {
	Value *enum.Tag
}

func NewTag(tag *enum.Tag) *Tag {
	var overrideTag enum.Tag
	if tag == nil {
		overrideTag = enum.TagPressureMedium
		return &Tag{
			Value: &overrideTag,
		}
	}
	return &Tag{
		Value: tag,
	}

}

func (c *Tag) Get() string {
	return c.Value.String()
}

type Tags []*Tag

func (c *Tags) Get() []*Tag {
	if c == nil {
		return []*Tag{}
	}
	return *c
}

func (c *Tags) Add(tag enum.Tag) *Tags {
	appended := append(c.Get(), NewTag(&tag))
	newCis := Tags(appended)
	*c = newCis
	return c
}
func (c *Tags) Render() TagYaml {
	var values []string
	if c == nil {
		return []string{NewTag(nil).Get()}
	}

	if len(c.Get()) == 0 {
		return []string{NewTag(nil).Get()}
	}
	for _, ci := range c.Get() {
		values = append(values, ci.Get())
	}
	return values
}

type TagYaml []string
