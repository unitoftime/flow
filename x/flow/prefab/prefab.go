package prefab

import (
	"fmt"
	// "maps"

	// // "os"
	// "bytes"
	// "encoding/gob"
	// "github.com/BurntSushi/toml"
	"github.com/mitchellh/mapstructure"
	// "github.com/fatih/structs"
	// // "github.com/r3labs/diff/v3"
	"reflect"
	// "github.com/jinzhu/copier"
	"github.com/unitoftime/ecs"
)

type Mergable interface {
	MergeIn(any) Mergable
}

type CompList map[string]Mergable

func mergeCompList(base, over CompList) CompList {
	ret := make(CompList)
	for k, v := range base {
		ret[k] = v
	}
	for k, v := range over {
		existing, has := ret[k]
		if has {
			ret[k] = existing.MergeIn(v)
		} else {
			ret[k] = v
		}
	}

	return ret
}

type TypedMap struct {
	Type string
	Map map[string]any
}

type Box[T any] struct {
	Value T
}
func Boxed[T any](t T) Box[T] {
	return Box[T]{t}
}
func (b Box[T]) ToComponent() ecs.Component {
	return ecs.C(b.Value)
}

func (b Box[T]) ToMap() TypedMap {
	name := registry.TypeToName(b.Value)
	ret := make(map[string]any)
	err := mapstructure.Decode(b.Value, &ret)
	if err != nil { panic(err) }
	return TypedMap{
		Type: name,
		Map: ret,
	}
}

func (t TypedMap)FromMap() any {
	registry.NameToType(t.Type)
	err := mapstructure.Decode(b.Value, &ret)
	if err != nil { panic(err) }
	return ret
}

// Merges in the passed in type, provided it is the same type
func (b Box[T]) MergeIn(v any) Mergable {
	b2, match := v.(Box[T])
	if !match { return b }

	// Merge function for these types
	base := make(map[string]any)
	err := mapstructure.Decode(b, &base)
	if err != nil { panic(err) }

	over := make(map[string]any)
	err = mapstructure.Decode(b2, &over)
	if err != nil { panic(err) }

	return b
}

// func Marshal(v any) ([]byte, error) {
// 	dat := make([]byte, 0)
// 	value := reflect.ValueOf(v)
// 	switch value.Kind() {
// 	case reflect.Struct:

// 	default:
// 	}
// }


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

// func ToMap(base any) map[string]any {
// 	ret := make(map[string]any)
// 	switch base.Kind() {
// 	case reflect.Struct:
// 		for i := 0; i < base.NumField(); i++ {
// 			fieldName := base.Type().Field(i).Name
// 			subMap := ToMap(base.Field(i))
// 			ret[fieldName
// 			subDiff := findDiff(base.Field(i), over.Field(i))
// 			if subDiff != nil {
// 				diff[fieldName] = subDiff
// 			}
// 		}
// 		return diff
// 		// case reflect.Slice:
// 		// 	for i := 0; i < base.Len(); i++ {
// 		// 		findDiff(base.Field(i), over.Field(i), path+"."+"slice", diff)
// 		// 	}
// 	case reflect.Map:
// 		diff := make(map[string]any)
// 		iter := over.MapRange()
// 		for iter.Next() {
// 			k := iter.Key()
// 			overVal := iter.Value()
// 			baseVal := base.MapIndex(k)

// 			subDiff := findDiff(baseVal, overVal)
// 			if subDiff != nil {
// 				diff[k.String()] = subDiff
// 			}
// 		}
// 		return diff

// 	case reflect.Interface:
// 		if base.IsNil() && !over.IsNil() {
// 			return over.Interface()
// 		} else if !base.IsNil() && over.IsNil() {
// 			// This means that the type is missing from the over, which means they wanted to delete it
// 			return Deleted{}
// 		} else if !base.IsNil() && !over.IsNil() {
// 			subDiff := findDiff(base.Elem(), over.Elem())
// 			return subDiff
// 		}
// 	case reflect.Ptr:
// 		if base.IsNil() && !over.IsNil() {
// 			return over.Interface()
// 		} else if !base.IsNil() && over.IsNil() {
// 			// This means that the type is missing from the over, which means they wanted to delete it
// 			return Deleted{}
// 		} else if !base.IsNil() && !over.IsNil() {
// 			subDiff := findDiff(base.Elem(), over.Elem())
// 			return subDiff
// 		}
// 	default:
// 		fmt.Println("Default: ", base.Kind().String(), base, over)
// 		if base.CanInterface() && over.CanInterface() {
// 			if !reflect.DeepEqual(base.Interface(), over.Interface()) {
// 				return over.Interface()
// 			}
// 		}
// 	}
// }

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
	mergeDiff(reflect.ValueOf(base), diff)
	return base
}

func mergeDiff(base reflect.Value, diff map[string]any) {
	// baseCopy := reflect.New(base.Type()).Interface()
	// err := copier.CopyWithOption(baseCopy, base, copier.Option{DeepCopy: true})
	// if err != nil { panic(err) }
	// value := reflect.ValueOf(baseCopy.Elem())


	for k, v := range diff {
		fmt.Println("Loop: ", k, v, base.Kind().String())
		switch base.Kind() {
		case reflect.Struct:
			field := base.FieldByName(k)
			fmt.Println("Field: ", k, field)
			if field.Kind() == reflect.Struct {
				mergeDiff(field, v.(map[string]any))
			}
		default:
			fmt.Println("default", k, v)
		}
	}

	// switch base.Kind() {
	// case reflect.Struct:
	// 	for i := 0; i < base.NumField(); i++ {
	// 		fieldName := base.Type().Field(i).Name
	// 		diffVal, has := diff[fieldName]
	// 		if has {
	// 			field := base.Field(i)
	// 			if field.CanSet() {
	// 				field.Set(reflect.ValueOf(diffVal))
	// 			}
	// 		}
	// 	}
	// 	return base
	// case reflect.Map:
	// case reflect.Interface:
	// case reflect.Ptr:
	// }
	// return base
}


type typed struct {
	Name string
	Data any
}
// func newTyped(v any) typed {
// 	return typed{
// 		Name: registry.TypeToName(v),
// 		Data: v,
// 	}
// }
// func retypeMapData(base map[string]any) map[string]typed {
// 	ret := make(map[string]typed)
// 	for k, v := range base {
// 		switch t := v.(type) {
// 		case map[string]any:
// 			subMap := retypeMapData(t)
// 			if subMap != nil {
// 				ret[k] = subMap
// 			}
// 		default:
// 			fmt.Printf("Default: %v %T\n", v, v)
// 		}
// 	}
// 	return nil
// }


// func SM(v any) typed {
// 	t := reflect.ValueOf(v)
// 	return ToMap(t)
// }

// func ToMap(base reflect.Value) typed {
// 	ret := make(map[string]typed)
// 	switch base.Kind() {
// 	case reflect.Struct:
// 		for i := 0; i < base.NumField(); i++ {
// 			fieldName := base.Type().Field(i).Name
// 			subMap := ToMap(base.Field(i))
// 			if subMap != nil {
// 				ret[fieldName] = subMap
// 			}
// 		}
// 		// case reflect.Slice:
// 		// 	for i := 0; i < base.Len(); i++ {
// 		// 		findDiff(base.Field(i), over.Field(i), path+"."+"slice", diff)
// 		// 	}
// 	case reflect.Map:
// 		iter := over.MapRange()
// 		for iter.Next() {
// 			k := iter.Key()
// 			overVal := iter.Value()
// 			baseVal := base.MapIndex(k)

// 			subMap := ToMap(baseVal)
// 			if subMap != nil {
// 				ret[k.String()] = subMap
// 			}
// 		}
// 	case reflect.Interface:
// 		if base.IsNil() {
// 			return nil
// 		} else {
// 			return ToMap(base.Elem())
// 		}
// 	case reflect.Ptr:
// 		if base.IsNil() {
// 			return nil
// 		} else {
// 			return ToMap(base.Elem())
// 		}
// 	default:
// 		fmt.Println("ToMapLeaf: ", base.Kind().String(), base)
// 		if base.CanInterface() {
// 			return newTyped(base.Interface())
// 		}
// 	}
// }


// type typedAny map[string]any
// // TODO: Marshal/Unmarshal functions
// func TA(v any) typedAny {
// 	name := registry.TypeToName(v)
// 	return typedAny{
// 		name: v,
// 	}
// }



// type structMap map[string]typedAny // FieldName -> typedAny
// func mergeMaps(base, over structMap) structMap {
// 	final := make(structMap)
// 	// TODO You need to recursively do this on every key
// 	maps.Copy(final, base)
// 	maps.Copy(final, over)
// 	return final
// }

// func newStructMap(v any) structMap {
// 	ret := make(structMap)
// 	// TODO: typeswitch actually
// 	t := reflect.TypeOf(v)
// 	switch t.Kind() {
// 	case reflect.Struct:
// 		for i := 0; i < t.NumField(); i++ {
// 			fieldName := t.Type().Field(i).Name
// 			subMap := ToMap(t.Field(i))
// 			if subMap != nil {
// 				ret[fieldName] = subMap
// 			}
// 		}
// 	}

// 	return ret
// }

// type typedData map[string]any
// func typed(v any) typedData {
// 	name := registry.TypeToName(v)
// 	return typedData{
// 		name: v,
// 	}
// }

// func Encode(v any) any {
// 	value := reflect.ValueOf(v)
// 	switch value.Kind() {
// 	case reflect.Struct:
// 		ret :=  make(map[string]typedData)
// 		return ret
// 	default:
// 		if value.CanInterface() {
			
// 		}
// 	}

// 	return nil // Nothing to encode
// }

