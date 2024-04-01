//go:build !js

package storage

import (
	"os"
	"runtime"
	"runtime/pprof"
)

// Note: Desktop version of this doesn't work yet. Entire API will likely change as I develop my storage needs

func GetItem[T any](key string) (*T, error) {
	return nil, nil
}

func SetItem(key string, val any) error {
	return nil
}

func GetQueryString(key string) ([]string, error) {
	return []string{}, nil
}

func WriteMemoryProfile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close() // error handling omitted for example
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		return err
	}
	return nil
}

func WriteCpuProfile(file string) (func(), error) {
	f, err := os.Create(file)
	if err != nil {
		return func(){}, err
	}
	defer f.Close() // error handling omitted for example
	if err := pprof.StartCPUProfile(f); err != nil {
		return func(){}, err
	}
	return pprof.StopCPUProfile, nil
}
