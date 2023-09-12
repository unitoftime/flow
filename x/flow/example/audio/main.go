package main

import (
	"os"
	"time"

	"github.com/unitoftime/ecs"

	"github.com/unitoftime/flow/x/flow"
	"github.com/unitoftime/flow/x/flow/audio"
	"github.com/unitoftime/flow/asset"


	"github.com/unitoftime/glitch"
)

func main() {
	glitch.Run(run)
}

func run() {
	win, err := glitch.NewWindow(1920, 1080, "Audio", glitch.WindowConfig{
		Vsync: true,
	})
	if err != nil { panic(err) }

	app := flow.NewApp()
	flow.AddResource(app, win)

	audio.Initialize()

	assetServer := asset.NewServer(os.DirFS("./"))
	asset.Register(assetServer, audio.AssetLoader{})
	flow.AddResource(app, assetServer)

	app.AddSystems(flow.StageStartup,
		setupSystem,
	)

	app.AddSystems(flow.StageUpdate,
		renderSystem,
		inputSystem,
	)

	app.Run()
}

func setupSystem(world *ecs.World) ecs.System {
	assetServer := ecs.GetResource[asset.Server](world)

	return ecs.NewSystem(func(dt time.Duration) {
		audioFile := asset.Load[audio.Audio](assetServer,
			"https://upload.wikimedia.org/wikipedia/commons/c/c8/Example.ogg")
		audio.Play(audioFile.Get())
	})
}


func inputSystem(world *ecs.World) ecs.System {
	win := ecs.GetResource[glitch.Window](world)
	scheduler := ecs.GetResource[ecs.Scheduler](world)

	return ecs.NewSystem(func(dt time.Duration) {
		if win.JustPressed(glitch.KeyEscape) {
			scheduler.SetQuit(true)
		}
	})
}

func renderSystem(world *ecs.World) ecs.System {
	win := ecs.GetResource[glitch.Window](world)

	return ecs.NewSystem(func(dt time.Duration) {
		glitch.Clear(win, glitch.Black)

		win.Update()
	})
}
