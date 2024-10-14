package ds

type RingBuffer[T any] struct {
	idx     int
	readIdx int
	buffer  []T
}

func NewRingBuffer[T any](length int) *RingBuffer[T] {
	return &RingBuffer[T]{
		idx:    0,
		buffer: make([]T, length),
	}
}

func (b *RingBuffer[T]) Cap() int {
	return len(b.buffer)
}

func (b *RingBuffer[T]) Len() int {
	l := len(b.buffer)

	firstIdx := b.readIdx
	lastIdx := b.idx
	if lastIdx < firstIdx {
		lastIdx += l
	}
	return lastIdx - firstIdx

}

func (b *RingBuffer[T]) Add(t T) {
	b.buffer[b.idx] = t
	b.idx = (b.idx + 1) % len(b.buffer)

	if b.idx == b.readIdx {
		// If we just added one and the read index matches the write index, then we know that we are overwriting unread elements. So just shift the read index by one (as if we just read the one that was written)
		b.readIdx = (b.readIdx + 1) % len(b.buffer)
	}
}

// Returns the last element and false if the buffer is emptied
func (b *RingBuffer[T]) Remove() (T, bool) {
	ret := b.buffer[b.readIdx]
	newReadIdx := (b.readIdx + 1) % len(b.buffer)

	if newReadIdx == b.idx {
		// If the next index to read from is the current write index, then we know we've read the whole buffer. In this case don't increment the read index, this will cause the next Remove function to read the same value.
		return ret, false
	} else {
		// Else we do want to progress the index
		b.readIdx = newReadIdx
	}
	return ret, true
}

// TODO - Maybe convert this to an iterator
func (b *RingBuffer[T]) Buffer() []T {
	ret := make([]T, len(b.buffer))
	firstSliceLen := len(b.buffer) - b.idx
	copy(ret[:firstSliceLen], b.buffer[b.idx:len(b.buffer)])
	copy(ret[firstSliceLen:], b.buffer[0:b.idx])
	return ret
}
