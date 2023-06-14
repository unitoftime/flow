package main

import (
	"github.com/unitoftime/flow/cod"
)
//go:generate go run ../

//cod:struct
type Person struct {
	Name string
	Age uint8
	Id Id
	Array [2]uint16
	Slice []uint32
	DoubleSlice [][]uint8
	// Map map[string]uint64
	Map map[string][]uint64


	MultiMap map[string]map[uint32][]uint8

	// MyUnion cod.Union //`union:"Id, Thing1, Thing2"`

	MyUnion MyUnion
}

//cod:union Id, SpecialMap
type MyUnion cod.Union

//cod:struct
type SpecialMap map[string][]uint8

// //cod:struct
// type SpecialMap2 struct {
// 	Id
// }

//cod:struct
type Id struct {
	Val uint16
}

// for unions i guess i could either have a pointer list of items. then encode type and then value. or I could do an `any` or a `cod.Encodable` (ie something that implements the cod interface. Then I just call EncodeCod() on it to marshal, and then to unmarshal I get the correct type and run DecodeCod(). Maybe create some helpers?
// //cod:union
// type Thing struct {
// 	Name *string `cod:"128`
// 	Age *uint8
// 	Id *Id
// }
