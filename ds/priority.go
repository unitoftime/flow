package ds

import (
	"container/heap"
)

type PriorityQueue[T any] struct {
	heap heapQueue[T]
}

func NewPriorityQueue[T any]() *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		heap: make(heapQueue[T], 0),
	}

	heap.Init(&pq.heap)
	return pq
}

func (pq *PriorityQueue[T]) Len() int {
	return pq.heap.Len()
}

func (pq *PriorityQueue[T]) Push(item *Item[T]) {
	heap.Push(&pq.heap, item)
}

func (pq *PriorityQueue[T]) Pop() *Item[T] {
	item := heap.Pop(&pq.heap).(*Item[T])
	return item
}

func (pq *PriorityQueue[T]) Update(item *Item[T]) {
	heap.Fix(&pq.heap, item.index)
}

func (pq *PriorityQueue[T]) Remove(item *Item[T]) {
	heap.Remove(&pq.heap, item.index)
}

// An Item is something we manage in a priority queue.
type Item[T any] struct {
	// The index is needed by update and is maintained by the heap.Interface methods.
	index    int // The index of the item in the heap.
	Priority int // The priority of the item in the queue.

	Value    T   // The value of the item; arbitrary.
}
func NewItem[T any](value T, priority int) *Item[T] {
	return &Item[T]{
		Value: value,
		Priority: priority,
	}
}


// TODO: You could make this faster by replacing this with something hand-made

// A heap implements heap.Interface and holds Items.
type heapQueue[T any] []*Item[T]

func (pq heapQueue[T]) Len() int { return len(pq) }

func (pq heapQueue[T]) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq heapQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *heapQueue[T]) Push(x any) {
	n := len(*pq)
	item := x.(*Item[T])
	item.index = n
	*pq = append(*pq, item)
}

func (pq *heapQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

