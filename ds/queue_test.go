package ds

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	q := NewQueue[string](10)

	q.Add("a")
	fmt.Println(q.Len())
	r, ok := q.Remove()
	fmt.Println(r, ok)
}
