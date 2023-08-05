package pgen

import (
	"math/rand"
	"golang.org/x/exp/constraints"
)

type castable interface {
	constraints.Integer | constraints.Float
}

type Range[T castable] struct {
	Min, Max T
}
func (r *Range[T]) Get() T {
	width := float64(r.Max) - float64(r.Min)
	return T((rand.Float64() * width) + float64(r.Min))
}

func (r *Range[T]) SeededGet(rng *rand.Rand) T {
	width := float64(r.Max) - float64(r.Min)
	return T((rng.Float64() * width) + float64(r.Min))
}

// TODO: Should I separate int from float?
// type RngRange[T constraints.Integer]struct{
// 	Min, Max T
// }
// func NewRngRange[T constraints.Integer](min, max T) RngRange[T] {
// 	return RngRange[T]{min, max}
// }
// func (r RngRange[T]) Roll() T {
// 	delta := r.Max - r.Min
// 	if delta <= 0 {
// 		return r.Min
// 	}
// 	return T(rand.Intn(int(delta))) + r.Min
// }
