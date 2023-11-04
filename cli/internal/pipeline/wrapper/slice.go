package wrapper

type Slice[T any] interface {
	Add(str T) *Slice[T]
	AddSeveral(strSlice []T) *Slice[T]
	Get() []T
}
