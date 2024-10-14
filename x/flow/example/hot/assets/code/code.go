package code

import (
	"time"

	"github.com/unitoftime/ecs"
	"github.com/unitoftime/flow/phy2"
	"github.com/unitoftime/flow/x/flow"
	"github.com/unitoftime/glitch"
)

type Pos struct {
	X, Y float64
}

func GetFixedSystems(world *ecs.World) []ecs.System {
	println("abcd")
	// fmt.Println("Fixed2")
	return []ecs.System{
		fixedUpdate(world),
	}
}

func fixedUpdate(world *ecs.World) ecs.System {
	query := ecs.Query2[flow.Transform, phy2.Vel](world)
	return ecs.NewSystem(func(dt time.Duration) {
		// fmt.Print("")
		// fmt.Println("Fixed2")

		query.MapId(func(id ecs.Id, transform *flow.Transform, vel *phy2.Vel) {
			transform.Position.X += vel.X * dt.Seconds()
			transform.Position.Y += vel.Y * dt.Seconds()
		})
	})
}

func GetRenderSystems(world *ecs.World) []ecs.System {
	return []ecs.System{
		renderUpdate(world),
	}
}

func renderUpdate(world *ecs.World) ecs.System {
	win := ecs.GetResource[glitch.Window](world)
	scheduler := ecs.GetResource[ecs.Scheduler](world)

	camera := glitch.NewCameraOrtho()

	return ecs.NewSystem(func(dt time.Duration) {
    // fmt.Println("Render")

		// TODO: Input?
		if win.JustPressed(glitch.KeyEscape) {
			scheduler.SetQuit(true)
		}

		// TODO: Draw stuff

		// camera setup
		camera.SetOrtho2D(win.Bounds())
		camera.SetView2D(0, 0, 1.0, 1.0)

		glitch.Clear(win, glitch.RGBA{1, 0.2, 0.3, 1.0})
		// glitch.Clear(win, glitch.RGBA{1.0, 0.0, 0.0, 1.0})
		glitch.SetCamera(camera)
		win.Update()
	})
}


// type Game struct {
// 	Win *glitch.Window
// 	Pass *glitch.RenderPass

// 	Scheduler *ecs.Scheduler
// 	World *ecs.World
// }
// func getGame(gameAny any) *Game {
// 	rval := reflect.ValueOf(gameAny)
// 	iGame := rval.Convert(reflect.TypeOf(&Game{}))
// 	game := iGame.Interface().(*Game)
// 	return game
// }


// type Pos struct {
// 	X, Y float64
// }

// func GetFixedSystem(gameAny any) ecs.System {
// 	game := getGame(gameAny)

// 	query := ecs.Query2[flow.Transform, phy2.Vel](game.World)
// 	return ecs.NewSystem(func(dt time.Duration) {
// 		query.MapId(func(id ecs.Id, transform *flow.Transform, vel *phy2.Vel) {
// 			transform.Position.X += vel.X * dt.Seconds()
// 			transform.Position.Y += vel.Y * dt.Seconds()
// 		})
// 	})
// }

// func GetSystem(gameAny any) ecs.System {
// 	game := getGame(gameAny)

// 	pass := game.Pass
// 	win := game.Win
// 	scheduler := game.Scheduler

// 	camera := glitch.NewCameraOrtho()

// 	return ecs.NewSystem(func(dt time.Duration) {
//     fmt.Println("Render")

// 		// TODO: Input?
// 		if win.JustPressed(glitch.KeyEscape) {
// 			scheduler.SetQuit(true)
// 		}

// 		pass.Clear()

// 		// TODO: Draw stuff

// 		// camera setup
// 		camera.SetOrtho2D(win.Bounds())
// 		camera.SetView2D(0, 0, 1.0, 1.0)

// 		glitch.Clear(win, glitch.RGBA{0.1, 0.2, 0.3, 1.0})
// 		pass.SetUniform("projection", camera.Projection)
// 		pass.SetUniform("view", camera.View)
// 		pass.Draw(win)
// 		win.Update()
// 	})
// }
