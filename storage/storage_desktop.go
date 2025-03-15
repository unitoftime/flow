//go:build !js

package storage

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
)

var storageRoot = ""

// Sets the root directory for local file storage
func SetStorageRoot(root string) error {
	err := os.MkdirAll(root, 0700)
	if err != nil {
		return err
	}

	storageRoot = root

	return nil
}

// Gets a copy of the item out of storage and returns a pointer to it. Else returns nil
// If there is no item we will return nil
// If there is an error getting or deserializing the item we will return (nil, error)
func GetItem[T any](key string) (*T, error) {
	key = filepath.Join(storageRoot, key)
	dat, err := os.ReadFile(key)
	if err != nil {
		return nil, err
	}

	return deserialize[T](dat)
}

func GetItemWithDefault[T any](key string, def T) (*T, error) {
	key = filepath.Join(storageRoot, key)

	dat, err := os.ReadFile(key)
	if err != nil {
		return nil, err
	}

	return deserializeWithDefault(dat, def)
}

func SetItem(key string, val any) error {
	key = filepath.Join(storageRoot, key)

	buf, err := serialize(val)
	if err != nil {
		return err
	}

	// Perm: Read/Write owner, nothing for anyone else
	err = os.WriteFile(key, buf, 0600)
	return err
}

func GetQueryString(key string) ([]string, error) {
	return []string{}, nil
}

func WriteMemoryProfile(file string) error {
	file = filepath.Join(storageRoot, file)

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close() // error handling omitted for example
	runtime.GC()    // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		return err
	}
	return nil
}

func WriteCpuProfile(file string) (func(), error) {
	file = filepath.Join(storageRoot, file)

	f, err := os.Create(file)
	if err != nil {
		return func() {}, err
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		return func() {}, err
	}

	finisher := func() {
		defer f.Close() // error handling omitted for example
		pprof.StopCPUProfile()
	}
	return finisher, nil
}
