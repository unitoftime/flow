package prefab

import (
	"testing"
	"fmt"
	"bytes"
	// "maps"
	// "reflect"

	// // "os"
	// "bytes"
	// "encoding/gob"
	"github.com/BurntSushi/toml"
	"github.com/mitchellh/mapstructure"
	"reflect"
	// "github.com/jinzhu/copier"
	// "github.com/fatih/structs"
	// // "github.com/r3labs/diff/v3"

	"github.com/unitoftime/ecs"
	"github.com/unitoftime/glitch"
	"github.com/unitoftime/flow/phy2"
	"github.com/unitoftime/flow/asset"

)

type Interfacer interface {
	Thing()
}
type DoubleAny struct {
	V any
}
func (d DoubleAny) Thing() {
	fmt.Print("HERE")
}
type SubSubType struct {
	Sub SubType
}
func (d SubSubType) Thing() {
	fmt.Print("HERE")
}
type SubType struct {
	Z int
}
func (d SubType) Thing() {
	fmt.Print("HERE")
}
type ComplexType struct {
	Blah int
	X SubType
	A, B, C any
	Dub DoubleAny
	SubSub SubSubType
}

func RegisteredInterfaceMapHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
			fmt.Println("INTERFACE HOOK:", f, t, data)
			return data, nil

		// if f.Kind() != reflect.String {
		// 	return data, nil
		// }
		// result := reflect.New(t).Interface()
		// unmarshaller, ok := result.(encoding.TextUnmarshaler)
		// if !ok {
		// 	return data, nil
		// }
		// if err := unmarshaller.UnmarshalText([]byte(data.(string))); err != nil {
		// 	return nil, err
		// }
		// return result, nil
	}
}

// func RegisteredInterfaceMapHookFunc() mapstructure.DecodeHookFunc {
// 	return func(f reflect.Value, t reflect.Value) (any, error) {
// 		fmt.Println("INTERFACE HOOK:", f, t)
// 		if f.Kind() != reflect.Interface {
// 			return f.Interface(), nil
// 		}
// 		fmt.Println("AAAAAAAAAAAAAA")

// 		// var i interface{} = struct{}{}
// 		// if t.Type() != reflect.TypeOf(&i).Elem() {
// 		// 	return f.Interface(), nil
// 		// }

// 		// m := make(map[string]interface{})
// 		// t.Set(reflect.ValueOf(m))

// 		return f.Interface(), nil
// 	}
// }

// func TestMapStruct(t *testing.T) {
// 	input := ComplexType{
// 		Blah: 1,
// 		X: SubType{2},
// 		A: "hello",
// 		B: 3,
// 		C: SubType{4},
// 		Dub: DoubleAny{SubType{5}},
// 		SubSub: SubSubType{SubType{6}},
// 	}
// 	output := make(map[string]any)
// 	config := &mapstructure.DecoderConfig{
// 		Result:   &output,
// 		DecodeHook: RegisteredInterfaceMapHookFunc(),
// 		// DecodeHook: mapstructure.RecursiveStructToMapHookFunc(),
// 	}

// 	decoder, err := mapstructure.NewDecoder(config)
// 	if err != nil {
// 		panic(err)
// 	}

// 	decoder.Decode(input)
// 	fmt.Println(output)
// 	printMap(output)
// }

// 1. new Boxed component that has a method to convert to ecs.C()
// 2. main prefab holder object that has a list or map of this new boxed component
// 3. Diff by diffing on a field by field in the map or list
// 4. Merge by merging on a field by field in the map or list
// 5. wrapper thing can do type embedding in the serdes

func TestBlah(t *testing.T) {
	base := CompList{
		"phy2.Pos": Boxed(phy2.Pos{1, 2}),
		"sprite": Boxed(asset.Handle[glitch.Sprite]{Name:"gopher.png"}),
	}
	prefab := CompList{
		"phy2.Pos": Boxed(phy2.Pos{1, 200}),
		"phy2.Vel": Boxed(phy2.Vel{33, 44}),
	}

	final := mergeCompList(base, prefab)
	for k, v := range final {
		fmt.Println(k, ":", v)
	}
}

func init() {
	registry.Register(MyType{})
	registry.Register(OtherData{})
	registry.Register(ecs.C(phy2.Pos{}))
	registry.Register(ecs.C(phy2.Vel{}))
	registry.Register(ecs.C(asset.Handle[glitch.Sprite]{}))
}

func marshalToml(v any) string {
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	err := encoder.Encode(v)
	if err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}

// func ToMap(v any) map[string]string {
// 	ret := make(map[string]string)
// }

// func TestSimple(t *testing.T) {
// 	a := 5
// 	fmt.Println(marshalToml(a))
// 	b := "abcd"
// 	fmt.Println(marshalToml(b))
// }

// func TestPrefabDiff(t *testing.T) {
// 	base := PrefabAsset3{
// 		Id: 0x45e6874534cc94d9,
// 		Name: "base",
// 		Comp: make(map[string]any),
// 	}

// 	base.Add(
// 		OtherData{123},
// 		ecs.C(phy2.Pos{1, 2}),
// 		ecs.C(asset.Handle[glitch.Sprite]{
// 			Name: "gopher.png",
// 		}), // TODO: spritedata or something (cant run on server)
// 	)

// 	prefab := PrefabAsset3{
// 		Id: 0x3be5b087cd038114,
// 		Name: "prefab",
// 		Comp: make(map[string]any),
// 	}
// 	prefab.Add(
// 		OtherData{123},
// 		ecs.C(phy2.Pos{1, 222}),
// 		ecs.C(phy2.Vel{100, 200}),
// 		ecs.C(asset.Handle[glitch.Sprite]{
// 			Name: "gopher.png",
// 		}), // TODO: spritedata or something (cant run on server)
// 	)

// 	baseCopy := PrefabAsset3{}
// 	err := copier.CopyWithOption(&baseCopy, base, copier.Option{DeepCopy: true})
// 	if err != nil { panic(err) }

// 	prefabCopy := PrefabAsset3{}
// 	err = copier.CopyWithOption(&prefabCopy, prefab, copier.Option{DeepCopy: true})
// 	if err != nil { panic(err) }

// 	fmt.Println(baseCopy)
// 	fmt.Println(prefabCopy)

// 	fmt.Println("--------------------------------------------------------------------------------")
// 	baseMap := make(map[string]any)
// 	decoder, err := mapstructure.NewDecoder(
// 		&mapstructure.DecoderConfig{
// 			DecodeHook: mapstructure.RecursiveStructToMapHookFunc(),
// 			Result: &baseMap,
// 		},
// 	)
// 	if err != nil { panic(err) }
// 	err = decoder.Decode(baseCopy)
// 	if err != nil { panic(err) }
// 	printMap(baseMap)

// 	// baseMap := make(map[string]any)
// 	// err = mapstructure.Decode(baseCopy, &baseMap)
// 	// if err != nil { panic(err) }
// 	// printMap(baseMap)


// 	fmt.Println("--------------------------------------------------------------------------------")
// }

// func TestMapStructure(t *testing.T) {
// 	myType := MyType{
// 		Id: 1234,
// 		Data: map[string]any{
// 			"A": 5,
// 			"B": map[string]any{
// 				"A": 6,
// 				"B1": "abcd",
// 			},
// 			"C": []any{
// 				1, 2, 3,
// 				"a", "b", "c",
// 			},
// 		},
// 	}
// 	input := MyType{}
// 	err := copier.CopyWithOption(&input, myType, copier.Option{DeepCopy: true})
// 	if err != nil { panic(err) }

// 	output := make(map[string]any)
// 	err = mapstructure.Decode(input, &output)
// 	if err != nil { panic(err) }
// 	printMap(output)
// 	myType.Data["A"] = 6

// 	myType2 := MyType{}
// 	mapstructure.Decode(output, &myType2)

// 	fmt.Println("T1", myType)
// 	fmt.Println("T2", myType2)
// 	fmt.Println("DeepEqual: ", reflect.DeepEqual(myType, myType2))
// }
// func TestTypelessSerdes(t *testing.T) {
// 	myType := MyType{
// 		Id: 1234,
// 		Data: map[string]any{
// 			"A": 5,
// 			"B": map[string]any{
// 				"A": 6,
// 				"B1": "abcd",
// 			},
// 		},
// 	}
// }

type MyType struct {
	Id uint64
	Data map[string]any
}

// func TestConstruction(t *testing.T) {
// 	base := MyType{
// 		1234,
// 		Data: map[string]any{
// 			"A": 5,
// 			"B": structMap{
// 				"A": 6,
// 				"B1": 7,
// 			},
// 		},
// 	}

// 	baseMap := Encode(base)
// 	printMap(baseMap)
// }

// func TestMapMerge(t *testing.T) {
// 	base := structMap{
// 		"A": 5,
// 		"B": structMap{
// 			"A": 6,
// 			"B1": 7,
// 		},
// 	}

// 	over := structMap{
// 		"A": 5,
// 		"A1": 9999,
// 		"B": structMap{
// 			"A": 6666,
// 			"C": 8888,
// 		},
// 		"C": structMap{
// 			"DD": "hello there",
// 		},
// 	}

// 	final := mergeMaps(base, over)
// 	printMap(final)
// }

// func TestToMap(t *testing.T) {
// 	base := PrefabAsset3{
// 		Id: 0x45e6874534cc94d9,
// 		Name: "base",
// 		Comp: make(map[string]any),
// 	}

// 	base.Add(
// 		OtherData{123},
// 		ecs.C(phy2.Pos{1, 2}),
// 		ecs.C(asset.Handle[glitch.Sprite]{
// 			Name: "gopher.png",
// 		}), // TODO: spritedata or something (cant run on server)
// 	)

// 	prefab := PrefabAsset3{
// 		Id: 0x3be5b087cd038114,
// 		Name: "prefab",
// 		Comp: make(map[string]any),
// 	}
// 	prefab.Add(
// 		OtherData{123},
// 		ecs.C(phy2.Pos{1, 222}),
// 		ecs.C(phy2.Vel{100, 200}),
// 		ecs.C(asset.Handle[glitch.Sprite]{
// 			Name: "gopher.png",
// 		}), // TODO: spritedata or something (cant run on server)
// 	)

// 	baseMap := structs.Map(base)
// 	prefabMap := structs.Map(prefab)
// 	fmt.Println("------BaseMap")
// 	printMap(baseMap)
// 	fmt.Println("------PrefabMap")
// 	printMap(prefabMap)

// 	// fmt.Println("------Retype")
// 	// retypeMapData(prefabMap)
// 	fmt.Println("--------ToMap")
// 	ToMap(prefab)
// }

type OtherData struct {
	Value int
}
// func TestPrefabDiff(t *testing.T) {
// 	base := PrefabAsset3{
// 		Id: 0x45e6874534cc94d9,
// 		Name: "base",
// 		Comp: make(map[string]any),
// 	}

// 	base.Add(
// 		OtherData{123},
// 		ecs.C(phy2.Pos{1, 2}),
// 		ecs.C(asset.Handle[glitch.Sprite]{
// 			Name: "gopher.png",
// 		}), // TODO: spritedata or something (cant run on server)
// 	)

// 	prefab := PrefabAsset3{
// 		Id: 0x3be5b087cd038114,
// 		Name: "prefab",
// 		Comp: make(map[string]any),
// 	}
// 	prefab.Add(
// 		OtherData{123},
// 		ecs.C(phy2.Pos{1, 222}),
// 		ecs.C(phy2.Vel{100, 200}),
// 		ecs.C(asset.Handle[glitch.Sprite]{
// 			Name: "gopher.png",
// 		}), // TODO: spritedata or something (cant run on server)
// 	)

// 	diff := DiffPrefab(base, prefab)
// 	fmt.Println("--------------------------------------------------------------------------------")
// 	fmt.Println("--------------------------------------------------------------------------------")
// 	printMap(diff)
// 	fmt.Println("BASE: ", base)

// 	MergePrefab(base, diff)
// 	fmt.Println("================================================================================")
// 	fmt.Println("BASE: ", base)

// 	// printMap(diff)
// }

func printMap(m map[string]any) {
	for k, v := range m {
		switch t := v.(type) {
		case map[string]any:
			fmt.Printf(k + " . ")
			printMap(t)
			fmt.Printf("\n")
		default:
			fmt.Println(k, ":", v)
		}
	}
}
