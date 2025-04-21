package ds

type Opt[T any] struct {
	Val T
	Has bool
}

func Optional[T any](t T) Opt[T] {
	return Opt[T]{
		Val: t,
		Has: true,
	}
}
func (o *Opt[T]) Get() (T, bool) {
	return o.Val, o.Has
}
func (o *Opt[T]) GetOrDefault(def T) T {
	if !o.Has {
		return def
	}
	return o.Val
}
func (o *Opt[T]) Set(newVal T) {
	o.Has = true
	o.Val = newVal
}
func (o *Opt[T]) Clear() {
	o.Has = false
	var t T
	o.Val = t
}
