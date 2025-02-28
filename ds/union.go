package ds

// import (
// 	"fmt"
// 	"unsafe"
// )

// type String []uint8

// type StringArena struct {
// 	// idx int
// 	// list []uint8 // TODO: you could probably just slice out sections of one really big byte slice

// 	arena *SliceArena[byte]
// }

// func NewStringArena() *StringArena {
// 	return &StringArena{
// 		arena: NewSliceArena[byte](),
// 	}
// }

// func (a *StringArena) Sprintf(str string, args ...any) String {
// 	str := a.New([]byte{})

// 	// TODO: This doesn't reallly work bc the users slice gets reallocated dynamically but not put back into the slice, so you need to pass outward a pointer
// }

// func (a *StringArena) New(dat []byte) String {
// 	cachedDat := a.arena.New()
// 	cachedDat = append(cachedDat, dat...)
// 	return String(cachedDat)
// }

// func (a *StringArena) Reset() {
// 	a.arena.Reset()
// }

// //--------------------------------------------------------------------------------

// type SliceArena[T any] struct {
// 	idx int
// 	list [][]T
// }

// func NewSliceArena[T any]() *SliceArena[T] {
// 	return &SliceArena[T]{
// 		list: make([][]T, 0),
// 	}
// }

// // func (a *Arena[T]) Count() int {
// // 	return len(a.list)
// // }

// func (a *SliceArena[T]) New() []T {
// 	// TODO: This doesn't reallly work bc the users slice gets reallocated dynamically but not put back into the slice, so you need to pass outward a pointer

// 	if a.idx < len(a.list) {
// 		i := a.idx
// 		a.idx++

// 		return a.list[i][:0]
// 	}
// 	slice := make([]T, 0)
// 	a.list = append(a.list, slice)
// 	a.idx++
// 	return a.list[len(a.list) - 1][:0]
// }

// func (a *SliceArena[T]) Reset() {
// 	a.idx = 0
// }

// //--------------------------------------------------------------------------------

// type MapArena[K comparable, V any] struct {
// 	idx int
// 	list []map[K]V
// }

// func NewMapArena[K comparable, V any]() *MapArena[K, V] {
// 	return &MapArena[K, V]{
// 		list: make([]map[K]V, 0),
// 	}
// }

// // func (a *Arena[T]) Count() int {
// // 	return len(a.list)
// // }

// func (a *MapArena[K, V]) New() map[K]V {
// 	if a.idx < len(a.list) {
// 		i := a.idx
// 		a.idx++

// 		m := a.list[i]
// 		for k := range m {
// 			delete(m, k)
// 		}
// 		return m
// 	}
// 	m := make(map[K]V)
// 	a.list = append(a.list, m)
// 	a.idx++
// 	return m
// }

// func (a *MapArena[K, V]) Reset() {
// 	a.idx = 0
// }

// // type Union[T any] struct {
// // 	data T
// // 	// tag uint8
// // }

// // func NewUnion[T any]() Union[T] {
// // 	return Union[T]{}
// // }

// // func Put[T, A any](u *Union[T], val A) {
// // 	*(*A)(unsafe.Pointer(&u.data)) = val
// // }

// // func Get[T, A any](u *Union[T]) A {
// // 	return *(*A)(unsafe.Pointer(&u.data))
// // }

// // type Unionizer interface {
// // 	// GetTag(any) uint8
// // 	Encode(unsafe.Pointer, any) uint8
// // 	Decode(uint8, unsafe.Pointer) any
// // }

// // type Union2[T any, U Unionizer] struct {
// // 	data T
// // 	tag uint8
// // }
// // func NewUnion2[T any, U Unionizer]() Union2[T, U] {
// // 	return Union2[T, U]{}
// // }

// // func (u *Union2[T, U]) Put(val any) {
// // 	var e U
// // 	u.tag = e.Encode(unsafe.Pointer(&u.data), val)
// // }

// // func (u *Union2[T, U]) Get() any {
// // 	var e U
// // 	return e.Decode(u.tag, unsafe.Pointer(&u.data))
// // }

// // // func Put2[T any, U Unionizer[T], A any](u *Union2[T, U], val A) {
// // // 	var def U
// // // 	u.tag = def.GetTag(val)
// // // 	*(*A)(unsafe.Pointer(&u.data)) = val
// // // }

// // func Get2[T, A any](u *Union[T]) any {
// // 	return *(*A)(unsafe.Pointer(&u.data))
// // }

// // type Enum2[A, B any] struct {
// // }
// // // func (e Enum2[A, B]) Tag(t any) uint8 {
// // // 	switch t.(type) {
// // // 	case A:
// // // 		return 0
// // // 	case B:
// // // 		return 1
// // // 	}
// // // 	panic("AAA")
// // // }

// // func (e Enum2[A, B]) Decode(tag uint8, data unsafe.Pointer) any {
// // 	switch tag {
// // 	case 0:
// // 		return *(*A)(data)
// // 		// return *(*A)(unsafe.Pointer(&data))
// // 	case 1:
// // 		return *(*B)(data)
// // 		// return *(*B)(unsafe.Pointer(&data))
// // 	}
// // 	panic("AAA")
// // }

// // func (e Enum2[A, B]) Encode(data unsafe.Pointer, val any) uint8 {
// // 	switch v := val.(type) {
// // 	case A:
// // 		*(*A)(data) = v
// // 		return 0
// // 	case B:
// // 		*(*B)(data) = v
// // 		return 1
// // 	}
// // 	fmt.Printf("%T\n", val)
// // 	panic("AAA")
// // }

// // // type Union2[T, A, B any] struct {
// // // 	inner Union[T]
// // // 	tag uint8
// // // }

// // // func (u *Union[T]) Put0(a int64) {
// // // 	u.tag = 0
// // // 	*(*int64)(unsafe.Pointer(&u.data)) = a
// // // }

// // // func (u *Union[T]) Put2(a int16) {
// // // 	u.tag = 1
// // // 	*(*int16)(unsafe.Pointer(&u.data)) = a
// // // }

// // // func (u *Union[T]) Get() any {
// // // 	switch u.tag {
// // // 	case 0:
// // // 		return *(*int64)(unsafe.Pointer(&u.data))
// // // 	case 1:
// // // 		return *(*int16)(unsafe.Pointer(&u.data))
// // // 	}
// // // 	panic("AAAA")
// // // }

// // // func (u *Union[T]) Put1(a int64) {
// // // 	u.tag = 0
// // // 	*(*int64)(unsafe.Pointer(&u.data)) = a
// // // }

// // // func (u *Union[T]) Put2(a int16) {
// // // 	u.tag = 1
// // // 	*(*int16)(unsafe.Pointer(&u.data)) = a
// // // }

// // // func (u *Union[T]) Get() any {
// // // 	switch u.tag {
// // // 	case 0:
// // // 		return *(*int64)(unsafe.Pointer(&u.data))
// // // 	case 1:
// // // 		return *(*int16)(unsafe.Pointer(&u.data))
// // // 	}
// // // 	panic("AAAA")
// // // }

// //--------------------------------------------------------------------------------

// // Note: You can't store anything with a pointer in here bc itll lose its reference I guess?
// type Union2[T any, A, B any] struct {
// 	data T
// 	tag uint8
// }

// func (u *Union2[T, A, B]) Put(val any) {
// 	u.tag = u.Encode(unsafe.Pointer(&u.data), val)
// }

// func (u *Union2[T, A, B]) Get() any {
// 	return u.Decode(u.tag, unsafe.Pointer(&u.data))
// }

// func (e Union2[T, A, B]) Decode(tag uint8, data unsafe.Pointer) any {
// 	switch tag {
// 	case 0:
// 		return *(*A)(data)
// 	case 1:
// 		return *(*B)(data)
// 	}
// 	panic("AAA")
// }

// func (e Union2[T, A, B]) Encode(data unsafe.Pointer, val any) uint8 {
// 	switch v := val.(type) {
// 	case A:
// 		*(*A)(data) = v
// 		return 0
// 	case B:
// 		*(*B)(data) = v
// 		return 1
// 	}
// 	fmt.Printf("%T\n", val)
// 	panic("AAA")
// }

// //--------------------------------------------------------------------------------

// type Union3[T any, A, B, C any] struct {
// 	data T
// 	tag uint8
// }

// func (u *Union3[T, A, B, C]) Put(val any) {
// 	u.tag = u.Encode(unsafe.Pointer(&u.data), val)
// }

// func (u *Union3[T, A, B, C]) Get() any {
// 	return u.Decode(u.tag, unsafe.Pointer(&u.data))
// }

// func (e Union3[T, A, B, C]) Decode(tag uint8, data unsafe.Pointer) any {
// 	switch tag {
// 	case 0:
// 		return *(*A)(data)
// 	case 1:
// 		return *(*B)(data)
// 	case 2:
// 		return *(*C)(data)
// 	}
// 	panic("AAA")
// }

// func (e Union3[T, A, B, C]) Encode(data unsafe.Pointer, val any) uint8 {
// 	switch v := val.(type) {
// 	case A:
// 		*(*A)(data) = v
// 		return 0
// 	case B:
// 		*(*B)(data) = v
// 		return 1
// 	case C:
// 		*(*C)(data) = v
// 		return 2
// 	}
// 	fmt.Printf("%T\n", val)
// 	panic("AAA")
// }
