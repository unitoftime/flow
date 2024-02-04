package pgen

import (
	"fmt"
	"testing"
)

func TestRngTable(t *testing.T) {
	table := NewTable(
		NewItem(1, "a"),
		NewItem(1, "b"),
		NewItem(1, "c"),
		NewItem(1, "d"),
	)
	samples := make(map[string]int, 0)

	for i := 0; i < 1e6; i++ {
		s := table.Get()
		samples[s] += 1
	}
	compare(samples)
}

func compare(samples map[string]int) {
	total := 0
	for _, s := range samples {
		total += s
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Total:", total)
	fmt.Println("Samples:", samples)
	for k, s := range samples {
		fmt.Println(k, float64(s)/float64(total))
	}
}

func TestRngTable2(t *testing.T) {
	table := NewTable(
		NewItem(1, "a"),
		NewItem(24, "b"),
		NewItem(25, "c"),
		NewItem(50, "d"),
	)
	samples := make(map[string]int, 0)

	for i := 0; i < 1e6; i++ {
		s := table.Get()
		samples[s] += 1
	}
	compare(samples)
}
