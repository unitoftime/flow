package main

import (
	"time"

	"github.com/unitoftime/flow/example/hot/plugin"
	"github.com/unitoftime/flow/hot"
)

// Steps:
// 1. run rebuild.sh to build the initial plugin *.so file
// 2. Run main.go
// 3. Change something in the HelloWorld function in ./plugin/plugin.go
// 4. run rebuild.sh to rebuild the plugin *.so file. the main() function below will detect, reload the symbol and run it once

// Note: You need to run rebuild.sh to build the plugin *.so file

func main() {
	// You can directly import the plugin package to use it immediately
	plugin.HelloWorld()

	// This will search the provided directory for .so files and try to load them
	p := hot.NewPlugin("./plugin/build/lib/")

	for {
		time.Sleep(1 * time.Second)
		if !p.Check() {
			continue
		} // When this becomes true, it means a new plugin is loaded

		// With our new plugin, we can lookup our symbol `HelloWorld`
		sym, err := p.Lookup("HelloWorld")
		if err != nil {
			panic(err)
		}
		hello := sym.(func())

		// Then we can call our Looked up symbol
		hello()
	}
}
