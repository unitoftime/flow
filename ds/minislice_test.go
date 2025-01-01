package ds

import (
	"fmt"
	"testing"
)

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func TestMiniSlice(t *testing.T) {
	slice := MiniSlice[[8]int, int]{}

	// for i := range slice.Array {
	// 	slice.Array[i] = i * 100
	// 	fmt.Println(slice.Array[i])
	// }

	for i, val := range slice.All() {
		fmt.Println(i, val)
	}

	if slice.Len() != 0 {
		t.Errorf("Wrong Length\n")
	}

	slice.Append(100)
	if slice.Len() != 1 {
		t.Errorf("Wrong Length\n")
	}

	slice.Clear()
	if slice.Len() != 0 {
		t.Errorf("Wrong Length\n")
	}


	for i := range 100 {
		slice.Append(i * 100)
	}

	if slice.Len() != 100 {
		t.Errorf("Wrong Length\n")
	}

	for i, val := range slice.All() {
		if val != i * 100 {
			t.Errorf("Bad Iterated Value. Expect %v. Got %v\n", i*100, val)
		}
	}
}
