//go:build js || wasm

package storage

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"runtime"
	"runtime/pprof"
	"syscall/js"

	"github.com/mitchellh/mapstructure"
	"github.com/unitoftime/flow/browser"
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
	jsonDat, err := base64.StdEncoding.DecodeString(baseString)
	if err != nil {
		return nil, err
	}

	var ret T
	err = json.Unmarshal(jsonDat, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func GetItemWithDefault[T any](key string, def T) (*T, error) {
	defaultMap := make(map[string]any)
	err := mapstructure.Decode(def, &defaultMap)
	if err != nil {
		return nil, err
	}

	val := localStorage.Call("getItem", key)
	if val.IsNull() || val.IsUndefined() {
		return nil, nil
	}
	if val.Type() != js.TypeString {
		return nil, fmt.Errorf("failed to access data, must be a string")
	}

	baseString := val.String()
	jsonDat, err := base64.StdEncoding.DecodeString(baseString)
	if err != nil {
		return nil, err
	}

	decodedMap := make(map[string]any)
	err = json.Unmarshal(jsonDat, &decodedMap)
	if err != nil {
		return nil, err
	}

	maps.Copy(defaultMap, decodedMap)

	var ret T
	err = mapstructure.Decode(decodedMap, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func SetItem(key string, val any) error {
	valMap := make(map[string]any)
	err := mapstructure.Decode(val, &valMap)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(valMap)
	if err != nil {
		return err
	}

	baseString := base64.StdEncoding.EncodeToString(buf)

	localStorage.Call("setItem", key, baseString)
	return nil
}

// func RemoveItem[T any](key string) {
// }
// func ClearAll[T any]() {
// }

func GetQueryString(key string) ([]string, error) {
	href := js.Global().Get("location").Get("href").String()
	u, err := url.Parse(href)
	if err != nil {
		return nil, err
	}
	values, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}

	ret := values[key]
	return ret, nil
}

func WriteMemoryProfile(file string) error {
	buf := bytes.NewBuffer([]byte{})

	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(buf); err != nil {
		return err
	}
	dat := buf.Bytes()
	browser.Open("data:application/octet-stream;charset=utf-16le;base64,"+base64.StdEncoding.EncodeToString(dat), browser.OpenSameTab)
	return nil
}

func WriteCpuProfile(file string) (func(), error) {
	// Note: this code doesn't work. I guess cpu profiling isn't enabled in wasm
	// buf := bytes.NewBuffer([]byte{})

	// if err := pprof.StartCPUProfile(buf); err != nil {
	// 	return func(){}, err
	// }

	// finisher := func() {
	// 	pprof.StopCPUProfile()
	// 	dat := buf.Bytes()
	// 	browser.Open("data:application/octet-stream;charset=utf-16le;base64,"+base64.StdEncoding.EncodeToString(dat), false)
	// }

	// return finisher, nil
	return func() {}, nil
}
