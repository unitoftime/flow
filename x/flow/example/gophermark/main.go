package main

import (
	"fmt"
	"os"
	"time"
	"embed"
	"math/rand"

	"github.com/unitoftime/ecs"

	"github.com/unitoftime/flow/x/flow"
	"github.com/unitoftime/flow/phy2"
	"github.com/unitoftime/flow/asset"
	"github.com/unitoftime/flow/interp"


	"github.com/unitoftime/glitch"
	"github.com/unitoftime/glitch/shaders"
)

// components := []Component{
// 	phy2.Pos{1, 2},
// 	phy2.Vel{3, 4},
// 	Handle<Sprite>{"gopher.png"},
// }

// components := []Component{
// 	Transform{...},
// 	phy2.Vel{3, 4},
// 	Handle<Sprite>{"gopher.png"},
// }

// type ComponentBuilder interface {
// 	Get() []Component
// }
// components := []ComponentBuilder{
// 	Transform{...} -> (Pos, Scale, Rot)
// 	SpriteData{"gopher.png"} -> glitch.Sprite
// }



//go:embed assets/*
var EmbeddedFilesystem embed.FS

func main() {
	glitch.Run(run)
}

func run() {
	win, err := glitch.NewWindow(1920, 1080, "Gophermark", glitch.WindowConfig{
		Vsync: false,
	})
	if err != nil { panic(err) }
	shader, err := glitch.NewShader(shaders.PixelArtShader)
	if err != nil { panic(err) }
	pass := glitch.NewRenderPass(shader)

	app := flow.NewApp()
	flow.AddResource(app, win)
	flow.AddResource(app, pass)

	assetServer := asset.NewServer(os.DirFS("./"))
	asset.Register(assetServer, SpriteAssetLoader{})
	asset.Register(assetServer, PrefabAssetLoader{})
	flow.AddResource(app, assetServer)

	atlas, err := glitch.DefaultAtlas()
	if err != nil { panic(err) }
	flow.AddResource(app, atlas)

	app.AddSystems(flow.StageStartup,
		setupSystem,
	)

	app.AddSystems(flow.StageFixedUpdate,
		copyTransform,
		movementSystem,
		collisionSystem,
	)

	app.AddSystems(flow.StageUpdate,
		inputSystem,
		renderSystem,
	)

	app.Run()
}

// type Transform struct {
// 	Matrix glitch.Mat4
// }
// type Transform struct {
// 	Position Vec3
// 	Scale Vec3
// 	Rotation Vec4 // Quaternion
// }
// func (t Transform) RenderMatrix() glitch.Mat4 {
// 	mat := glitch.Mat4Ident
// 	return mat.
// 		Scale(t.Scale.X, t.Scale.Y, t.Scale.Z).
// 		Rotate()
// }


type Sprite struct {
	Sprite *asset.Handle[glitch.Sprite]
}
func (s *Sprite) Bounds() glitch.Rect {
	sprite := s.Sprite.Get()
	if sprite == nil {
		return glitch.Rect{}
	}
	return sprite.Bounds()
}
func (s *Sprite) Draw(pass *glitch.RenderPass, mat glitch.Mat4) {
	sprite := s.Sprite.Get()
	if sprite == nil {
		return
	}

	sprite.Draw(pass, mat)
}

type Text struct {
	Text *glitch.Text
}
func (t *Text) Draw(pass *glitch.RenderPass, mat glitch.Mat4) {

}


// TODO: Command buffer injection
func setupSystem(world *ecs.World) ecs.System {
	assetServer := ecs.GetResource[asset.Server](world)

	gopherCount := 20000
	vScale := 200.0

	return ecs.NewSystem(func(dt time.Duration) {
		// text := atlas.Text("hello world", 1)
		// id := world.NewId()
		// ecs.Write(world, id,
		// 	ecs.C(flow.NewTransform()),
		// 	ecs.C(Text{
		// 		text,
		// 	}),
		// )

		gopherSprite := asset.Load[glitch.Sprite](assetServer, "assets/gopher.png")

		for i := 0; i < gopherCount; i++ {
			id := world.NewId()
			transform := flow.NewTransform()
			transform.Position.X = 1920/2
			transform.Position.Y = 1080/2
			ecs.Write(world, id,
				ecs.C(transform),
				ecs.C(flow.LastTransform(transform)),
				// ecs.C(phy2.Pos{1920/2, 1080/2}),
				ecs.C(phy2.Vel{
					float64(2*vScale * (rand.Float64()-0.5)),
					float64(2*vScale * (rand.Float64()-0.5)),
				}),
				ecs.C(Sprite{
					gopherSprite,
				}),
			)
		}

		prefab := asset.Load[PrefabAsset](assetServer, "assets/test.prefab.toml")
		fmt.Println(prefab.Get())
	})
}

func movementSystem(world *ecs.World) ecs.System {
	query := ecs.Query2[flow.Transform, phy2.Vel](world)

	// lastFrameTime := time.Now()
	return ecs.NewSystem(func(dt time.Duration) {
		query.MapId(func(id ecs.Id, transform *flow.Transform, vel *phy2.Vel) {
			transform.Position.X += vel.X * dt.Seconds()
			transform.Position.Y += vel.Y * dt.Seconds()
		})

		// fmt.Println("Phy:", time.Since(lastFrameTime))
		// lastFrameTime = time.Now()
	})
}

func copyTransform(world *ecs.World) ecs.System {
	query := ecs.Query2[flow.Transform, flow.LastTransform](world)
	return ecs.NewSystem(func(dt time.Duration) {
		query.MapId(func(id ecs.Id, transform *flow.Transform, lastTransform *flow.LastTransform) {
			*lastTransform = flow.LastTransform(*transform)
		})
	})
}

func collisionSystem(world *ecs.World) ecs.System {
	query := ecs.Query3[flow.Transform, phy2.Vel, Sprite](world)
	win := ecs.GetResource[glitch.Window](world)
	bounds := win.Bounds()

	return ecs.NewSystem(func(dt time.Duration) {
		query.MapId(func(id ecs.Id, transform *flow.Transform, vel *phy2.Vel, sprite *Sprite) {
			w := sprite.Bounds().W() / 2
			h := sprite.Bounds().H() / 2
			pos := transform.Position
			if pos.X <= 0 || (pos.X + w) >= bounds.W() {
				vel.X = -vel.X
			}
			if pos.Y <= 0 || (pos.Y + h) >= bounds.H() {
				vel.Y = -vel.Y
			}
		})
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
	query := ecs.Query3[flow.LastTransform, flow.Transform, Sprite](world)
	pass := ecs.GetResource[glitch.RenderPass](world)
	win := ecs.GetResource[glitch.Window](world)
	scheduler := ecs.GetResource[ecs.Scheduler](world)

	camera := glitch.NewCameraOrtho()
	lastFrameTime := time.Now()

	return ecs.NewSystem(func(dt time.Duration) {
		interpVal := scheduler.GetRenderInterp()
		fmt.Println("Interp:", interpVal)
		pass.Clear()
		mat := glitch.Mat4Ident
		query.MapId(func(id ecs.Id, lastTransform *flow.LastTransform, transform *flow.Transform, sprite *Sprite) {
			mat = glitch.Mat4Ident
			// mat.Scale(0.25, 0.25, 1.0).Translate(pos.X, pos.Y, 0)
			xPos := interp.Linear.Float64(lastTransform.Position.X, transform.Position.X, interpVal)
			yPos := interp.Linear.Float64(lastTransform.Position.Y, transform.Position.Y, interpVal)
			// fmt.Println(lastTransform.Position.X, transform.Position.X)
			mat.Scale(0.25, 0.25, 1.0)
			mat.Translate(xPos, yPos, 0)


			// mat.Translate(transform.Position.X, transform.Position.Y, 0)
			sprite.Draw(pass, mat)
		})

		// camera setup
		camera.SetOrtho2D(win.Bounds())
		camera.SetView2D(0, 0, 1.0, 1.0)

		glitch.Clear(win, glitch.RGBA{0.1, 0.2, 0.3, 1.0})
		pass.SetUniform("projection", camera.Projection)
		pass.SetUniform("view", camera.View)
		pass.Draw(win)
		win.Update()

		fmt.Println("Ren", time.Since(lastFrameTime))
		// fmt.Println("RenderInterp: ", scheduler.GetRenderInterp())
		lastFrameTime = time.Now()
	})
}

// var MovementSystem = flow.System1(movementSystemFunc)
// func movementSystemFunc(dt time.Duration, query *ecs.View2[phy2.Pos, phy2.Vel]) {
// 	query.MapId(func(id ecs.Id, pos *phy2.Pos, vel *phy2.Vel) {
// 		fmt.Println("map move")
// 		pos.X += vel.X * dt.Seconds()
// 		pos.Y += vel.Y * dt.Seconds()
// 	})
// }


// // var movementSystem = flow.System1(
// // 	func(dt time.Duration, query *ecs.View2[phy2.Pos, phy2.Vel]) {
// // 		query.MapId(func(id ecs.Id, pos *phy2.Pos, vel *phy2.Vel) {
// // 			fmt.Println("map move")
// // 			pos.X += vel.X * dt.Seconds()
// // 			pos.Y += vel.Y * dt.Seconds()
// // 		})
// // 	}
// // )

// type MovementSystem struct{
// 	query *ecs.View2[phy2.Pos, phy2.Vel]
// }
// // func (s *MovementSystem) Setup(world *ecs.World) {
// // 	s.query := ecs.Query2[phy2.Pos, phy2.Vel](world)
// // }

// func (s *MovementSystem) Run(dt time.Duration) {
// 	s.query.MapId(func(id ecs.Id, pos *phy2.Pos, vel *phy2.Vel) {
// 		fmt.Println("map move")
// 		pos.X += vel.X * dt.Seconds()
// 		pos.Y += vel.Y * dt.Seconds()
// 	})
// }

// var movementSystem = flow.System1(movementSystemFunc)
// func movementSystemFunc(dt time.Duration, query *ecs.View2[phy2.Pos, phy2.Vel]) {
// 	query.MapId(func(id ecs.Id, pos *phy2.Pos, vel *phy2.Vel) {
// 		fmt.Println("map move")
// 		pos.X += vel.X * dt.Seconds()
// 		pos.Y += vel.Y * dt.Seconds()
// 	})
// }

// func movementSystem(world *ecs.World) ecs.System {
// 	return ecs.NewSystem1(world,
// 		func(dt time.Duration, query *ecs.View2[phy2.Pos, phy2.Vel]) {
// 			query.MapId(func(id ecs.Id, pos *phy2.Pos, vel *phy2.Vel) {
// 				fmt.Println("map move")
// 				pos.X += vel.X * dt.Seconds()
// 				pos.Y += vel.Y * dt.Seconds()
// 			})
// 		},
// 	)
// }

// func movementSystem(world *ecs.World) ecs.System {
// 	query := ecs.Query2[phy2.Pos, phy2.Vel](world)
// 	return ecs.NewSystem(func(dt time.Duration) {
// 		movementSystemFunc(dt, query)
// 	})


// 	return ecs.NewSystem(func(dt time.Duration) {
// 		query.MapId(func(id ecs.Id, pos *phy2.Pos, vel *phy2.Vel) {
// 			fmt.Println("map move")
// 			pos.X += vel.X * dt.Seconds()
// 			pos.Y += vel.Y * dt.Seconds()
// 		})
// 	})
// }

// var TransformBundle = ecs.NewBundle2[phy2.Pos, phy2.Vel]()
// var TimerBundle = ecs.NewBundle2[CustomTimer, Activator]()

// func setupSystem(world *ecs.World) ecs.System {
// 	command := ecs.NewCommand(world)

// 	// posWriter := ecs.NewWriter[phy2.Pos]()
// 	// velWriter := ecs.NewWriter[phy2.Vel]()

// 	return ecs.NewSystem(func(dt time.Duration) {
// 		fmt.Println("setupSystem")

// 		// bundle := TransformBundle{
// 		// 	phy2.Pos{},
// 		// 	phy2.Vel{},
// 		// }

// 		// bundle := ecs.DynamicBundle(
// 		// 	TransformBundle.With(
// 		// 		phy2.Pos{},
// 		// 		phy2.Vel{},
// 		// 	),
// 		// )

// 		command.Spawn(
// 			TransformBundle.With(
// 				phy2.Pos{},
// 				phy2.Vel{},
// 			),
// 			TimerBundle.With(
// 				CustomTimer{timer.New(1 * time.Second, 1.0)},
// 				Activator{},
// 			),
// 		)

// 		// ent := NewEnt()
// 		// posWriter.
// 		// 	Write(ent, phy2.Pos{})


// 		// id := world.NewId()
// 		// ecs.WriteCmd(command, id, phy2.Pos{})
// 		// ecs.WriteCmd(command, id, phy2.Vel{})
// 		// ecs.WriteCmd(command, id, timer.New(1 * time.Second, 1.0))
// 		// ecs.WriteCmd(command, id, CustomTimer{timer.New(1 * time.Second, 1.0)})
// 		// ecs.WriteCmd(command, id, Activator{})

// 		// id := world.NewId()
// 		// ecs.Write(world, id,
// 		// 	ecs.C(phy2.Pos{}),
// 		// 	ecs.C(phy2.Vel{}),
// 		// 	ecs.C(timer.New(1 * time.Second, 1.0)),
// 		// 	ecs.C(CustomTimer{timer.New(1 * time.Second, 1.0)}),
// 		// 	ecs.C(Activator{}),
// 		// )

// 		command.Execute()
// 	})
// }

// type PTimer[T any] interface {
// 	*T
// 	Update(time.Duration) bool
// }

// type PActivator[A any] interface {
// 	*A
// 	Activate()
// }

// type CustomTimer struct {
// 	timer.Timer
// }
// type Activator struct {
// }
// func (a *Activator) Activate(/*commandbuffer*/) {
// 	fmt.Println("DOTHING")
// 	// do thing
// }

// func timerSystem[T any, A any, PT PTimer[T], PA PActivator[A]](world *ecs.World) ecs.System {
// 	query := ecs.Query2[T, A](world)
// 	return ecs.NewSystem(func(dt time.Duration) {
// 		query.MapId(func(id ecs.Id, t *T, activator *A) {
// 			fmt.Println("map timer", t, dt.Seconds())
// 			if PT(t).Update(dt) {
// 				PA(activator).Activate(/*world.commandbuffer()*/)
// 			}
// 		})
// 	})
// }

// // func movement2(dt time.Duration, query *flow.Query2[*phy2.Pos, *phy2.Vel, ecs.Filter1[ecs.With[phy2.Rot]]]) {
// // // func movement2(dt time.Duration, query *ecs.Query2[phy2.Pos, phy2.Vel, ecs.Filter1[ecs.With[phy2.Rot]]]) {
// // // func movement2(dt time.Duration, query *ecs.Query2[phy2.Pos, phy2.Vel, ecs.Filter[ecs.All]]) {
// // // func movement2(dt time.Duration, query *ecs.Query2[phy2.Pos, phy2.Vel, ecs.With[phy2.Rot]]) {
// // // func movement(dt time.Duration, query *ecs.View2[phy2.Pos, phy2.Vel]) {
// // func movement(dt time.Duration, query *ecs.View2[phy2.Pos, phy2.Vel]) {
// // 	query.Map2(func(pos *phy2.Pos, vel *phy2.Vel) {
// // 		pos.X += vel.X * dt
// // 		pos.Y += vel.Y * dt
// // 	})
// // }

// // func movement(world *ecs.World) ecs.System {
// // 	var query ecs.Query2[phy2.Pos, phy2.Vel, ecs.Filter1[ecs.With[phy2.Rot]]]
// // 	query.Init(world)
// // 	return ecs.NewSystem(func(dt time.Duration) {
// // 		movement2(dt, query)
// // 	})
// // }

// // func movement2(dt time.Duration, query *ecs.View2[phy2.Pos, phy2.Vel, ecs.Filter1[ecs.With[phy2.Rot]]]) {
// // 	query.Map2(func(pos *phy2.Pos, vel *phy2.Vel) {
// // 		pos.X += vel.X * dt
// // 		pos.Y += vel.Y * dt
// // 	})
// // }
