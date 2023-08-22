package ds

import (
	"container/heap"
)

type mapData[K comparable,T any] struct {
	key K
	val T
}
func newMapData[K comparable, T any](key K, val T) mapData[K,T] {
	return mapData[K,T]{
		key: key,
		val: val,
	}
}


// Note: Higher value is pulled out first
type PriorityMap[K comparable, T any] struct {
	lookup map[K]*Item[mapData[K,T]]
	queue *PriorityQueue[mapData[K,T]]
}
func NewPriorityMap[K comparable, T any]() *PriorityMap[K, T] {
	return &PriorityMap[K,T]{
		lookup: make(map[K]*Item[mapData[K,T]]),
		queue: NewPriorityQueue[mapData[K,T]](),
	}
}

// Adds the item, overwriting the old one if needed
func (q *PriorityMap[K, T]) Put(key K, val T, priority int) {
	item, ok := q.lookup[key]
	if ok {
		item.Priority = priority
		item.Value.val = val
		q.queue.Update(item)
	} else {
		item := NewItem(newMapData(key, val), priority)
		q.lookup[key] = item
		q.queue.Push(item)
	}
}

// Gets the item, doesnt remove it
func (q *PriorityMap[K, T]) Get(key K) (T, int, bool) {
	item, ok := q.lookup[key]
	if ok {
		return item.Value.val, item.Priority, true
	}

	var t T
	return t, 0, false
}

// Gets the value and removes it from the queue
func (q *PriorityMap[K, T]) Remove(key K) (T, bool) {
	item, ok := q.lookup[key]
	if ok {
		q.queue.Remove(item)
		delete(q.lookup, key)
		return item.Value.val, true
	}

	var t T
	return t, false
}

func (q *PriorityMap[K, T]) Len() int {
	return q.queue.Len()
}

// Pops the highest priority item
func (q *PriorityMap[K, T]) Pop() (K, T, int, bool) {
	if q.queue.Len() <= 0 {
		var key K
		var t T
		return key, t, 0, false
	}

	item := q.queue.Pop()
	delete(q.lookup, item.Value.key)
	return item.Value.key, item.Value.val, item.Priority, true
}

// Removes all items from the queue
func (q *PriorityMap[K, T]) Clear() {
	clear(q.lookup)
	q.queue.Clear()
}

//--------------------------------------------------------------------------------

type PriorityQueue[T any] struct {
	heap heapQueue[T]
}

// Note: Higher value is pulled out first
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

func (pq *PriorityQueue[T]) Clear() {
	pq.heap.Clear()
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
	return pq[i].Priority > pq[j].Priority
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

func (pq *heapQueue[T]) Clear() {
	*pq = (*pq)[:0]
}
