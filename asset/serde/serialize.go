package serde

import (
	"bytes"
	"encoding/gob"
)

func Register(value any) {
	gob.Register(value)
}

func RegisterName(name string, value any) {
	gob.RegisterName(name, value)
}

func Marshal[T any](t T) ([]byte, error) {
	var dat bytes.Buffer
	enc := gob.NewEncoder(&dat)
	err := enc.Encode(t)
	if err != nil {
		return nil, err
	}
	return dat.Bytes(), nil
}

func Unmarshal[T any](dat []byte) (T, error) {
	dec := gob.NewDecoder(bytes.NewReader(dat))

	var t T
	err := dec.Decode(&t)
	if err != nil {
		return t, err
	}
	return t, err
}
