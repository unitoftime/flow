package ds

import "testing"

func TestArrayMap(t *testing.T) {
	m := NewArrayMap[int, string]()

	// Add and check 100
	m.Put(100, "100")
	v, ok := m.Get(100)
	check(t, ok)
	compare(t, v, "100")

	// Doesn't have 50
	v, ok = m.Get(50)
	check(t, !ok)
	compare(t, v, "")

	// Add and check 50 (inside current bounds)
	m.Put(50, "50")
	v, ok = m.Get(50)
	check(t, ok)
	compare(t, v, "50")

	// Doesn't have 150
	v, ok = m.Get(150)
	check(t, !ok)
	compare(t, v, "")

	// Add and check 150 (outside current bounds)
	m.Put(150, "150")
	v, ok = m.Get(150)
	check(t, ok)
	compare(t, v, "150")

	// Iterate and check all expectations
	expectedInts := []int{
		100, 50, 150,
	}
	expectedStrings := []string{
		"100", "50", "150",
	}

	i := 0
	for k, v := range m.All() {
		compare(t, k, expectedInts[i])
		compare(t, v, expectedStrings[i])
		i++
	}

	// Deleting 20 changes nothing
	m.Delete(20)
	i = 0
	for k, v := range m.All() {
		compare(t, k, expectedInts[i])
		compare(t, v, expectedStrings[i])
		i++
	}

	// Deleting the first reorders the list
	m.Delete(50)
	expectedInts = []int{
		100, 150,
	}
	expectedStrings = []string{
		"100", "150",
	}
	i = 0
	for k, v := range m.All() {
		compare(t, k, expectedInts[i])
		compare(t, v, expectedStrings[i])
		i++
	}

	m.Delete(100)
	m.Delete(150)

	for range m.All() {
		check(t, false) // always fails, length should be 0 now
	}
}

func TestIndexMap(t *testing.T) {
	m := NewIndexMap[int, string]()

	// Add and check 100
	m.Put(100, "100")
	v, ok := m.Get(100)
	check(t, ok)
	compare(t, v, "100")

	// Doesn't have 50
	v, ok = m.Get(50)
	check(t, !ok)
	compare(t, v, "")

	// Add and check 50 (inside current bounds)
	m.Put(50, "50")
	v, ok = m.Get(50)
	check(t, ok)
	compare(t, v, "50")

	// Doesn't have 150
	v, ok = m.Get(150)
	check(t, !ok)
	compare(t, v, "")

	// Add and check 150 (outside current bounds)
	m.Put(150, "150")
	v, ok = m.Get(150)
	check(t, ok)
	compare(t, v, "150")

	// Iterate and check all expectations
	expectedInts := []int{
		50, 100, 150,
	}
	expectedStrings := []string{
		"50", "100", "150",
	}

	i := 0
	for k, v := range m.All() {
		compare(t, k, expectedInts[i])
		compare(t, v, expectedStrings[i])
		i++
	}

	// Deleting 20 changes nothing
	m.Delete(20)
	i = 0
	for k, v := range m.All() {
		compare(t, k, expectedInts[i])
		compare(t, v, expectedStrings[i])
		i++
	}

	// Deleting the first makes us skip it in the expected list
	m.Delete(50)
	i = 1
	for k, v := range m.All() {
		compare(t, k, expectedInts[i])
		compare(t, v, expectedStrings[i])
		i++
	}

}
