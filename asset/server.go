package asset

import (
	"fmt"
	"strings"
	"sync/atomic"
)

// TODO: if I wanted this to be more "ecs-like" I would make a resource per asset type then use some kind of integer handle (ie `type Handle[T] uint32` or something). Then I use that handle to index into the asset type resource (ie `assets.Get(handle)` and `assets := ecs.GetResource[T](world)`)

// TODO: Finalizers on handles to deallocate assets that are no longer used?
type Handle[T any] struct {
	ptr atomic.Pointer[T]
	name string
}
func newHandle[T any](name string) *Handle[T] {
	return &Handle[T]{
		name: name,
	}
}

func (h *Handle[T]) Set(val *T) {
	h.ptr.Store(val)
}
func (h *Handle[T]) Get() *T {
	return h.ptr.Load()
}

type assetHandler interface {}

type Loader[T any] interface {
	Ext() []string
	Load(*Server, []byte) (*T, error)
}


type Server struct {
	load *Load

	extToLoader map[string]any // Map file extension strings to the loader that loads them

	nameToHandle map[string]assetHandler // Map the full filepath name to the asset handle
}
func NewServer(load *Load) *Server {
	return &Server{
		load: load,
		extToLoader: make(map[string]any),

		nameToHandle: make(map[string]assetHandler),
	}
}

func Register[T any](s *Server, loader Loader[T]) {
	extensions := loader.Ext()
	for _, ext := range extensions {
		_, exists := s.extToLoader[ext]
		if exists {
			panic(fmt.Sprintf("duplicate loader registration: %s", ext))
		}

		s.extToLoader[ext] = loader
	}
}

func LoadAsset[T any](server *Server, name string) *Handle[T] {
	// Check if already loaded
	anyHandle, ok := server.nameToHandle[name]
	if ok {
		handle := anyHandle.(*Handle[T])
		return handle
	}

	// Find a loader for it
	_, ext, found := strings.Cut(name, ".")
	if !found {
		ext = name
	}

	anyLoader, ok := server.extToLoader[ext]
	if !ok {
		panic(fmt.Sprintf("could not find loader for extension: %s", ext))
	}
	loader, ok := anyLoader.(Loader[T])
	if !ok {
		panic(fmt.Sprintf("wrong type for registered loader on extension: %s", ext))
	}

	handle := newHandle[T](name)
	server.nameToHandle[name] = handle

	go func() {
		data, err := server.load.Data(name)
		if err != nil {
			// TODO: Store error
			return
		}

		val, err := loader.Load(server, data)
		if err != nil {
			// TODO: Store error
			return
		}

		handle.Set(val)
	}()

	// Success
	return handle
}

// func LoadAsset[T any](server *Server, name string) *Handle[T] {
// 	handle := newHandle[T](name)
// }

// Note: This is more like how bevy works
// type UntypedHandle uint64

// type Handle[T any] struct {
// 	UntypedHandle
// }

// type Asset struct {
// 	Error error
// 	Value any
// }

// type Loader interface {
// 	Ext() []string
// 	Load(data []byte) (any, error)
// }


// type Server struct {
// 	load *Load

// 	extToLoader map[string]Loader

// 	nameToHandle map[string]UntypedHandle
// 	assets []Asset
// }
// func NewServer(load *Load) *Server {
// 	return &Server{
// 		load: load,
// 		extToLoader: make(map[string]Loader),

// 		nameToHandle: make(map[string]UntypedHandle),
// 		assets: make([]Asset, 0),
// 	}
// }


// func (s *Server) Register(loader Loader) {
// 	extensions := loader.Ext()
// 	for _, ext := range extensions {
// 		_, exists := s.extToLoader[ext]
// 		if exists {
// 			panic(fmt.Sprintf("duplicate loader registration: %s", ext))
// 		}

// 		s.extToLoader[ext] = loader
// 	}
// }

// func (s *Server) addAsset(name string) (*Asset, UntypedHandle) {
// 	s.assets = append(s.assets, Asset{})
// 	handle := UntypedHandle(len(s.assets) - 1)
// 	s.nameToHandle[name] = handle

// 	return &s.assets[handle], handle
// }

// func (s *Server) LoadUntyped(name string) UntypedHandle {
// 	// Check if already loaded
// 	handle, ok := s.nameToHandle[name]
// 	if ok {
// 		return handle
// 	}

// 	// Find a loader for it
// 	_, ext, found := strings.Cut(name, ".")
// 	if !found {
// 		ext = name
// 	}

// 	loader, ok := s.extToLoader["."+ext]
// 	if !ok {
// 		panic(fmt.Sprintf("could not find loader for extension: %s", ext))
// 	}

// 	asset, handle := s.addAsset(name)

// 	// TODO: load dynamically (maybe chan?)
// 	data, err := s.load.Data(name)
// 	if err != nil {
// 		asset.Error = err
// 		return handle
// 	}

// 	loadedVal, err := loader.Load(data)
// 	if err != nil {
// 		asset.Error = err
// 		return handle
// 	}

// 	// Success
// 	asset.Value = loadedVal
// 	return handle
// }

// func (s *Server) Get(handle UntypedHandle) (any, error) {
// 	asset := s.assets[handle]
// 	return asset.Value, asset.Error
// }

// func LoadAsset[T any](server *Server, name string) Handle[T] {
// 	uHandle := server.LoadUntyped(name)
// 	return Handle[T]{uHandle}
// }

// func GetAsset[T any](server *Server, handle Handle[T]) (T, error) {
// 	asset, err := server.Get(handle.UntypedHandle)
// 	if err != nil {
// 		var t T
// 		return t, err
// 	}
// 	asset.
// }
