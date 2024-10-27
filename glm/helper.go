package glm

import "cmp"

func Clamp[T cmp.Ordered](low, high, val T) T {
	return min(high, max(low, val))
}
