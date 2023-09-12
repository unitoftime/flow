package main

import (
	"bytes"

	"encoding/gob"
	"github.com/BurntSushi/toml"

	"dario.cat/mergo"
	"github.com/jinzhu/copier"

	"github.com/unitoftime/flow/asset"
	"github.com/unitoftime/ecs"
)

type EntConfig struct {
	Components []ecs.Component
}

// type EntConfig[T any] struct {
// 	Data T
// }

// func (e EntConfig) With(e2 EntConfig) EntConfig {
// 	err := mergo.Merge(&e, e2, mergo.WithOverride)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return e
// }

type PrefabAsset struct {
	Id string
	Base *asset.Handle[PrefabAsset]
	Override EntConfig
}

func (p *PrefabAsset) Get() EntConfig {
	var ret EntConfig
	if p == nil { return ret }

	if p.Base != nil {
		prefab := p.Base.Get()
		ret = prefab.Get()
	}

	err := mergo.Merge(&ret, p.Override, mergo.WithOverride)
	if err != nil {
		panic(err) // From what I understand, the only errors are for passing in bad values
	}

	var final EntConfig
	err = copier.Copy(&final, &ret)
	if err != nil {
		panic(err)
	}

	return final
}

type PrefabAssetLoader struct {
}
func (l PrefabAssetLoader) Ext() []string {
	return []string{".prefab.toml", ".prefab.gob"}
}
func (l PrefabAssetLoader) Load(server *asset.Server, data []byte) (*PrefabAsset, error) {
	var ass PrefabAsset
	err := unmarshalPrefab(&ass, data)

	if ass.Base != nil {
		handle := asset.Load[PrefabAsset](server, ass.Base.Name)
		ass.Base = handle
	}

	return &ass, err
}

func (l PrefabAssetLoader) Store(server *asset.Server, ass *PrefabAsset) ([]byte, error) {
	data, err := marshalPrefab(ass)
	return data, err
}

func unmarshalPrefab(v any, data []byte) error {
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(v)
	if err == nil {
		return nil
	}

	return toml.Unmarshal(data, v)
}

func marshalPrefab(v any) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(v)
	if err == nil {
		return buf.Bytes(), err
	}

	buf = new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	err = encoder.Encode(v)
	if err == nil {
		return buf.Bytes(), nil
	}
	return nil, err
}
