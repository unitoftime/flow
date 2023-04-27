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
