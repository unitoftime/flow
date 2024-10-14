package main

import (
	"fmt"

	"github.com/unitoftime/ecs"
	"github.com/unitoftime/glitch"

	"github.com/unitoftime/flow/hot"
	"github.com/unitoftime/flow/x/flow"
)

func main() {
	glitch.Run(run)
}

func run() {
	win, err := glitch.NewWindow(1920, 1080, "Gophermark", glitch.WindowConfig{
		Vsync: true,
	})
	if err != nil { panic(err) }

	app := flow.NewApp()
	flow.AddResource(app, win)
	// flow.AddResource(app, app.GetScheduler())


	// // --------------------------------------------------------------------------------
	// Yaegi

	// i := interp.New(interp.Options{
	// 	// GoPath: build.Default.GOPATH,
	// })
	// err = i.Use(stdlib.Symbols)
	// if err != nil { panic(err) }
	// err = i.Use(assets.Symbols)
	// if err != nil { panic(err) }

	// _, err = i.EvalPath("./assets/code/code.go")
	// if err != nil { panic(err) }

	// fixedAny, err := i.Eval("code.GetFixedSystem")
	// if err != nil { panic(err) }
	// fixedSystem := fixedAny.Interface().(func(any) ecs.System)

	// renderAny, err := i.Eval("code.GetSystem")
	// if err != nil { panic(err) }
	// renderSystem := renderAny.Interface().(func(any) ecs.System)

	// game := &code.Game{
	// 	Win: win,
	// 	Pass: pass,
	// 	Scheduler: app.GetScheduler(),
	// 	World: app.GetWorld(),
	// }

	// --------------------------------------------------------------------------------


	// assetServer := asset.NewServer(os.DirFS("./"))
	// asset.Register(assetServer, SpriteAssetLoader{})
	// flow.AddResource(app, assetServer)

	// // --------------------------------------------------------------------------------
	// // plugins := []string{"./assets/code/code.so"}
	// go func() {
	// 	for {
	// 		fixedSys := make([]ecs.System, 0)
	// 		renderSys := make([]ecs.System, 0)

	// 		entries, err := os.ReadDir("./assets/code/")
	// 		if err != nil { panic(err) }

	// 		plugins := make([]string, 0)
	// 		for _, e := range entries {
	// 			if strings.HasSuffix(e.Name(), ".so") {
	// 				fmt.Println(e.Name())
	// 				plugins = append(plugins, "./assets/code/" + e.Name())
	// 			}
	// 		}

	// 		for i := range plugins {
	// 			p, err := plugin.Open(plugins[i])
	// 			if err != nil { panic(err) }

	// 			fixedAny, err := p.Lookup("GetFixedSystems")
	// 			if err != nil { panic(err) }
	// 			fixedSystem := fixedAny.(func(*ecs.World) []ecs.System)

	// 			renderAny, err := p.Lookup("GetRenderSystems")
	// 			if err != nil { panic(err) }
	// 			renderSystem := renderAny.(func(*ecs.World) []ecs.System)

	// 			fixedSys = append(fixedSys, fixedSystem(app.GetWorld())...)
	// 			renderSys = append(renderSys, renderSystem(app.GetWorld())...)
	// 		}

	// 		fmt.Println("Reloading", len(fixedSys), len(renderSys))
	// 		scheduler := app.GetScheduler()
	// 		scheduler.SetPhysics(fixedSys...)
	// 		scheduler.SetRender(renderSys...)

	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()


	// app.Run()

	// // --------------------------------------------------------------------------------
	// // rm -f ../plugin/*.so && VAR=$RANDOM && echo $VAR && rm -rf ./build/* && mkdir ./build/tmp$VAR && cp reloader.go ./build/tmp$VAR && go build -buildmode=plugin -o ../plugin/tmp$VAR.so ./build/tmp$VAR
	// go func() {
	// 	currentPlugin := ""
	// 	nextPlugin := ""
	// 	for {
	// 		time.Sleep(1 * time.Second)

	// 		fixedSys := make([]ecs.System, 0)
	// 		renderSys := make([]ecs.System, 0)

	// 		entries, err := os.ReadDir("./assets/plugin/")
	// 		if err != nil { panic(err) }

	// 		// plugins := make([]string, 0)
	// 		for _, e := range entries {
	// 			if strings.HasSuffix(e.Name(), ".so") {
	// 				// plugins = append(plugins, "./assets/plugin/" + e.Name())
	// 				nextPlugin = "./assets/plugin/" + e.Name()
	// 				break
	// 			}
	// 		}

	// 		reload := nextPlugin != currentPlugin
	// 		if !reload {
	// 			fmt.Println(".")
	// 			continue
	// 		}

	// 		fmt.Println("Found New Plugin:", nextPlugin)
	// 		lastPlugin := currentPlugin
	// 		currentPlugin = nextPlugin

	// 		p, err := plugin.Open(currentPlugin)
	// 		if err != nil {
	// 			fmt.Println("Plugin already loaded(last, curr):", lastPlugin, currentPlugin)
	// 			panic(err)
	// 			continue
	// 		}
	// 		// if err != nil { panic(err) }

	// 		fixedAny, err := p.Lookup("GetFixedSystems")
	// 		if err != nil { panic(err) }
	// 		fixedSystem := fixedAny.(func(*ecs.World) []ecs.System)

	// 		renderAny, err := p.Lookup("GetRenderSystems")
	// 		if err != nil { panic(err) }
	// 		renderSystem := renderAny.(func(*ecs.World) []ecs.System)

	// 		fixedSys = append(fixedSys, fixedSystem(app.GetWorld())...)
	// 		renderSys = append(renderSys, renderSystem(app.GetWorld())...)

	// 		fmt.Println("Reloading", len(fixedSys), len(renderSys))
	// 		scheduler := app.GetScheduler()
	// 		scheduler.SetPhysics(fixedSys...)
	// 		scheduler.SetRender(renderSys...)
	// 	}
	// }()

	// --------------------------------------------------------------------------------
	// rm -f ../plugin/*.so && VAR=$RANDOM && echo $VAR && rm -rf ./build/* && mkdir ./build/tmp$VAR && cp reloader.go ./build/tmp$VAR && go build -buildmode=plugin -o ../plugin/tmp$VAR.so ./build/tmp$VAR

	// go build -ldflags "-pluginpath=plugin/hot-$(date +%s)" -buildmode=plugin -o hotload.so hotload.go

	// go build -ldflags "-pluginpath=plugin/hot-$(date +%s)" -buildmode=plugin -o ../plugin/plugin.so .

	plugin := hot.NewPlugin("./assets/plugin/")
	// plugin.Start()

	go func() {
		for {
			if !plugin.Check() { continue }
			// <-plugin.Refresh()

			fixedAny, err := plugin.Lookup("GetFixedSystems")
			if err != nil { panic(err) }
			fixedSystem := fixedAny.(func(*ecs.World) []ecs.System)

			renderAny, err := plugin.Lookup("GetRenderSystems")
			if err != nil { panic(err) }
			renderSystem := renderAny.(func(*ecs.World) []ecs.System)

			fixedSys := fixedSystem(app.GetWorld())
			renderSys := renderSystem(app.GetWorld())

			fmt.Println("Reloading", len(fixedSys), len(renderSys))
			scheduler := app.GetScheduler()
			scheduler.SetPhysics(fixedSys...)
			scheduler.SetRender(renderSys...)
		}
	}()

	app.Run()


	// --------------------------------------------------------------------------------
	// p, err := plugin.Open("./assets/code/code.so")
	// if err != nil { panic(err) }

	// fixedAny, err := p.Lookup("GetFixedSystem")
	// if err != nil { panic(err) }
	// fixedSystem := fixedAny.(func(*ecs.World) ecs.System)

	// renderAny, err := p.Lookup("GetSystem")
	// if err != nil { panic(err) }
	// renderSystem := renderAny.(func(*ecs.World) ecs.System)

	// --------------------------------------------------------------------------------
	// app.AddSystems(flow.StageStartup,
	// 	setupSystem,
	// )

	// app.AddSystems(flow.StageFixedUpdate,
	// 	// fixedUpdateSystem,
	// 	fixedSystem,
	// )

	// app.AddSystems(flow.StageUpdate,
	// 	// updateSystem,
	// 	renderSystem,
	// )

	// app.Run()
}

// func setupSystem(world *ecs.World) ecs.System {
// 	return ecs.NewSystem(func(dt time.Duration) {
// 	})
// }

// func updateSystem(world *ecs.World) ecs.System {
// 	// query := ecs.Query2[flow.Transform, phy2.Vel](world)

// 	return ecs.NewSystem(func(dt time.Duration) {
// 		fmt.Println("UpdateSystem")
// 		time.Sleep(10*time.Millisecond)
// 	})
// }

// func fixedUpdateSystem(world *ecs.World) ecs.System {
// 	// query := ecs.Query2[flow.Transform, phy2.Vel](world)

// 	return ecs.NewSystem(func(dt time.Duration) {
// 		fmt.Println("FixedUpdateSystem")
// 	})
// }

// type EcsPlugin struct {
	
// }

// func LoadPlugin(name string) {
// 	p, err := plugin.Open(name)
// 	if err != nil { panic(err) }

// 	fixedAny, err := p.Lookup("GetFixedSystem")
// 	if err != nil { panic(err) }
// 	fixedSystem := fixedAny.(func(*ecs.World) ecs.System)

// 	renderAny, err := p.Lookup("GetSystem")
// 	if err != nil { panic(err) }
// 	renderSystem := renderAny.(func(*ecs.World) ecs.System)
// }


// type SpriteAssetLoader struct {
// }
// func (l SpriteAssetLoader) Ext() []string {
// 	return []string{".png"}
// }
// func (l SpriteAssetLoader) Load(server *asset.Server, data []byte) (*glitch.Sprite, error) {
// 	smooth := true

// 	img, _, err := image.Decode(bytes.NewBuffer(data))
// 	if err != nil {
// 		return nil, err
// 	}

// 	texture := glitch.NewTexture(img, smooth)

// 	return glitch.NewSprite(texture, texture.Bounds()), nil
// }
// func (l SpriteAssetLoader) Store(server *asset.Server, sprite *glitch.Sprite) ([]byte, error) {
// 	return nil, errors.New("sprites do not support writeback")
// }

// type CodeAsset struct {

// }

// type CodeAssetLoader struct {
// }
// func (l CodeAssetLoader) Ext() []string {
// 	return []string{".so"}
// }
// func (l CodeAssetLoader) Load(server *asset.Server, data []byte) (*Code, error) {
// 	p, err := plugin.Open("./assets/code/code.so")
// 	if err != nil { panic(err) }

// 	fixedAny, err := p.Lookup("GetFixedSystem")
// 	if err != nil { panic(err) }
// 	fixedSystem := fixedAny.(func(*ecs.World) ecs.System)

// 	renderAny, err := p.Lookup("GetSystem")
// 	if err != nil { panic(err) }
// 	renderSystem := renderAny.(func(*ecs.World) ecs.System)

// 	return &Code{

// 	}
// }
// func (l CodeAssetLoader) Store(server *asset.Server, sprite *Code) ([]byte, error) {
// 	return nil, errors.New("code does not support writeback")
// }


