package ds

import (
	"fmt"
	"testing"
)

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func TestPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue[string]()

	{
		item := NewItem("Num5", 5)
		pq.Push(item)
	}

	{
		item := NewItem("Num7", 7)
		pq.Push(item)
	}

	{
		item := NewItem("Num1", 1)
		pq.Push(item)

		item.Priority = 6
		pq.Update(item)
	}

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := pq.Pop()
		fmt.Printf("%.2d:%s \n", item.Priority, item.Value)
	}
}

func TestPriorityMap(t *testing.T) {
	pq := NewPriorityMap[int, string]()

	pq.Put(1, "1", 1)
	pq.Put(2, "2", 2)
	pq.Put(3, "3", 3)

	pq.Put(4, "4.Priority=1", 1)

	pq.Put(5, "5.Priority=2", 2)
	pq.Put(5, "5.DIFFSTRING=2", 2)

	pq.Put(6, "5.Priority=3", 3)
	pq.Put(6, "5.Priority=4", 4)

	// Verify both lengths are equal
	fmt.Println("Length: ", pq.Len(), pq.queue.Len(), len(pq.lookup))

	// Take the items out; they arrive in decreasing priority order.
	for {
		key, val, priority, ok := pq.Pop()
		if !ok {
			break
		}

		fmt.Printf("%d:%s: %d \n", key, val, priority)
	}
}

func TestPriorityMapClear(t *testing.T) {
	pq := NewPriorityMap[int, string]()

	pq.Put(1, "1", 1)
	pq.Put(2, "2", 2)
	pq.Put(3, "3", 3)

	pq.Put(4, "4.Priority=1", 1)

	pq.Put(5, "5.Priority=2", 2)
	pq.Put(5, "5.DIFFSTRING=2", 2)

	pq.Put(6, "5.Priority=3", 3)
	pq.Put(6, "5.Priority=4", 4)

	pq.Clear()

	fmt.Println("Length: ", pq.Len())
	// Take the items out; they arrive in decreasing priority order.
	for {
		key, val, priority, ok := pq.Pop()
		if !ok {
			break
		}

		fmt.Printf("%d:%s: %d \n", key, val, priority)
		panic("SHOULD NOT GET HERE!")
	}
}
