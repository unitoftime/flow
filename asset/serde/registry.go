package serde

// // Maybe?
// // var hashTable = crc64.MakeTable(crc64.ISO)
// // func crc(label string) uint64 {
// // 	return crc64.Checksum([]byte(label), hashTable)
// // }

// //--------------------------------------------------------------------------------
// var registry *TypeRegistry
// func init() {
// 	registry = NewTypeRegistry()
// 	Register(int(0))
// 	Register(int8(0))
// 	Register(int16(0))
// 	Register(int32(0))
// 	Register(int64(0))
// 	Register(uint(0))
// 	Register(uint8(0))
// 	Register(uint16(0))
// 	Register(uint32(0))
// 	Register(uint64(0))
// 	Register(float32(0))
// 	Register(float64(0))
// 	Register(complex64(0i))
// 	Register(complex128(0i))
// 	Register(uintptr(0))
// 	Register(false)
// 	Register("")
// 	Register([]byte(nil)) // TODO: Same as []uint8
// 	Register([]int(nil))
// 	Register([]int8(nil))
// 	Register([]int16(nil))
// 	Register([]int32(nil))
// 	Register([]int64(nil))
// 	Register([]uint(nil))
// 	// Register([]uint8(nil)) // TODO: Same as []byte
// 	Register([]uint16(nil))
// 	Register([]uint32(nil))
// 	Register([]uint64(nil))
// 	Register([]float32(nil))
// 	Register([]float64(nil))
// 	Register([]complex64(nil))
// 	Register([]complex128(nil))
// 	Register([]uintptr(nil))
// 	Register([]bool(nil))
// 	Register([]string(nil))
// }

// type TypeRegistry struct {
// 	typeToName map[reflect.Type]string
// 	nameToTarget map[string]decodeTargeter
// }
// func NewTypeRegistry() *TypeRegistry {
// 	return &TypeRegistry{
// 		nameToTarget: make(map[string]decodeTargeter),
// 		typeToName: make(map[reflect.Type]string),
// 	}
// }
// func (r *TypeRegistry) Get(name string) (decodeTargeter, bool) {
// 	t, ok := registry.nameToTarget[name]
// 	return t, ok
// }

// func (r *TypeRegistry) fromType(val any) (string, bool) {
// 	name, ok := registry.typeToName[reflect.TypeOf(val)]
// 	return name, ok
// }

// func register[T any](value T) {
// 	// fmt.Println("Register", reflect.TypeOf(value).String())
// 	registerName(reflect.TypeOf(value).String(), value)
// }

// func registerName[T any](name string, value T) {
// 	_, exists := registry.nameToTarget[name]
// 	if exists {
// 		panic(fmt.Sprintf("Duplicate Register: %s", name))
// 	}

// 	fmt.Println("Register", name)
// 	target := NewDT[T](value)
// 	registry.nameToTarget[name] = target
// 	registry.typeToName[reflect.TypeOf(value)] = name
// }

// // // Mostly for testing
// // func (r *TypeRegistry) Clear() {
// // 	clear(r.nameToTarget)
// // 	clear(r.typeToName)
// // }

// //--------------------------------------------------------------------------------
// type decodeTargeter interface {
// 	New() decodeTargeter
// 	Get() any
// 	Ptr() any
// 	Type() reflect.Type
// }

// type decodeTarget[T any] struct {
// 	Value T
// }
// func NewDT[T any](value T) *decodeTarget[T] {
// 	return &decodeTarget[T]{value}
// }
// func (t *decodeTarget[T]) New() decodeTargeter {
// 	return &decodeTarget[T]{t.Value}
// }
// func (t *decodeTarget[T]) Get() any {
// 	return t.Value
// }
// func (t *decodeTarget[T]) Ptr() any {
// 	return &t.Value
// }
// func (t *decodeTarget[T]) Type() reflect.Type {
// 	return reflect.TypeOf(t.Value)
// }
