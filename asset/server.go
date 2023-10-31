package asset

import (
	"fmt"
	"os"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// TODO: if I wanted this to be more "ecs-like" I would make a resource per asset type then use some kind of integer handle (ie `type Handle[T] uint32` or something). Then I use that handle to index into the asset type resource (ie `assets.Get(handle)` and `assets := ecs.GetResource[T](world)`)

// TODO: Finalizers on handles to deallocate assets that are no longer used?
type Handle[T any] struct {
	ptr atomic.Pointer[T]
	Name string
	err error
	doneChan chan struct{}
	done atomic.Bool
	modTime time.Time
}
func newHandle[T any](name string) *Handle[T] {
	return &Handle[T]{
		Name: name,
		doneChan: make(chan struct{}),
	}
}

func (h *Handle[T]) Set(val *T) {
	h.err = nil
	h.done.Store(true)
	h.ptr.Store(val)
}
func (h *Handle[T]) Get() (*T, error) {
	h.Wait()
	return h.ptr.Load(), h.err
}
func (h *Handle[T]) Err() error {
	return h.err
}

// Returns true if the asset is done loading
// At this point either an error or the asset will be available
func (h *Handle[T]) Done() bool {
	return h.done.Load()
}

// Blocks until the handle, or an error is set
func (h *Handle[T]) Wait() {
	<- h.doneChan
}

type assetHandler interface {}

type Loader[T any] interface {
	Ext() []string
	Load(*Server, []byte) (*T, error)
	Store(*Server, *T) ([]byte, error)
}

type Server struct {
	fsPath string
	filesystem fs.FS // TODO: Maybe use: https://pkg.go.dev/github.com/ungerik/go-fs

	mu sync.Mutex
	extToLoader map[string]any // Map file extension strings to the loader that loads them
	nameToHandle map[string]assetHandler // Map the full filepath name to the asset handle
}
func NewServerFromPath(fsPath string) *Server {
	filesystem := os.DirFS(fsPath)
	return &Server{
		fsPath: fsPath,
		filesystem: filesystem,
		extToLoader: make(map[string]any),

		nameToHandle: make(map[string]assetHandler),
	}
}
func NewServer(filesystem fs.FS) *Server {
	return &Server{
		filesystem: filesystem,
		extToLoader: make(map[string]any),

		nameToHandle: make(map[string]assetHandler),
	}
}


func getScheme(path string) string {
	u, err := url.Parse(path)
	if err != nil {
		return ""
	}
	return u.Scheme
}

func (s *Server) getModTime(fpath string) (time.Time, error) {
	// TODO: Wont work for networked files
	file, err := s.filesystem.Open(fpath)
	if err != nil { return time.Time{}, err }

	info, err := file.Stat()
	if err != nil { return time.Time{}, err }

	return info.ModTime(), nil
}

func (s *Server) readRaw(fpath string) ([]byte, time.Time, error) {
	scheme := getScheme(fpath)

	var rc io.ReadCloser
	var err error
	var modTime time.Time
	if scheme == "https" || scheme == "http" {
		rc, err = s.getHttp(fpath)
	} else {
		rc, modTime, err = s.getFile(fpath)
	}
	if err != nil {
		return nil, modTime, err
	}
	defer rc.Close()

	dat, err := ioutil.ReadAll(rc)
	return dat, modTime, err
}

func (s *Server) getFile(fpath string) (io.ReadCloser, time.Time, error) {
	file, err := s.filesystem.Open(fpath)
	if err != nil { return nil, time.Time{}, err }
	info, err := file.Stat()
	if err != nil { return nil, time.Time{}, err }

	return file, info.ModTime(), nil
}

func (s *Server) getHttp(fpath string) (io.ReadCloser, error) {
	resp, err := http.Get(fpath)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (s *Server) writeRaw(fpath string, dat []byte) error {
	// file, err := s.filesystem.Open(fpath)
	// if err != nil {
	// 	return nil, err
	// }
	// defer file.Close()

	fullFilepath := filepath.Join(s.fsPath, fpath)

	// Build entire filepath
	err := os.MkdirAll(filepath.Dir(fullFilepath), 0750)
	if err != nil {
		return err
	}

	// TODO: verify file is writable.
	return os.WriteFile(fullFilepath, dat, 0755)
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

// TODO: Extension filters?
// TODO: Should this return a single handle that gives us access to subhandles in the directory?
// Loads a directory that contains the same asset type. Returns a slice filled with all asset handles. Does not search recursively
func LoadDir[T any](server *Server, fpath string) []*Handle[T] {
	fpath = filepath.Clean(fpath)

	dirEntries, err := fs.ReadDir(server.filesystem, fpath)
	if err != nil {
		return nil // TODO!!! : You're just snuffing an error here, which obviously isn't good
	}


	ret := make([]*Handle[T], 0, len(dirEntries))
	for _, e := range dirEntries {
		if e.IsDir() { continue } // TODO: Recursive?

		handle := Load[T](server, filepath.Join(fpath, e.Name()))
		ret = append(ret, handle)
	}

	return ret
}

// Gets the handle, returns true if the handle has already started loading
func getHandle[T any](server *Server, name string) (*Handle[T], bool) {
	server.mu.Lock()
	defer server.mu.Unlock()

	// Check if already loaded
	anyHandle, ok := server.nameToHandle[name]
	if ok {
		handle := anyHandle.(*Handle[T])
		return handle, true
	}

	handle := newHandle[T](name)
	server.nameToHandle[name] = handle
	return handle, false
}

// Loads a single file
func Load[T any](server *Server, name string) *Handle[T] {
	handle, loaded := getHandle[T](server, name)
	if loaded {
		return handle
	}

	// Find a loader for it
	ext := getExtension(name)

	anyLoader, ok := server.extToLoader[ext]
	if !ok {
		panic(fmt.Sprintf("could not find loader for extension: %s", ext))
	}
	loader, ok := anyLoader.(Loader[T])
	if !ok {
		panic(fmt.Sprintf("wrong type for registered loader on extension: %s", ext))
	}

	go func() {
		// TODO: Recover?
		defer func() {
			handle.done.Store(true)
			close(handle.doneChan)
		}()

		data, modTime, err := server.readRaw(name)
		if err != nil {
			handle.err = err
			return
		}

		handle.modTime = modTime // TODO: Data race here if reload is called simultaneously with load

		val, err := loader.Load(server, data)
		if err != nil {
			handle.err = err
			return
		}

		handle.Set(val)
	}()

	// Success
	return handle
}

// Loads a single file
func Reload[T any](server *Server, handle *Handle[T]) {
	name := handle.Name

	// Find a loader for it
	ext := getExtension(name)

	anyLoader, ok := server.extToLoader[ext]
	if !ok {
		panic(fmt.Sprintf("could not find loader for extension: %s", ext))
	}
	loader, ok := anyLoader.(Loader[T])
	if !ok {
		panic(fmt.Sprintf("wrong type for registered loader on extension: %s", ext))
	}

	go func() {
		// TODO: Recover?

		modTime, err := server.getModTime(name)
		if err != nil {
			handle.err = err
			return
		}
		if handle.modTime == modTime {
			// Same file, don't reload
			return
		}

		data, modTime, err := server.readRaw(name)
		if err != nil {
			handle.err = err
			return
		}
		handle.modTime = modTime

		val, err := loader.Load(server, data)
		if err != nil {
			handle.err = err
			return
		}

		handle.Set(val)
	}()
}

// Writes the asset handle back to the file
func Store[T any](server *Server, handle *Handle[T]) error {
	name := handle.Name
	// Find a loader for it
	ext := getExtension(name)

	anyLoader, ok := server.extToLoader[ext]
	if !ok {
		panic(fmt.Sprintf("could not find loader for extension: %s", ext))
	}
	loader, ok := anyLoader.(Loader[T])
	if !ok {
		panic(fmt.Sprintf("wrong type for registered loader on extension: %s", ext))
	}

	val, _ := handle.Get()
	// Note: We skip error checking here b/c we really dont care if there was an error loading. All we want to do is write data to a file. Hence val shouldn't be nil
	if val == nil {
		return fmt.Errorf("handle data can't be nil when storing")
	}
	// if err != nil {
	// 	return err
	// }

	dat, err := loader.Store(server, val)
	if err != nil {
		return err
	}

	return server.writeRaw(handle.Name, dat)
}

func getExtension(name string) string {
	idx := -1
	for i := len(name) - 1; i >= 0 && name[i] != '/'; i-- {
		if name[i] == '.' {
			idx = i
		}
	}
	if idx > 0 {
		return name[idx:]
	}
	return ""

	// Note: Does't properly cut slashes
	// _, ext, found := strings.Cut(name, ".")
	// if !found {
	// 	ext = name
	// }

	// Note: Only returns the very final extension
	// ext := filepath.Ext(name)
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
