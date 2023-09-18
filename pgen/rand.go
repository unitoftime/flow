package pgen

import (
	"math/rand"
	"golang.org/x/exp/constraints"

	"github.com/unitoftime/flow/phy2"
)

type castable interface {
	constraints.Integer | constraints.Float
}

type Range[T castable] struct {
	Min, Max T
}
func (r Range[T]) Get() T {
	width := float64(r.Max) - float64(r.Min)
	return T((rand.Float64() * width) + float64(r.Min))
}

func (r Range[T]) SeededGet(rng *rand.Rand) T {
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

func RandomPositionInRect(r phy2.Rect) phy2.Pos {
	randX := Range[float64]{r.Min.X, r.Max.X}.Get()
	randY := Range[float64]{r.Min.Y, r.Max.Y}.Get()
	return phy2.Pos{randX, randY}
}

//--------------------------------------------------------------------------------
// - Tables
//--------------------------------------------------------------------------------

type Item[T any] struct{
	Weight int
	Item T
}
func NewItem[T any](weight int, item T) Item[T] {
	return Item[T]{
		Weight: weight,
		Item: item,
	}
}

type Table[T any] struct {
	Total int
	Items []Item[T]
}

func NewTable[T any](items ...Item[T]) *Table[T] {
	total := 0
	for i := range items {
		total += items[i].Weight
	}

	// TODO - Seeding?

	return &Table[T]{
		Total: total,
		Items: items, // TODO - maybe sort this. it might make the search a little faster?
	}
}

// Returns the item if successful, else returns nil
func (t *Table[T]) Get() T {
	roll := rand.Intn(t.Total + 1)

	// Essentially we just loop forward incrementing the `current` value. and once we pass it, we know that we are in that current section of the distribution.
	current := 0
	for i := range t.Items {
		current += t.Items[i].Weight
		if roll <= current {
			return t.Items[i].Item
		}
	}

	// TODO: is there a way to write this so it never fails?
	// Else just return the first item, something went wrong with the search
	return t.Items[0].Item
}
