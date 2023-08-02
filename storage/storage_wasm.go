//go:build js || wasm
package storage

import (
	"fmt"
	"bytes"
	"syscall/js"
	"encoding/gob"
	"encoding/base64"
)

// TODO Maybe: https://developer.mozilla.org/en-US/docs/Web/API/IndexedDB_API

// TODO: Should these go into init? they could both theoretically panic?
var window = js.Global().Get("window")
var localStorage = window.Get("localStorage")

// Gets a copy of the item out of storage and returns a pointer to it. Else returns nil
// If there is no item we will return nil
// If there is an error getting or deserializing the item we will return (nil, error)
func GetItem[T any](key string) (*T, error) {
	val := localStorage.Call("getItem", key)
	if val.IsNull() || val.IsUndefined() {
		return nil, nil
	}
	if val.Type() != js.TypeString {
		return nil, fmt.Errorf("failed to access data, must be a string")
	}

	baseString := val.String()
	gobDat, err := base64.StdEncoding.DecodeString(baseString)
	if err != nil {
		return nil, err
	}

	var ret T
	buf := bytes.NewBuffer(gobDat)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func SetItem(key string, val any) error {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(val)
	if err != nil {
		return err
	}

	baseString := base64.StdEncoding.EncodeToString(buf.Bytes())

	localStorage.Call("setItem", key, baseString)
	return nil
}

// func RemoveItem[T any](key string) {
// }
// func ClearAll[T any]() {
// }
