package images

type PullPolicy int

const (
	Always PullPolicy = iota
	IfNotPresent
	Never
)

var values = map[PullPolicy]string{
	Always:       "always",
	IfNotPresent: "if-not-present",
	Never:        "never",
}

func (c PullPolicy) String() string {
	if v, ok := values[c]; ok {
		return v
	}

	return values[Always]
}
