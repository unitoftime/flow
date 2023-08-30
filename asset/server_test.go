package asset

import (
	"testing"

	"os"
	"time"
	"encoding/json"
)

type MyAsset struct {
	Health int
}
type CustomAssetLoader struct {
}
func (l CustomAssetLoader) Ext() []string {
	return []string{"1.json"}
}
func (l CustomAssetLoader) Load(server *Server, data []byte) (*MyAsset, error) {
	var myAsset MyAsset
	err := json.Unmarshal(data, &myAsset)
	return &myAsset, err
}

type MyAsset2 struct {
	HealthAsset *Handle[MyAsset]
}
type CustomAssetLoader2 struct {
}
func (l CustomAssetLoader2) Ext() []string {
	return []string{"2.json"}
}

func (l CustomAssetLoader2) Load(server *Server, data []byte) (*MyAsset2, error) {
	assetMap := make(map[string]string)
	err := json.Unmarshal(data, &assetMap)
	if err != nil {
		return nil, err
	}

	healthAssetName := assetMap["Health"]
	myAsset := MyAsset2{
		HealthAsset: LoadAsset[MyAsset](server, healthAssetName),
	}
	return &myAsset, err
}

func TestAssetServerBasic(t *testing.T) {
	server := NewServer(NewLoad(os.DirFS("./test-data")))
	Register(server, CustomAssetLoader{})

	handle1 := LoadAsset[MyAsset](server, "test.1.json")
	handle2 := LoadAsset[MyAsset](server, "test2.1.json")

	time.Sleep(1 * time.Second)

	t.Log(handle1)
	t.Log(handle2)

	asset1 := handle1.Get()
	asset2 := handle2.Get()

	t.Log(asset1)
	t.Log(asset2)
}


func TestAssetServerNested(t *testing.T) {
	server := NewServer(NewLoad(os.DirFS("./test-data")))
	Register(server, CustomAssetLoader{})
	Register(server, CustomAssetLoader2{})

	handle1 := LoadAsset[MyAsset2](server, "test.2.json")

	time.Sleep(1 * time.Second)

	t.Log(handle1)

	asset1 := handle1.Get()

	t.Log(asset1)
	t.Log(asset1.HealthAsset.Get())
}

// func TestAssetServerTyped(t *testing.T) {
// 	server := NewServer(NewLoad(os.DirFS("./test-data")))
// 	Register(server, CustomAssetLoader{})

// 	handle1 := LoadAsset[MyAsset](server, "test.json")
// 	handle2 := LoadAsset[MyAsset](server, "test2.json")

// 	t.Log(handle1)
// 	t.Log(handle2)

// 	asset1, err1 := GetAsset(server, handle1)
// 	asset2, err2 := GetAsset(server, handle2)

// 	t.Log(asset1, err1)
// 	t.Log(asset2, err2)
// }


// type MyAsset struct {
// 	Health int
// }
// type CustomAssetLoader struct {
// }
// func (l CustomAssetLoader) Ext() []string {
// 	return []string{".json"}
// }
// func (l CustomAssetLoader) Load(data []byte) (*MyAsset, error) {
// 	var myAsset MyAsset
// 	err := json.Unmarshal(data, &myAsset)
// 	return &myAsset, err
// }

// func TestAssetServerUntyped(t *testing.T) {
// 	server := NewServer(NewLoad(os.DirFS("./test-data")))
// 	server.Register(CustomAssetLoader{})

// 	handle1 := server.LoadUntyped("test.json")
// 	handle2 := server.LoadUntyped("test2.json")

// 	t.Log(handle1)
// 	t.Log(handle2)

// 	asset1, err1 := server.Get(handle1)
// 	asset2, err2 := server.Get(handle2)

// 	t.Log(asset1, err1)
// 	t.Log(asset2, err2)

// 	asset
// }

// func TestAssetServerTyped(t *testing.T) {
// 	server := NewServer(NewLoad(os.DirFS("./test-data")))
// 	server.Register(CustomAssetLoader{})

// 	handle1 := LoadAsset[MyAsset](server, "test.json")
// 	handle2 := LoadAsset[MyAsset](server, "test2.json")

// 	t.Log(handle1)
// 	t.Log(handle2)

// 	asset1, err1 := GetAsset(server, handle1)
// 	asset2, err2 := GetAsset(server, handle2)

// 	t.Log(asset1, err1)
// 	t.Log(asset2, err2)
// }
