package ds

// import (
// 	"fmt"
// 	"testing"
// )

// func BenchmarkStringArena(b *testing.B) {
// 	a := NewStringArena()

// 	for i := 0; i < b.N; i++ {
// 		s := a.Sprintf()

// 		if len(s) != 0 {
// 			b.Errorf("Arena didnt return empty slice")
// 		}

// 		for k := range i%20 {
// 			s = append(s, k)
// 		}

// 		// fmt.Println(a.Count(), a.idx)
// 		if i%10 == 0 {
// 			a.Reset()
// 		}
// 	}
// }

// func BenchmarkSliceArena(b *testing.B) {
// 	a := NewSliceArena[int]()

// 	for i := 0; i < b.N; i++ {
// 		s := a.New()

// 		if len(s) != 0 {
// 			b.Errorf("Arena didnt return empty slice")
// 		}

// 		for k := range i%20 {
// 			s = append(s, k)
// 		}

// 		// fmt.Println(a.Count(), a.idx)
// 		if i%10 == 0 {
// 			a.Reset()
// 		}
// 	}
// }

// func BenchmarkMapArena(b *testing.B) {
// 	a := NewMapArena[int, string]()

// 	for i := 0; i < b.N; i++ {
// 		m := a.New()

// 		if len(m) != 0 {
// 			b.Errorf("Arena didnt return empty map")
// 		}

// 		for k := range i%20 {
// 			m[k] = "1"
// 		}

// 		// fmt.Println(a.Count(), a.idx)
// 		if i%10 == 0 {
// 			a.Reset()
// 		}
// 	}
// }

// // go test -run=TestUnion -bench=BenchmarkUnion . -v -benchmem -benchtime=60x
// func BenchmarkUnion(b *testing.B) {
// 	// u := NewUnion2[[4]uint8, Enum2[int16, int64]]()
// 	u := Union2[[1]uint64, int16, int64]{}

// 	var a int16
// 	var aa int64
// 	for i := 0; i < b.N; i++ {
// 		u.Put(int16(55))
// 		u.Put(int64(77))

// 		switch t := u.Get().(type) {
// 		case int16:
// 			a = t
// 		case int64:
// 			aa = t
// 		}
// 	}
// 	fmt.Println(a, aa)
// }

// func TestUnion(t *testing.T) {
// 	// u := NewUnion[[1]uint64]()
// 	// fmt.Println("A", u)
// 	// // u.Put(16)
// 	// // fmt.Println("B", u)
// 	// // u.Put(99)
// 	// // fmt.Println("C", u)

// 	// switch t := u.Get().(type) {
// 	// case int64:
// 	// 	fmt.Println("int64", t)
// 	// case int16:
// 	// 	fmt.Println("int16", t)
// 	// }

// 	//--------------------------------------------------------------------------------

// 	// u := NewUnion2[[1]uint64, BlahUnion]()
// 	// Put2(&u, int16(77))

// 	// switch t := u.Get().(type) {
// 	// case int16:
// 	// 	fmt.Println("int16", t)
// 	// case uint64:
// 	// 	fmt.Println("int64", t)
// 	// }
// 	//--------------------------------------------------------------------------------

// 	// u := NewUnion2[[4]uint8, Enum2[int16, int64]]()
// 	u := Union2[[4]uint8, int16, int64]{}
// 	u.Put(int16(55))
// 	u.Put(int64(77))

// 	switch t := u.Get().(type) {
// 	case int16:
// 		fmt.Println("int16", t)
// 	case uint64:
// 		fmt.Println("int64", t)
// 	}
// }

// // type BlahUnion struct {
// // }

// // func (u BlahUnion) GetTag(t any) uint8 {
// // 	switch t.(type) {
// // 	case int16:
// // 		return 0
// // 	case uint64:
// // 		return 1
// // 	}
// // 	panic("AAA")
// // }
