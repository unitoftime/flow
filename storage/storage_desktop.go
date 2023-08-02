//go:build !js
package storage

// Note: Desktop version of this doesn't work yet. Entire API will likely change as I develop my storage needs

func GetItem[T any](key string) (*T, error) {
	return nil, nil
}

func SetItem(key string, val any) error {
	return nil
}
