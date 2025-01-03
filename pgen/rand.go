package pgen

import (
	"math/rand"
	"slices"

	"github.com/unitoftime/flow/glm"
	"golang.org/x/exp/constraints"
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

// Pick a random item out of a list
func GetList[T any](list []T) (T, bool) {
	if len(list) <= 0 {
		var t T
		return t, false
	}

	return list[rand.Intn(len(list))], true
}

// Returns a random element of the list, based on the provided rng
func SeededList[T any](rng *rand.Rand, list []T) T {
	return list[rng.Intn(len(list))]
}

// func ListItem[T any](rng *rand.Rand, list []T) (T, bool) {
// 	if len(list) <= 0 {
// 		var t T
// 		return t, false
// 	}

// 	return list[rng.Intn(len(list))], true
// }

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

func RandomPositionInRect(r glm.Rect) glm.Vec2 {
	randX := Range[float64]{r.Min.X, r.Max.X}.Get()
	randY := Range[float64]{r.Min.Y, r.Max.Y}.Get()
	return glm.Vec2{randX, randY}
}

//--------------------------------------------------------------------------------
// - Tables
//--------------------------------------------------------------------------------

type Item[T any] struct {
	Weight int
	Item   T
}

func NewItem[T any](weight int, item T) Item[T] {
	return Item[T]{
		Weight: weight,
		Item:   item,
	}
}

type Table[T any] struct {
	Total int
	Items []Item[T]
}

func NewTable[T any](items ...Item[T]) *Table[T] {
	// TODO - Seeding?
	t := &Table[T]{
		Items: items, // TODO - maybe sort this. it might make the search a little faster?
	}
	t.regenerate()

	return t
}

func (t *Table[T]) regenerate() {
	total := 0
	for i := range t.Items {
		if t.Items[i].Weight <= 0 {
			continue
		} // Skip if the weight of this item is <= 0

		total += t.Items[i].Weight
	}
	t.Total = total
}

func (t *Table[T]) getIndex() int {
	if t.Total == 0 {
		t.regenerate()
	}
	roll := rand.Intn(t.Total)

	// Essentially we just loop forward incrementing the `current` value. and once we pass it, we know that we are in that current section of the distribution.
	current := 0
	for i := range t.Items {
		current += t.Items[i].Weight
		if roll < current {
			return i
		}
	}

	// Else just return the first item, something went wrong with the search
	// TODO: is this okay? Or should I return a bool and handle it further up?
	return 0
}

// Returns the item rolled
func (t *Table[T]) Get() T {
	index := t.getIndex()

	// TODO: is there a way to write this so it never fails?
	return t.Items[index].Item
}

// TODO: Needs testing
// // Returns the item rolled and removes it from the table
// func (t *Table[T]) GetAndRemove() (T, bool) {
// 	var ret T
// 	if len(t.Items) <= 0 {
// 		return ret, false
// 	}

// 	if t.Total == 0 {
// 		t.regenerate()
// 	}

// 	roll := rand.Intn(t.Total)

// 	// Essentially we just loop forward incrementing the `current` value. and once we pass it, we know that we are in that current section of the distribution.
// 	current := 0
// 	idx := -1
// 	for i := range t.Items {
// 		current += t.Items[i].Weight
// 		if roll < current {
// 			idx = i
// 		}
// 	}

// 	if idx < 0 {
// 		// If we couldn't find the index for some reason, then it fails
// 		return ret, false
// 	}

// 	// Get Item
// 	ret = t.Items[idx].Item

// 	// Remove Item and regenerate
// 	t.Items[idx] = t.Items[len(t.Items)-1]
// 	t.Items = t.Items[:len(t.Items)-1]
// 	t.regenerate()

// 	return ret, true
// }

// Returns returns count unique items, if there are less items in the loot table, only returns what is available to satisfy the uniqueness
func (t *Table[T]) GetUnique(count int) []T {
	if count <= 0 {
		return []T{}
	}
	ret := make([]T, 0, count)

	// If there are less items than we are requesting, then just return them all
	if count >= len(t.Items) {
		for i := range t.Items {
			ret = append(ret, t.Items[i].Item)
		}
		return ret
	}

	indexes := make([]int, 0, count)
	for i := 0; i < count; i++ {
		idx := t.getIndex()
		if slices.Contains(indexes, idx) {
			// If we have already found this index, try again
			i-- // Note: You are guaranteed that count < len(t.Items) so this will exit
			continue
		}

		indexes = append(indexes, idx)
	}

	for _, idx := range indexes {
		ret = append(ret, t.Items[idx].Item)
	}

	return ret
}
