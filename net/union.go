package net

import (
	"fmt"
	"reflect"

	"github.com/unitoftime/binary"
)

type Union struct {
	Type uint8
	Payload []byte
}

type UnionBuilder struct {
	types map[reflect.Type]uint8
	impl []any
}

func NewUnion(structs ...any) *UnionBuilder {
	if len(structs) > 256 {
		panic("TOO MANY STRUCTS")
	}

	types := make(map[reflect.Type]uint8)
	for i := range structs {
		typeStr := reflect.TypeOf(structs[i])
		types[typeStr] = uint8(i)
	}

	return &UnionBuilder {
		types: types,
		impl: structs,
	}
}

// Converts the underlying value inside the to a pointer and returns an interface for that
func valToPtr(val any) any {
	v := reflect.ValueOf(val)
	rVal := reflect.New(v.Type())
	rVal.Elem().Set(v)
	ptrVal := rVal.Interface()
	return ptrVal
}
// Converts the underlying interface with pointer to just the value
func ptrToVal(valPtr any) any {
	return reflect.Indirect(reflect.ValueOf(valPtr)).Interface()
}

func (u *UnionBuilder) Make(val any) (Union, error) {
	typeStr := reflect.TypeOf(val)
	typeId, ok := u.types[typeStr]
	if !ok {
		return Union{}, fmt.Errorf("Unknown Type: %T", val)
	}

	// TODO - can optimize the double serialize
	serializedVal, err := binary.Marshal(val)
	if err != nil {
		return Union{}, err
	}
	union := Union{
		Type: typeId,
		Payload: serializedVal,
	}
	return union, nil
}

func (u *UnionBuilder) Unmake(union Union) (any, error) {
	idx := int(union.Type)
	if idx >= len(u.impl) {
		return nil, fmt.Errorf("Unknown message opcode %d max: %d", idx, len(u.impl)-1)
	}
	val := u.impl[idx]
	valPtr := valToPtr(val)

	err := binary.Unmarshal(union.Payload, valPtr)

	return ptrToVal(valPtr), err
}

func (u *UnionBuilder) Serialize(val any) ([]byte, error) {
	union, err := u.Make(val)
	if err != nil {
		return nil, err
	}

	serializedUnion, err := binary.Marshal(union)
	return serializedUnion, err
}

func (u *UnionBuilder) Deserialize(dat []byte) (any, error) {
	union := Union{}
	err := binary.Unmarshal(dat, &union)
	if err != nil { return nil, err }

	return u.Unmake(union)
}
