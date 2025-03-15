package storage

import (
	"encoding/json"
	"maps"

	"github.com/mitchellh/mapstructure"
)

func serialize(val any) ([]byte, error) {
	valMap := make(map[string]any)
	err := mapstructure.Decode(val, &valMap)
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(valMap)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func deserialize[T any](dat []byte) (*T, error) {
	var ret T
	err := json.Unmarshal(dat, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func deserializeWithDefault[T any](dat []byte, def T) (*T, error) {
	defaultMap := make(map[string]any)
	err := mapstructure.Decode(def, &defaultMap)
	if err != nil {
		return nil, err
	}

	decodedMap := make(map[string]any)
	err = json.Unmarshal(dat, &decodedMap)
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
