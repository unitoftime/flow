package ds

// Safely adds the value at the slice index provided
func GrowAdd[K any](slice []K, idx int, val K) []K {
	requiredLength := idx + 1
	growAmount := requiredLength - len(slice)
	if growAmount <= 0 {
		slice[idx] = val
		return slice // No need to grow if the sliceIdx is already in bounds
	}

	slice = append(slice, make([]K, growAmount)...)
	slice[idx] = val
	return slice
}

// Safely gets and returns a value from a slice, and a boolen to indicate boundscheck
func SafeGet[K any](slice []K, idx int) (K, bool) {
	if idx < 0 || idx >= len(slice) {
		var k K
		return k, false
	}
	return slice[idx], true
}
