package main

import (
	"testing"
	"fmt"

	// // "os"
	// "bytes"
	// "encoding/gob"
	// // "github.com/BurntSushi/toml"
	// // "github.com/mitchellh/mapstructure"
	// "github.com/fatih/structs"
	// // "github.com/r3labs/diff/v3"
	"reflect"


	"github.com/unitoftime/ecs"
	"github.com/unitoftime/glitch"
	"github.com/unitoftime/flow/phy2"
	"github.com/unitoftime/flow/asset"

)

var registry = NewTypeRegistry()

type TypeRegistry struct {
	nameToType map[string]reflect.Type
	typeToName map[reflect.Type]string
}
func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{
		nameToType: make(map[string]reflect.Type),
		typeToName: make(map[reflect.Type]string),
	}
}
func (r *TypeRegistry) Register(value any) {
	// This is how Gob registers: https://cs.opensource.google/go/go/+/refs/tags/go1.21.1:src/encoding/gob/type.go;l=833
	rt := reflect.TypeOf(value)
	name := rt.String()
	star := ""
	if rt.Name() == "" {
		if rt.Kind() == reflect.Pointer {
			star = "*"
			rt = rt.Elem() // Dereference the pointer
		}
	}
	if rt.Name() != "" {
		if rt.PkgPath() == "" {
			name = star + rt.Name()
		} else {
			name = star + rt.PkgPath() + "." + rt.Name()
		}
	}
	r.RegisterName(name, value)
}
func (r *TypeRegistry) RegisterName(name string, value any) {
	rt := reflect.TypeOf(value)
	r.nameToType[name] = rt
	r.typeToName[rt] = name
}

func (r *TypeRegistry) TypeToName(value any) string {
	name, ok := r.typeToName[reflect.TypeOf(value)]
	if !ok {
		panic(fmt.Sprintf("unregistered type in interface: %T", value))
	}
	return name
}

// Used to mark something as deleted in a diff
type Deleted struct {}

type PrefabAsset3 struct {
	Id uint64
	Name string
	Comp map[string]any
}
func (p *PrefabAsset3) Add(comps ...any) {
	for _, c := range comps {
		name := registry.TypeToName(c)
		p.Comp[name] = c
	}
}

func ToMap(base any) map[string]any {
	ret := make(map[string]any)
	switch base.Kind() {
	case reflect.Struct:
		for i := 0; i < base.NumField(); i++ {
			fieldName := base.Type().Field(i).Name
			subMap := ToMap(base.Field(i))
			ret[fieldName
			subDiff := findDiff(base.Field(i), over.Field(i))
			if subDiff != nil {
				diff[fieldName] = subDiff
			}
		}
		return diff
		// case reflect.Slice:
		// 	for i := 0; i < base.Len(); i++ {
		// 		findDiff(base.Field(i), over.Field(i), path+"."+"slice", diff)
		// 	}
	case reflect.Map:
		diff := make(map[string]any)
		iter := over.MapRange()
		for iter.Next() {
			k := iter.Key()
			overVal := iter.Value()
			baseVal := base.MapIndex(k)

			subDiff := findDiff(baseVal, overVal)
			if subDiff != nil {
				diff[k.String()] = subDiff
			}
		}
		return diff

	case reflect.Interface:
		if base.IsNil() && !over.IsNil() {
			return over.Interface()
		} else if !base.IsNil() && over.IsNil() {
			// This means that the type is missing from the over, which means they wanted to delete it
			return Deleted{}
		} else if !base.IsNil() && !over.IsNil() {
			subDiff := findDiff(base.Elem(), over.Elem())
			return subDiff
		}
	case reflect.Ptr:
		if base.IsNil() && !over.IsNil() {
			return over.Interface()
		} else if !base.IsNil() && over.IsNil() {
			// This means that the type is missing from the over, which means they wanted to delete it
			return Deleted{}
		} else if !base.IsNil() && !over.IsNil() {
			subDiff := findDiff(base.Elem(), over.Elem())
			return subDiff
		}
	default:
		fmt.Println("Default: ", base.Kind().String(), base, over)
		if base.CanInterface() && over.CanInterface() {
			if !reflect.DeepEqual(base.Interface(), over.Interface()) {
				return over.Interface()
			}
		}
	}
}

func DiffPrefab(base, over PrefabAsset3) map[string]any {
	return findDiff(reflect.ValueOf(base), reflect.ValueOf(over)).(map[string]any)
}

func findDiff(base, over reflect.Value) any {
	if base.Kind() != over.Kind() {
		// If mismatched kinds, then just replace the entire thing
		return over.Interface()
	}

	switch base.Kind() {
	case reflect.Struct:
		diff := make(map[string]any)
		for i := 0; i < base.NumField(); i++ {
			fieldName := base.Type().Field(i).Name
			subDiff := findDiff(base.Field(i), over.Field(i))
			if subDiff != nil {
				diff[fieldName] = subDiff
			}
		}
		return diff
	// case reflect.Slice:
	// 	for i := 0; i < base.Len(); i++ {
	// 		findDiff(base.Field(i), over.Field(i), path+"."+"slice", diff)
	// 	}
	case reflect.Map:
		diff := make(map[string]any)
		iter := over.MapRange()
		for iter.Next() {
			k := iter.Key()
			overVal := iter.Value()
			baseVal := base.MapIndex(k)

			subDiff := findDiff(baseVal, overVal)
			if subDiff != nil {
				diff[k.String()] = subDiff
			}
		}
		return diff

	case reflect.Interface:
		if base.IsNil() && !over.IsNil() {
			return over.Interface()
		} else if !base.IsNil() && over.IsNil() {
			// This means that the type is missing from the over, which means they wanted to delete it
			return Deleted{}
		} else if !base.IsNil() && !over.IsNil() {
			subDiff := findDiff(base.Elem(), over.Elem())
			return subDiff
		}
	case reflect.Ptr:
		if base.IsNil() && !over.IsNil() {
			return over.Interface()
		} else if !base.IsNil() && over.IsNil() {
			// This means that the type is missing from the over, which means they wanted to delete it
			return Deleted{}
		} else if !base.IsNil() && !over.IsNil() {
			subDiff := findDiff(base.Elem(), over.Elem())
			return subDiff
		}
	default:
		fmt.Println("Default: ", base.Kind().String(), base, over)
		if base.CanInterface() && over.CanInterface() {
			if !reflect.DeepEqual(base.Interface(), over.Interface()) {
				return over.Interface()
			}
		}
	}

	// Else, they must've matched. Just return nil
	return nil
}

func MergePrefab(base PrefabAsset3, diff map[string]any) PrefabAsset3 {
	return mergeDiff(reflect.ValueOf(base), diff).Interface().(PrefabAsset3)
}

func mergeDiff(base reflect.Value, diff map[string]any) reflect.Value {
	switch base.Kind() {
	case reflect.Struct:
		for i := 0; i < base.NumField(); i++ {
			fieldName := base.Type().Field(i).Name
			diffVal, has := diff[fieldName]
			if has {
				field := base.Field(i)
				if field.CanSet() {
					field.Set(reflect.ValueOf(diffVal))
				}
			}
		}
		return base
	case reflect.Map:
		
	case reflect.Interface:
		
	case reflect.Ptr:
		
	}
	return 
}

type OtherData struct {
	Value int
}
func TestPrefabDiff(t *testing.T) {
	registry.Register(OtherData{})
	registry.Register(ecs.C(phy2.Pos{}))
	registry.Register(ecs.C(phy2.Vel{}))
	registry.Register(ecs.C(asset.Handle[glitch.Sprite]{}))

	base := PrefabAsset3{
		Id: 0x45e6874534cc94d9,
		Name: "base",
		Comp: make(map[string]any),
	}

	base.Add(
		OtherData{123},
		ecs.C(phy2.Pos{1, 2}),
		ecs.C(asset.Handle[glitch.Sprite]{
			Name: "gopher.png",
		}), // TODO: spritedata or something (cant run on server)
	)

	prefab := PrefabAsset3{
		Id: 0x3be5b087cd038114,
		Name: "prefab",
		Comp: make(map[string]any),
	}
	prefab.Add(
		OtherData{123},
		ecs.C(phy2.Pos{1, 222}),
		ecs.C(phy2.Vel{100, 200}),
		ecs.C(asset.Handle[glitch.Sprite]{
			Name: "gopher.png",
		}), // TODO: spritedata or something (cant run on server)
	)

	diff := DiffPrefab(base, prefab)
	fmt.Println("--------------------------------------------------------------------------------")
	printMap(diff)

	// merged := MergeDiff(base, diff)
}

func printMap(m map[string]any) {
	for k, v := range m {
		switch t := v.(type) {
		case map[string]any:
			fmt.Printf(k + " . ")
			printMap(t)
		default:
			fmt.Println(k, ":", v)
		}
	}
}






// // func TestBasicPrefab(t *testing.T) {
// // 	components := []any{
// // 		phy2.Pos{1, 2},
// // 		phy2.Vel{3, 4},
// // 		asset.Handle[glitch.Sprite]{
// // 			Name: "gopher.png",
// // 		}, // TODO: spritedata or something (cant run on server)
// // 	}

// // 	prefab := PrefabAsset{
// // 		Id: "abcd",
// // 		Override: EntConfig{
// // 			Data: components,
// // 		},
// // 	}

// // 	buf := new(bytes.Buffer)
// // 	encoder := toml.NewEncoder(buf)
// // 	err := encoder.Encode(prefab)
// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	os.WriteFile("test.toml", buf.Bytes(), 0644)
// // }

// type PrefabAsset2 struct {
// 	Id string
// 	// Base *asset.Handle[PrefabAsset]
// 	Override EntConfig
// }

// type Bundle struct {
// 	Components []ecs.Component
// }

// func TestBasicPrefab2(t *testing.T) {
// 	gob.Register(ecs.C(phy2.Pos{}))
// 	gob.Register(ecs.C(phy2.Vel{}))
// 	gob.Register(ecs.C(asset.Handle[glitch.Sprite]{}))
// 	gob.Register(map[string]any{})
// 	gob.Register([]ecs.Component{})

// 	c1 := []ecs.Component{
// 		ecs.C(phy2.Pos{1, 2}),
// 		ecs.C(asset.Handle[glitch.Sprite]{
// 			Name: "gopher.png",
// 		}), // TODO: spritedata or something (cant run on server)
// 	}

// 	c2 := []ecs.Component{
// 		ecs.C(phy2.Pos{1, 101}),
// 		ecs.C(phy2.Vel{3, 4}),
// 		ecs.C(asset.Handle[glitch.Sprite]{
// 			Name: "gopher.png",
// 		}), // TODO: spritedata or something (cant run on server)
// 	}

// 	base := PrefabAsset2{
// 		Id: "base",
// 		Override: EntConfig{
// 			Components: c1,
// 		},
// 	}

// 	prefab := PrefabAsset2{
// 		Id: "instance",
// 		Override: EntConfig{
// 			Components: c2,
// 		},
// 	}

// 	fmt.Println("Base: ", base)
// 	fmt.Println("Prefab: ", prefab)

// 	// baseMap := StructToMap(base)
// 	// prefabMap := StructToMap(prefab)
// 	// baseMap := structs.Map(base)
// 	// prefabMap := structs.Map(prefab)
// 	// fmt.Println(baseMap)
// 	// fmt.Println(prefabMap)

// 	// {
// 	// 	diff := diffMaps(baseMap, prefabMap)
// 	// 	fmt.Println("Diff:", diff)

// 	// 	buf := new(bytes.Buffer)
// 	// 	encoder := toml.NewEncoder(buf)
// 	// 	err := encoder.Encode(diff)
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// 	os.WriteFile("diff.toml", buf.Bytes(), 0644)
// 	// }


// 	changelog, _ := diff.Diff(base, prefab)
// 	fmt.Println("Change: ", changelog)

// 	c := make(map[string]any)
// 	patchlog := diff.Patch(changelog, &c)
// 	fmt.Println("Patch: ", patchlog)
// 	fmt.Println("Final: ", c)

// 	// {
// 	// 	fmt.Println(mapStruct)

// 	// 	buf := new(bytes.Buffer)
// 	// 	encoder := gob.NewEncoder(buf)
// 	// 	err := encoder.Encode(mapStruct)
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}

// 	// 	os.WriteFile("test2.gob", buf.Bytes(), 0644)
// 	// }
// }

// func diffMaps(base, over map[string]any) map[string]any {
// 	diff := make(map[string]any)

// 	for k, overVal := range over {
// 		baseVal, ok := base[k]
// 		if !ok {
// 			// If base didnt have this field, then it must be fully added in the over
// 			diff[k] = overVal
// 			continue
// 		}

// 		switch overMapType := overVal.(type) {
// 		case map[string]any:
// 			baseMapType, matches := baseVal.(map[string]any)
// 			if !matches {
// 				// If baseVal is the wrong type, then it must be replaced
// 				diff[k] = overVal
// 				continue
// 			} else {
// 				// Else we found a submap, so we need to recurse
// 				subDiff := diffMaps(baseMapType, overMapType)
// 				if len(subDiff) > 0 {
// 					diff[k] = subDiff
// 				}
// 			}

// 		default:
// 			fmt.Println("DiffMapType: ", reflect.TypeOf(baseVal), reflect.TypeOf(overVal))
// 			if !reflect.DeepEqual(baseVal, overVal) {
// 				// If the types dont match then we need to replace the base with the diff
// 				diff[k] = overVal
// 				continue
// 			}
// 		}
// 	}

// 	// TODO: loop through base to find deleted fields

// 	return diff
// }

// // func StructToMap(input interface{}) map[string]interface{} {
// // 	result := make(map[string]interface{})
// // 	structValue := reflect.ValueOf(input)
// // 	structType := structValue.Type()

// // 	for i := 0; i < structValue.NumField(); i++ {
// // 		field := structValue.Field(i)
// // 		fieldType := structType.Field(i)

// // 		// Check if the field is exportable (starts with an uppercase letter)
// // 		if fieldType.PkgPath == "" {
// // 			fieldName := fieldType.Name

// // 			// Check if the field is a struct and should be recursively converted
// // 			if field.Kind() == reflect.Struct {
// // 				result[fieldName] = StructToMap(field.Interface())
// // 			} else if field.Kind() == reflect.Slice { // TODO: Or array?
// // 				length := field.Len()
// // 				for i := 0; i < length; i++ {
// // 					sliceVal := field.Index(i)
					
// // 				}
// // 			} else {
// // 				fmt.Println("Else: ", fieldName, field.Interface())
// // 				// For non-struct fields, convert to interface{}
// // 				result[fieldName] = field.Interface()
// // 			}
// // 		}
// // 	}

// // 	return result
// // }

// func getMap(v any) map[string]any {
// 	m := structs.Map(v)

// 	buf := new(bytes.Buffer)
// 	encoder := gob.NewEncoder(buf)
// 	err := encoder.Encode(m)
// 	if err != nil {
// 		panic(err)
// 	}

// 	mapStruct := make(map[string]any)
// 	dec := gob.NewDecoder(buf)
// 	err = dec.Decode(mapStruct)
// 	if err == nil {
// 		return nil
// 	}
// 	return mapStruct
// }

