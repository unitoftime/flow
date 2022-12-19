package ds

type RingBuffer[T any] struct {
	idx int
	buffer []T
}
func NewRingBuffer[T any](length int) *RingBuffer[T] {
	return &RingBuffer[T]{
		idx: 0,
		buffer: make([]T, length),
	}
}

func (b *RingBuffer[T]) Len() int {
	return len(b.buffer)
}

func (b *RingBuffer[T]) Add(t T) {
	b.buffer[b.idx] = t
	b.idx = (b.idx + 1) % len(b.buffer)
}

// TODO - Maybe convert this to an iterator
func (b *RingBuffer[T]) Buffer() []T {
	ret := make([]T, len(b.buffer))
	firstSliceLen := len(b.buffer) - b.idx
	copy(ret[:firstSliceLen], b.buffer[b.idx:len(b.buffer)])
	copy(ret[firstSliceLen:], b.buffer[0:b.idx])
	return ret
}
