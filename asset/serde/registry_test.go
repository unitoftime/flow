package serde

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/unitoftime/gotiny"
)

// // {
// // 	"D": {
// // 		"MyData": {
// // 			"B": "MyData - D", "D": {
// // 				"MyData3": {
// // 					"Number": 12345
// // 				}
// // 			}
// // 		}
// // 	}
// // }

// // {
// // 	"D": {
// // 		"__TYPE__": "MyData",
// // 		"B": "MyData - D",
// // 		"D": {
// // 			"__TYPE__": "MyData3"
// // 			Number: 12345,
// // 		},
// // 	}
// // }

// //--------------------------------------------------------------------------------
// type inter interface {
// 	Int() int
// }

// type MyData struct {
// 	A int
// 	B string
// 	C []any
// 	D any
// 	E []inter
// }
// func (d MyData) Int() int {
// 	return 1
// }

// type MyData2 struct {
// 	Name string
// }
// func (d MyData2) Int() int {
// 	return 2
// }

// type MyData3 struct {
// 	Number int
// }
// func (d MyData3) Int() int {
// 	return 3
// }

// func TestThing(t *testing.T) {
// // A : 555
// // B : 6666
// // C : [1 2 a b {MyData2Struct}]
// // D : {0 MyData - D [] {12345}}

// 	// registry.Clear()
// 	Register(MyData{})
// 	Register(MyData2{})
// 	Register(MyData3{})
// 	// Register([]inter{})

// 	myData := MyData{
// 		A: 555,
// 		B: "STRINGBBBBB",
// 		C: []any{1, 2, "a", "b", MyData2{"MyData2Struct"}},
// 		D: MyData{
// 			B: "MyData - D",
// 			C: []any{},
// 			D: MyData3{12345},
// 		},
// 		// D: 5,
// 		E: []inter{
// 			MyData{
// 				C: []any{},
// 				D: 5,
// 			},
// 			MyData2{},
// 			MyData3{},
// 		},
// 	}

// 	// myDataMap := make(map[string]any)
// 	// Encode(myDataMap, myData)
// 	myDataMap := Encode(myData)
// 	data, err := json.Marshal(myDataMap)
// 	if err != nil { panic(err) }
// 	fmt.Println(string(data))
// 	// printMap(myDataMap)

// 	myDataRet := MyData{}
// 	Decode(&myDataRet, myDataMap.(map[string]any))
// 	{
// 		data, err := json.Marshal(myDataRet)
// 		if err != nil { panic(err) }
// 		fmt.Println(string(data))
// 	}

// 	fmt.Println("--------------------------------------------------------------------------------")
// 	printMap(myDataMap.(map[string]any))

// 	fmt.Println("--------------------------------------------------------------------------------")
// 	fmt.Println("ReflEquals:", reflect.DeepEqual(myData, myDataRet))
// 	// fmt.Println("EncEquals: ", encEquals(myData, myDataRet))
// 	fmt.Println("MyData:    ", myData)
// 	fmt.Println("MyDataRet: ", myDataRet)

// 	j1, _ := json.Marshal(myData)
// 	j2, _ := json.Marshal(myData)
// 	fmt.Println("MyData:    ", string(j1))
// 	fmt.Println("MyDataRet: ", string(j2))
// 	fmt.Println("--------------------------------------------------------------------------------")

// 	// data, err := json.Marshal(myData)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// m := make(map[string]any)
// 	// err := yaml.Unmarshal([]byte(data), &m)
// 	// if err != nil { panic(err) }
// 	// printMap(m)

// 	// result := MyData{}
// 	// err = Decode(m, &result)
// 	// if err != nil { panic(err) }

// 	// fmt.Println("Result:", result)
// 	// fmt.Printf("Type: %T\n", result.D)
// }

// // func TestMerge(t *testing.T) {
// // 	registry.Clear()
// // 	Register("MyData", MyData{})
// // 	Register("MyData2", MyData2{})
// // 	Register("MyData3", MyData3{})

// // 	data := `
// // A: 5
// // B: "hello"
// // C:
// //   - {MyData2: {Name: "slicestring"}}
// //   - {MyData3: {Number: 77}}
// // D: {MyData2: {Name: "secondstring"}}
// // `

// // 	data2 := `
// // B: "hellothere-overwrittenbydata2"
// // D: {MyData2: {Name: "overwrittenbydata2"}}
// // `

// // 	m := make(map[string]any)
// // 	err := yaml.Unmarshal([]byte(data), &m)
// // 	if err != nil { panic(err) }
// // 	fmt.Println("Before")
// // 	printMap(m)

// // 	yaml.Unmarshal([]byte(data2), &m)
// // 	fmt.Println("After")
// // 	printMap(m)

// // 	result := MyData{}
// // 	err = Decode(m, &result)
// // 	if err != nil { panic(err) }

// // 	fmt.Println("Result:", result)
// // 	fmt.Printf("Type: %T\n", result.D)
// // }

// func printMap(m map[string]any) {
// 	for k, v := range m {
// 		switch t := v.(type) {
// 		case map[string]any:
// 			fmt.Printf(k + " . ")
// 			printMap(t)
// 			fmt.Printf("\n")
// 		case []any:
// 			fmt.Printf("[ ")
// 			for i := range t {
// 				fmt.Printf("%d (%T): %v, ", i, t[i], t[i])
// 			}
// 			fmt.Printf("]\n")
// 		default:
// 			fmt.Printf("%s (%T): %v\n", k, v, v)
// 		}
// 	}
// }

// func TestCodec(t *testing.T) {
// 	myData := MyData{
// 		A: 555,
// 		B: "STRINGBBBBB",
// 		// C: []any{1, 2, "a", "b", MyData2{"MyData2Struct"}},
// 		// D: 5,
// 		D: MyData{
// 			B: "MyData - D",
// 			// C: []any{},
// 			// D: MyData3{12345},
// 		},
// 		// // D: 5,
// 		// E: []inter{
// 		// 	MyData{
// 		// 		C: []any{},
// 		// 		D: 5,
// 		// 	},
// 		// 	MyData2{},
// 		// 	MyData3{},
// 		// },
// 	}

// 	var jh codec.JsonHandle
// 	h := &jh

// 	err := h.SetInterfaceExt(reflect.TypeOf(MyData{}), 0, GenericExt[MyData]{})
// 	if err != nil { panic(err) }
// 	// err = h.SetInterfaceExt(reflect.TypeOf(MyData2{}), 0, GenericExt[MyData2]{})
// 	// if err != nil { panic(err) }

// 	var b []byte = make([]byte, 0, 64)
// 	enc := codec.NewEncoderBytes(&b, h)
// 	err = enc.Encode(myData)
// 	if err != nil { panic(err) }

// 	fmt.Println("Encoded:", string(b))

// 	myDataRet := MyData{}
// 	dec := codec.NewDecoderBytes(b, h)
// 	err = dec.Decode(&myDataRet)
// 	if err != nil { panic(err) }

// 	fmt.Println("--------------------------------------------------------------------------------")
// 	fmt.Println("--------------------------------------------------------------------------------")
// 	fmt.Println("ReflEquals:", reflect.DeepEqual(myData, myDataRet))
// 	// fmt.Println("EncEquals: ", encEquals(myData, myDataRet))
// 	fmt.Println("MyData:    ", myData)
// 	fmt.Println("MyDataRet: ", myDataRet)

// 	j1, _ := json.Marshal(myData)
// 	j2, _ := json.Marshal(myDataRet)
// 	fmt.Println("MyData:    ", string(j1))
// 	fmt.Println("MyDataRet: ", string(j2))
// 	fmt.Println("--------------------------------------------------------------------------------")
// }

// type GenericExt[T any] struct {

// }

// func (g GenericExt[T]) ConvertExt(v any) any {
// 	fmt.Printf("ConvertExt %T\n", v)
// 	rv := reflect.ValueOf(v)
// 	if rv.Kind() == reflect.Pointer {
// 		rv = rv.Elem()
// 	}

// 	m := make(map[string]any)
// 	m["TYPE"] = rv.Type().String() // TODO: Should come from registration
// 	switch rv.Kind() {
// 	case reflect.Struct:
// 		numFields := rv.NumField()
// 		for i := 0; i < numFields; i++ {
// 			sf := rv.Type().Field(i)
// 			m[sf.Name] = rv.Field(i).Interface()
// 		}
// 	default:
// 		panic(fmt.Sprintf("Unhandled Kind: %s", rv.Kind()))
// 	}

// 	return m
// }
// func (g GenericExt[T]) UpdateExt(dst any, src any) {
// 	fmt.Printf("UpdateExt dst: %T | src: %T\n", dst, src)
// 	{
// 		a := dst.(*MyData)
// 		fmt.Println(a)
// 	}

// 	rv := reflect.ValueOf(dst)
// 	if rv.Kind() == reflect.Pointer {
// 		rv = rv.Elem()
// 	}

// 	switch rv.Kind() {
// 	case reflect.Struct:
// 		numFields := rv.NumField()
// 		for i := 0; i < numFields; i++ {
// 			sf := rv.Type().Field(i)
// 			mapVal, ok := src.(map[string]any)[sf.Name]
// 			if !ok { panic("AAA") }
// 			fmt.Println("CANSET", rv.Field(i).CanSet())
// 			fmt.Printf("DATA: %T\n", mapVal)
// 			if mapVal == nil { continue }
// 			rv.Field(i).Set(reflect.ValueOf(mapVal))
// 		}
// 	default:
// 		panic(fmt.Sprintf("Unhandled Kind: %s", rv.Kind()))
// 	}

// }

//--------------------------------------------------------------------------------
type inter interface {
	Int() int
}

type MyData struct {
	A int
	B string
	C []any
	D any
	E []inter
}
func (d MyData) Int() int {
	return 1
}

type MyData2 struct {
	Name string
}
func (d MyData2) Int() int {
	return 2
}

type MyData3 struct {
	Number int
}
func (d MyData3) Int() int {
	return 3
}

func TestGotiny(t *testing.T) {
// A : 555
// B : 6666
// C : [1 2 a b {MyData2Struct}]
// D : {0 MyData - D [] {12345}}

	// registry.Clear()
	gotiny.Register(MyData{})
	gotiny.Register(MyData2{})
	gotiny.Register(MyData3{})
	// Register([]inter{})

	myData := MyData{
		A: 555,
		B: "STRINGBBBBB",
		C: []any{1, 2, "a", "b", MyData2{"MyData2Struct"}},
		D: MyData{
			B: "MyData - D",
			C: []any{},
			D: MyData3{12345},
		},
		// D: 5,
		E: []inter{
			MyData{
				C: []any{},
				D: 5,
			},
			MyData2{},
			MyData3{},
		},
	}

	// myDataMap := make(map[string]any)
	// Encode(myDataMap, myData)
	data := gotiny.Marshal(&myData)
	// if err != nil { panic(err) }
	fmt.Println("Marshalled:", string(data))
	// printMap(myDataMap)

	myDataRet := MyData{}
	gotiny.Unmarshal(data, &myDataRet)
	{
		data, err := json.Marshal(myDataRet)
		if err != nil { panic(err) }
		fmt.Println(string(data))
	}

	fmt.Println("--------------------------------------------------------------------------------")
	// printMap(myDataMap.(map[string]any))

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("ReflEquals:", reflect.DeepEqual(myData, myDataRet))
	// fmt.Println("EncEquals: ", encEquals(myData, myDataRet))
	fmt.Println("MyData:    ", myData)
	fmt.Println("MyDataRet: ", myDataRet)

	j1, _ := json.Marshal(myData)
	j2, _ := json.Marshal(myData)
	fmt.Println("MyData:    ", string(j1))
	fmt.Println("MyDataRet: ", string(j2))
	fmt.Println("--------------------------------------------------------------------------------")

	// data, err := json.Marshal(myData)
	// if err != nil {
	// 	panic(err)
	// }

	// m := make(map[string]any)
	// err := yaml.Unmarshal([]byte(data), &m)
	// if err != nil { panic(err) }
	// printMap(m)

	// result := MyData{}
	// err = Decode(m, &result)
	// if err != nil { panic(err) }

	// fmt.Println("Result:", result)
	// fmt.Printf("Type: %T\n", result.D)
}
