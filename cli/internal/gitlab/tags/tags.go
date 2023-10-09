package tags

type Tag int

const (
	PRESSURE_MEDIUM Tag = iota
	PRESSURE_BUILDKIT
)

var values = map[Tag]string{
	PRESSURE_MEDIUM:   "pressure::medium",
	PRESSURE_BUILDKIT: "pressure::buildkit",
}

func (c *Tag) Tag() string {
	if value, ok := values[*c]; ok {
		return value
	}
	return values[PRESSURE_MEDIUM]
}

type Ci struct {
	Value *Tag
}

func New(tag *Tag) *Ci {
	var overrideTag Tag
	if tag == nil {
		overrideTag = PRESSURE_MEDIUM
		return &Ci{
			Value: &overrideTag,
		}
	}
	return &Ci{
		Value: tag,
	}

}

func (c *Ci) Get() string {
	return c.Value.Tag()
}

type Cis []*Ci

func (c *Cis) Get() []*Ci {
	if c == nil {
		return []*Ci{}
	}
	return *c
}

func (c *Cis) Add(tag Tag) *Cis {
	appended := append(c.Get(), New(&tag))
	newCis := Cis(appended)
	c = &newCis
	return c
}
func (c *Cis) Render() Yaml {
	var values []string
	if c == nil {
		return []string{New(nil).Get()}
	}

	if len(c.Get()) == 0 {
		return []string{New(nil).Get()}
	}
	for _, ci := range c.Get() {
		values = append(values, ci.Get())
	}
	return values
}

type Yaml []string
