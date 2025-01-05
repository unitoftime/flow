package main

import (
	"image/color"
	"math"
	"time"

	"github.com/unitoftime/ecs"
	"github.com/unitoftime/flow"
	"github.com/unitoftime/flow/render"
	"github.com/unitoftime/flow/transform"
	"github.com/unitoftime/glitch"
)

// Ideas
// - ecs.PutResource -> world.AddResource(thing) but then have thing implement some interface

func main() {
	glitch.Run(run)
}

func run() {
	app := flow.NewApp()
	app.AddPlugin(render.DefaultPlugin{})
	app.AddPlugin(transform.DefaultPlugin{})

	app.AddSystems(ecs.StageStartup,
		ecs.NewSystem1(setup),
	)

	app.AddSystems(ecs.StageFixedUpdate,
		ecs.NewSystem1(rotate),
	)

	app.AddSystems(ecs.StageUpdate,
		ecs.NewSystem1(printStuff),
		ecs.NewSystem2(escapeExit),
	)

	app.Run()
}

func setup(dt time.Duration, commands *ecs.CommandQueue) {
	// TODO: I'd like to rewrite this to be internally managed, but for now you must manually call Execute()
	defer commands.Execute()

	texture := glitch.NewRGBATexture(128, 128, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, false)
	sprite := glitch.NewSprite(texture, texture.Bounds())

	commands.SpawnEmpty().
		Insert(render.Sprite{sprite}).
		Insert(render.Visibility{}).
		Insert(render.CalculatedVisibility{}).
		Insert(render.Target{}).
		Insert(transform.Local{transform.Default()}).
		Insert(transform.Global{transform.Default()})
}

func rotate(dt time.Duration, query *ecs.View1[transform.Local]) {
	query.MapId(func(id ecs.Id, lt *transform.Local) {
		lt.Rot += math.Pi * dt.Seconds()

		// angle := math.Pi * dt.Seconds()
		// lt.Rot.RotateZ(angle)
	})
}

func printStuff(dt time.Duration, query *ecs.View1[transform.Local]) {
	// query.MapId(func(id ecs.Id, lt *transform.Local) {
	// 	fmt.Println(id, lt)
	// })
}

func escapeExit(dt time.Duration, scheduler *ecs.Scheduler, query *ecs.View1[render.Window]) {
	query.MapId(func(_ ecs.Id, win *render.Window) {
		if win.JustPressed(glitch.KeyEscape) {
			win.Close()
			scheduler.SetQuit(true)
		}
	})
}
