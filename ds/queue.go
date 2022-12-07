package ds

type Queue[T any] struct {
	Buffer []T
	ReadIdx int
	WriteIdx int
}

func NewQueue[T any](length int) *Queue[T] {
	return &Queue[T]{
		Buffer: make([]T, length),
		ReadIdx: 0,
		WriteIdx: 0,
	}
}

func (q *Queue[T]) Len() int {
	l := len(q.Buffer)

	firstIdx := q.ReadIdx
	lastIdx := q.WriteIdx
	if lastIdx < firstIdx {
		lastIdx += l
	}
	return lastIdx - firstIdx
}

func (q *Queue[T]) Add(t T) {
	if (q.WriteIdx + 1) % len(q.Buffer) == q.ReadIdx {
		panic("QUEUE IS FULL!") // TODO - Not sure how to handle this, maybe just keep writing like I would in a ringbuffer?
	}
	q.Buffer[q.WriteIdx] = t
	q.WriteIdx = (q.WriteIdx + 1) % len(q.Buffer)
}

func (q *Queue[T]) Peek() (T, bool) {
	if q.ReadIdx == q.WriteIdx {
		var ret T
		return ret, false
	}
	return q.Buffer[q.ReadIdx], true
}
func (q *Queue[T]) PeekLast() (T, bool) {
	if q.ReadIdx == q.WriteIdx {
		var ret T
		return ret, false
	}
	idx := (q.WriteIdx + len(q.Buffer) - 1) % len(q.Buffer)
	return q.Buffer[idx], true
}
func (q *Queue[T]) Remove() (T, bool) {
	if q.ReadIdx == q.WriteIdx {
		var ret T
		return ret, false
	}
	val := q.Buffer[q.ReadIdx]
	q.ReadIdx = (q.ReadIdx + 1) % len(q.Buffer)
	return val, true
}

// func (n *NextTransform) Map(fn func(t ServerTransform)) {
// 	if n.ReadIdx == n.WriteIdx {
// 		return // Empty
// 	}

// 	l := len(n.Transforms)
// 	firstIdx := n.ReadIdx
// 	// lastIdx := n.WriteIdx
// 	lastIdx := (n.WriteIdx + len(n.Transforms) - 1) % len(n.Transforms)

// 	cnt := 0
// 	// TODO - this might be simpler in two loops?
// 	for i := firstIdx; i != lastIdx; i=(i + 1) % l {
// 		fn(n.Transforms[i])
// 		cnt++
// 	}
// 	// log.Print("Mapped: ", cnt)
// }
