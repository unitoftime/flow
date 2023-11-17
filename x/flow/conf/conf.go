package conf

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

func Decode(m map[string]any, v any) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		// DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			registeredMapHookFunc(),
		),
		Result: v,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	err = decoder.Decode(m)
	if err != nil {
		return err
	}
	// fmt.Printf("Decoded: %T: %v\n", targ.Get(), targ.Get())

	return nil
}

func registeredMapHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type,
		data interface{}) (interface{}, error) {
			if f.Kind() != reflect.Map {
				return data, nil
			}
			// fmt.Printf("Type: %T\n", data)
			// fmt.Printf("From: %s To: %s\n", f.String(), t.String())
			if t.Kind() != reflect.Interface {
				return data, nil
			}

			m, ok := data.(map[string]any)
			if !ok {
				return data, nil
			}
			if len(m) != 1 {
				return data, nil
			}
			for k, v := range m {
				fmt.Println("Special: ", k, v)
				subMap, ok := v.(map[string]any)
				if !ok { return data, nil }

				targ, ok := registry.Get(k)
				if !ok {
					return data, nil // TODO: warn? Unregistered type?
				}

				ret := targ.New()
				err := Decode(subMap, ret.Ptr())
				if err != nil {
					return data, err // TODO: return nil? or err?
				}
				fmt.Println("Special Ret: ", k, ret.Get())
				return ret.Get(), nil
			}
			return data, nil
		}
}

//--------------------------------------------------------------------------------
type decodeTargeter interface {
	New() decodeTargeter
	Get() any
	Ptr() any
}

type decodeTarget[T any] struct {
	Value T
}
func NewDT[T any](value T) *decodeTarget[T] {
	return &decodeTarget[T]{value}
}
func (t *decodeTarget[T]) New() decodeTargeter {
	return &decodeTarget[T]{t.Value}
}
func (t *decodeTarget[T]) Get() any {
	return t.Value
}
func (t *decodeTarget[T]) Ptr() any {
	return &t.Value
}

//--------------------------------------------------------------------------------
var registry = NewTypeRegistry()

type TypeRegistry struct {
	nameToTarget map[string]decodeTargeter
}
func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{
		nameToTarget: make(map[string]decodeTargeter),
	}
}
func (r *TypeRegistry) Get(name string) (decodeTargeter, bool) {
	t, ok := registry.nameToTarget[name]
	return t, ok
}
func Register[T any](name string, value T) {
	_, exists := registry.nameToTarget[name]
	if exists {
		panic(fmt.Sprintf("Duplicate Register: %s", name))
	}

	target := NewDT[T](value)
	registry.nameToTarget[name] = target
}

// Mostly for testing
func (r *TypeRegistry) Clear() {
	clear(r.nameToTarget)
}
//--------------------------------------------------------------------------------

// var registry = NewTypeRegistry()

// type TypeRegistry struct {
// 	nameToType map[string]reflect.Type
// 	typeToName map[reflect.Type]string
// }
// func NewTypeRegistry() *TypeRegistry {
// 	return &TypeRegistry{
// 		nameToType: make(map[string]reflect.Type),
// 		typeToName: make(map[reflect.Type]string),
// 	}
// }
// func (r *TypeRegistry) Register(value any) {
// 	// This is how Gob registers: https://cs.opensource.google/go/go/+/refs/tags/go1.21.1:src/encoding/gob/type.go;l=833
// 	rt := reflect.TypeOf(value)
// 	name := rt.String()
// 	star := ""
// 	if rt.Name() == "" {
// 		if rt.Kind() == reflect.Pointer {
// 			star = "*"
// 			rt = rt.Elem() // Dereference the pointer
// 		}
// 	}
// 	if rt.Name() != "" {
// 		if rt.PkgPath() == "" {
// 			name = star + rt.Name()
// 		} else {
// 			name = star + rt.PkgPath() + "." + rt.Name()
// 		}
// 	}
// 	r.RegisterName(name, value)
// }
// func (r *TypeRegistry) RegisterName(name string, value any) {
// 	rt := reflect.TypeOf(value)
// 	r.nameToType[name] = rt
// 	r.typeToName[rt] = name
// }

// func (r *TypeRegistry) TypeToName(value any) string {
// 	name, ok := r.typeToName[reflect.TypeOf(value)]
// 	if !ok {
// 		panic(fmt.Sprintf("unregistered type in interface: %T", value))
// 	}
// 	return name
// }
