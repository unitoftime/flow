package ds

import (
	"runtime"
	"testing"
)

// --------------------------------------------------------------------------------
// TODO: Centralize these testing helpers
// --------------------------------------------------------------------------------
func check(t *testing.T, b bool) {
	if !b {
		_, f, l, _ := runtime.Caller(1)
		t.Errorf("%s:%d - checked boolean is false!", f, l)
	}
}

// Check two things match, if they don't, throw an error
func compare[T comparable](t *testing.T, actual, expected T) {
	if expected != actual {
		_, f, l, _ := runtime.Caller(1)
		t.Errorf("%s:%d - actual(%v) did not match expected(%v)", f, l, actual, expected)
	}
}

//--------------------------------------------------------------------------------

func TestMiniSlice(t *testing.T) {
	slice := MiniSlice[[8]int, int]{}

	compare(t, slice.Len(), 0)

	slice.Append(100)
	compare(t, slice.Len(), 1)

	slice.Clear()
	compare(t, slice.Len(), 0)

	for i := range 100 {
		slice.Append(i * 100)
	}

	compare(t, slice.Len(), 100)

	for i, val := range slice.All() {
		compare(t, val, i*100)
	}

	compare(t, slice.Find(1*100), 1)
	compare(t, slice.Find(2*100), 2)
	compare(t, slice.Find(19*100), 19)
	compare(t, slice.Find(99*100), 99)
	compare(t, slice.Find(100*100), -1) // Doesn't exist

	compare(t, slice.Get(7), 7*100)
	compare(t, slice.Get(8), 8*100)
	compare(t, slice.Get(9), 9*100)

	slice.Set(7, 777)
	slice.Set(8, 888)
	slice.Set(9, 999)

	compare(t, slice.Get(7), 777)
	compare(t, slice.Get(8), 888)
	compare(t, slice.Get(9), 999)

	slice.Delete(7)
	compare(t, slice.Get(7), 99*100)

	slice.Delete(8)
	compare(t, slice.Get(8), 98*100)

	slice.Delete(9)
	compare(t, slice.Get(9), 97*100)
}

func TestMiniSliceSmall(t *testing.T) {
	slice := MiniSlice[[8]int, int]{}

	slice.Append(0)
	slice.Append(1)
	slice.Append(2)

	slice.Delete(0)
	slice.Append(3)

	compare(t, slice.Get(0), 2)
	compare(t, slice.Get(1), 1)
	compare(t, slice.Get(2), 3)
}

func TestMiniSliceSkipDelete(t *testing.T) {
	slice := MiniSlice[[4]int, int]{}

	slice.Append(8)

	slice.Delete(-1)
	slice.Delete(2)

	compare(t, slice.Len(), 1)
	slice.Delete(0)
	compare(t, slice.Len(), 0)
}
