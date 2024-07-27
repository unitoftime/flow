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

// rm -f ../plugin/*.so && VAR=$RANDOM && echo $VAR && rm -rf ./build/* && mkdir ./build/tmp$VAR && cp reloader.go ./build/tmp$VAR && go build -buildmode=plugin -o ../plugin/tmp$VAR.so ./build/tmp$VAR

type Plugin struct {
	path string
	internal *plugin.Plugin
	startOnce sync.Once
	// gen uint64
	refresh chan struct{}
}

func NewPlugin(path string) Plugin {
	return Plugin{
		path: path,
		refresh: make(chan struct{}),
	}
}

// func (p Plugin) Generation() uint64 {
// 	return p.gen
// }

func (p Plugin) Refresh() chan struct{} {
	return p.refresh
}

func (p Plugin) Lookup(symName string) (any, error) {
	if p.internal == nil {
		return nil, errors.New("plugin not yet loaded")
	}
	val, err := p.internal.Lookup(symName)
	return val, err
}

func (p *Plugin) Start() {
	p.startOnce.Do(func() {
		p.start()
	})
}

func (p *Plugin) start() {
	path := p.path

	sleepDur := 1 * time.Second

	go func() {
		currentPlugin := ""
		nextPlugin := ""
		for {
			time.Sleep(sleepDur)

			entries, err := os.ReadDir(path)
			if err != nil { panic(err) }

			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".so") {
					nextPlugin = path + e.Name()
					break
				}
			}

			reload := nextPlugin != currentPlugin
			if !reload {
				fmt.Println(".")
				continue
			}

			fmt.Println("Found New Plugin:", nextPlugin)
			// lastPlugin := currentPlugin
			currentPlugin = nextPlugin

			iPlugin, err := plugin.Open(currentPlugin)
			if err != nil {
				fmt.Println("Error Loading Plugin:", currentPlugin, err)
				// fmt.Println("Plugin already loaded(last, curr):", lastPlugin, currentPlugin)
				continue
			}

			fmt.Println("Successfully Loaded Plugin:", currentPlugin)
			p.internal = iPlugin
			// p.gen++
			p.refresh <- struct{}{}
		}
	}()
}
