package pgen

import (
	"fmt"
	"testing"
)

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

func TestRngTable3(t *testing.T) {
	nothing := NewTable(
		NewItem(1, "nothing"),
	)
	item := NewTable(
		NewItem(3, "a"),
		NewItem(1, "b"),
		// NewItem(1, "c"),
		// NewItem(1, "d"),
	)
	outerTable := NewTable(
		NewItem(90, nothing),
		NewItem(10, item),
	)

	samples := make(map[string]int, 0)

	for i := 0; i < 1e6; i++ {
		table := outerTable.Get()
		s := table.Get()
		samples[s] += 1
	}
	compare(samples)
}

func TestTableGetUnique(t *testing.T) {
	table := NewTable(
		NewItem(1, "a"),
		NewItem(1, "b"),
		NewItem(1, "c"),
		NewItem(1, "d"),
	)

	samples := make(map[string]int, 0)

	for i := 0; i < 1e6; i++ {
		uniques := table.GetUnique(5)
		for _, s := range uniques {
			samples[s] += 1
		}
	}
	compare(samples)
}

func TestTableGetUnique2(t *testing.T) {
	table := NewTable(
		NewItem(1, "a"),
		NewItem(1, "b"),
		NewItem(1, "c"),
		NewItem(1, "d"),
	)

	samples := make(map[string]int, 0)

	for i := 0; i < 1e6; i++ {
		uniques := table.GetUnique(2)
		for _, s := range uniques {
			samples[s] += 1
		}
	}
	compare(samples)
}

func TestTableGetUnique3(t *testing.T) {
	table := NewTable(
		NewItem(80, "a"),
		NewItem(10, "b"),
		NewItem(7, "c"),
		NewItem(3, "d"),
	)

	samples := make(map[string]int, 0)

	for i := 0; i < 1e6; i++ {
		uniques := table.GetUnique(2)
		for _, s := range uniques {
			samples[s] += 1
		}
	}
	compareUnique(samples, 2)
}

func compareUnique(samples map[string]int, count int) {
	total := 0
	for _, s := range samples {
		total += s
	}
	total = total / count

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Total:", total)
	fmt.Println("Samples:", samples)
	fmt.Println("Count:", count)
	for k, s := range samples {
		fmt.Println(k, float64(s)/float64(total))
	}
}
