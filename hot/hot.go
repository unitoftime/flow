//go:build !js

package hot

import (
	"errors"
	"fmt"
	"os"
	"plugin"
	"strings"
	"sync"
	"time"
)

//--------------------------------------------------------------------------------
// Notes
//--------------------------------------------------------------------------------
// 1. https://github.com/golang/go/issues/19004#issuecomment-288923294
// 2. https://github.com/edwingeng/hotswap/tree/main
// 3. https://stackoverflow.com/questions/70033200/plugin-already-loaded
// 4. https://segmentfault.com/a/1190000042104429/en#item-4
// 5. https://github.com/golang/go/issues/68349#issuecomment-2217718002
// 6. Helps checking to make sure the plugin and main binary match: `go version -m plugin.so` and `go version -m binary.bin`
// 7. Maybe helps (never used it): `go clean -cache`
// 8. Maybe helps (unsure the cases where it helps): `trimpath`
// 9. Didn't work: `go build -ldflags "-pluginpath=plugin/hot-$(date +%s)" -buildmode=plugin -o hotload.so hotload.go`
// 10. Might help with checking binary versions: `readelf -aW ../driver/plugin1.so | grep PluginFunc`

//--------------------------------------------------------------------------------
// Tips
//--------------------------------------------------------------------------------
// 1. If you reference to a package outside of the plugin directory: Changing the location of any global def (like a function definition) will break the plugin reloading bc it changes the version of the original package
// 2. If you split plugin into several small plugins: it makes it hard to have one plugin depend on another
// 3. If you make one big mega plugin, it would increase compile times (probably), but is easy to reference several things inside the plugin
//--------------------------------------------------------------------------------

var cache map[string]*Plugin

// rm -f ../plugin/*.so && VAR=$RANDOM && echo $VAR && rm -rf ./build/* && mkdir ./build/tmp$VAR && cp reloader.go ./build/tmp$VAR && go build -buildmode=plugin -o ../plugin/tmp$VAR.so ./build/tmp$VAR

type Plugin struct {
	path      string
	internal  *plugin.Plugin
	startOnce sync.Once
	// gen uint64
	refresh chan struct{}

	currentPlugin string
}

func NewPlugin(path string) *Plugin {
	if cache == nil {
		cache = make(map[string]*Plugin)
	}
	p, ok := cache[path]
	if ok {
		return p
	}

	newPlugin := Plugin{
		path:    path,
		refresh: make(chan struct{}),
	}
	cache[path] = &newPlugin
	return &newPlugin
}

// func (p Plugin) Generation() uint64 {
// 	return p.gen
// }

// func (p *Plugin) Refresh() chan struct{} {
// 	return p.refresh
// }

func (p *Plugin) Lookup(symName string) (any, error) {
	if p.internal == nil {
		return nil, errors.New("plugin not yet loaded")
	}
	val, err := p.internal.Lookup(symName)
	return val, err
}

// Check to see if there is a new plugin to load
// Returns true if there is a new one
func (p *Plugin) Check() bool {
	entries, err := os.ReadDir(p.path)
	if err != nil {
		panic(err)
	}

	nextPlugin := ""
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".so") {
			nextPlugin = p.path + e.Name()
			break
		}
	}

	if nextPlugin == "" {
		return false // Nothing new
	}

	samePlugin := nextPlugin == p.currentPlugin
	if samePlugin {
		return false
	}
	p.currentPlugin = nextPlugin
	fmt.Println("Found New Plugin:", nextPlugin)

	// var iPlugin *plugin.Plugin
	// func() {
	// 	defer func() {
	// 		if r := recover(); r != nil {
	// 			fmt.Println("RECOVERED:", r)
	// 		}
	// 	}()
	// 	time.Sleep(100 * time.Millisecond)
	// 	// Note: I have to sleep here to ensure that all of the glitch CGO calls have completed for the frame. 100ms is arbitrary, and is unecessary if you dont make CGO calls.
	// 	var err error
	// 	iPlugin, err = plugin.Open(p.currentPlugin)
	// 	if err != nil {
	// 		fmt.Println("Error Loading Plugin:", err)
	// 	}
	// }()

	// Note: I have to sleep here to ensure that all of the glitch CGO calls have completed for the frame. 100ms is arbitrary, and is unecessary if you dont make CGO calls.
	time.Sleep(100 * time.Millisecond)
	iPlugin, err := plugin.Open(p.currentPlugin)
	if err != nil {
		fmt.Println("Error Loading Plugin:", err)
		return false
	}

	if iPlugin == nil {
		fmt.Println("Error Loading Plugin")
		return false
	}

	fmt.Println("Successfully Loaded Plugin:", p.currentPlugin)
	p.internal = iPlugin
	return true
}

// Old idea:
// - Problem - can't synchronize with CGO execution which can cause SIGBUS: bus errors
// // Starts a watcher process in the background
// func (p *Plugin) Start() {
// 	p.startOnce.Do(func() {
// 		p.start()
// 	})
// }

// func (p *Plugin) start() {
// 	path := p.path

// 	sleepDur := 100 * time.Millisecond

// 	go func() {
// 		nextPlugin := ""
// 		for {
// 			time.Sleep(sleepDur)

// 			entries, err := os.ReadDir(path)
// 			if err != nil { panic(err) }

// 			for _, e := range entries {
// 				if strings.HasSuffix(e.Name(), ".so") {
// 					nextPlugin = path + e.Name()
// 					break
// 				}
// 			}

// 			reload := nextPlugin != p.currentPlugin
// 			if !reload {
// 				// fmt.Println(".")
// 				continue
// 			}

// 			fmt.Println("Found New Plugin:", nextPlugin)
// 			// lastPlugin := currentPlugin
// 			p.currentPlugin = nextPlugin

// 			var iPlugin *plugin.Plugin
// 			func() {
// 				defer func() {
// 					if r := recover(); r != nil {
// 						fmt.Println("RECOVERED:", r)
// 					}
// 				}()
// 				var err error
// 				iPlugin, err = plugin.Open(p.currentPlugin)
// 				if err != nil {
// 					panic(err)
// 					// fmt.Println("Error Loading Plugin:", currentPlugin, err)
// 					// // fmt.Println("Plugin already loaded(last, curr):", lastPlugin, currentPlugin)
// 					// continue
// 				}
// 			}()
// 			if iPlugin == nil {
// 				fmt.Println("Error Loading Plugin")
// 				continue
// 			}

// 			fmt.Println("Successfully Loaded Plugin:", p.currentPlugin)
// 			p.internal = iPlugin
// 			// p.gen++
// 			p.refresh <- struct{}{}
// 		}
// 	}()
// }
