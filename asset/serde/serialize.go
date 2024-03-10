package serde

import (
	"reflect"

	"github.com/raszia/gotiny"
)

func Register[T any](value T) {
	// gob.Register(value)
	gotiny.Register(value)
}
func RegisterName[T any](name string, value T) {
	// gob.RegisterName(name, value)
	gotiny.RegisterName(name, reflect.TypeOf(value))
}

func Marshal[T any](t T) ([]byte, error) {
	data := gotiny.Marshal(&t)
	return data, nil
}

func Unmarshal[T any](dat []byte) (T, error) {
	var t T
	gotiny.Unmarshal(dat, &t)
	return t, nil
}

// func Marshal[T any](t T) ([]byte, error) {
// 	var dat bytes.Buffer
// 	// enc := gob.NewEncoder(&dat)
// 	enc := stablegob.NewEncoder(&dat)
// 	err := enc.Encode(t)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return dat.Bytes(), nil
// }

// func Unmarshal[T any](dat []byte) (T, error) {
// 	dec := stablegob.NewDecoder(bytes.NewReader(dat))
// 	// dec := gob.NewDecoder(bytes.NewReader(dat))

// 	var t T
// 	err := dec.Decode(&t)
// 	if err != nil {
// 		return t, err
// 	}
// 	return t, err
// }



// func Register[T any](value T) {
// 	register(value)
// }

// func RegisterName[T any](name string, value T) {
// 	registerName(name, value)
// }

// func Marshal[T any](t T) ([]byte, error) {
// 	m := Encode(t)
// 	return json.Marshal(m)
// }

// func Unmarshal[T any](dat []byte) (T, error) {
// 	var t T
// 	m := make(map[string]any)
// 	err := json.Unmarshal(dat, &m)
// 	if err != nil {
// 		return t, err
// 	}

// 	Decode(&t, m)
// 	return t, err
// }
