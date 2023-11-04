package wrapper

type Bool struct {
	Value bool
}

func (b *Bool) Get() any {
	if b == nil {
		return nil
	}
	return b.Value
}

func (b *Bool) Set(value bool) *Bool {
	b.Value = value
	return b
}
