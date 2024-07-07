package cache

type Policy int

const (
	Pull Policy = iota
	Push
	PullPush
)

var values = map[Policy]string{
	Pull:     "pull",
	Push:     "push",
	PullPush: "pull-push",
}

func (c Policy) String() string {
	if v, ok := values[c]; ok {
		return v
	}

	return values[Pull]
}
