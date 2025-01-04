package flow

import (
	"github.com/unitoftime/ecs"
)

// type Stage uint8

// const (
// 	StageStartup Stage = iota
// 	StagePreFixedUpdate
// 	StageFixedUpdate
// 	StagePostFixedUpdate
// 	StageUpdate
// )

type App struct {
	world     *ecs.World
	scheduler *ecs.Scheduler
}

func NewApp() *App {
	world := ecs.NewWorld()
	scheduler := ecs.NewScheduler(world)
	ecs.PutResource(world, scheduler)
	// scheduler.SetFixedTimeStep(4 * time.Millisecond)

	app := &App{
		world:          world,
		scheduler:      scheduler,
	}

	return app
}

func (a *App) Run() {
	// fmt.Printf("%+v", a.startupSystems)
	// fmt.Printf("%+v", a.scheduler)

	a.scheduler.Run()
}

type Plugin interface {
	Initialize(world *ecs.World)
}
func (a *App) AddPlugin(plugin Plugin) {
	plugin.Initialize(a.world)
}

func (a *App) AddSystems(stage ecs.Stage, systems ...ecs.SystemBuilder) {
	a.scheduler.AddSystems(stage, systems...)
	// for _, sys := range systems {
	// 	system := sys.Build(a.world)
	// 	switch stage {
	// 	case StageStartup:
	// 		a.startupSystems = append(a.startupSystems, system)
	// 	case StageFixedUpdate:
	// 		a.scheduler.AppendPhysics(system)
	// 	case StageUpdate:
	// 		a.scheduler.AppendRender(system)
	// 	}
	// }
}

// func (a *App) AddSystems2(stage Stage, systems ...ecs.System) {
// 	for _, system := range systems {
// 		switch stage {
// 		case StageStartup:
// 			a.startupSystems = append(a.startupSystems, system)
// 		case StageFixedUpdate:
// 			a.scheduler.AppendPhysics(system)
// 		case StageUpdate:
// 			a.scheduler.AppendRender(system)
// 		}
// 	}
// }

// func (a *App) AddSystems(stage Stage, systems ...func(*ecs.World) ecs.System) {
// 	for _, sys := range systems {
// 		system := sys(a.world)
// 		switch stage {
// 		case StageStartup:
// 			a.startupSystems = append(a.startupSystems, system)
// 		case StageFixedUpdate:
// 			a.scheduler.AppendPhysics(system)
// 		case StageUpdate:
// 			a.scheduler.AppendRender(system)
// 		}
// 	}
// }

// func (a *App) SetSystems(stage Stage, systems ...func(*ecs.World) ecs.System) {
// 	for _, sys := range systems {
// 		system := sys(a.world)
// 		switch stage {
// 		case StageStartup:
// 			a.startupSystems = append(a.startupSystems, system)
// 		case StageFixedUpdate:
// 			a.scheduler.AppendPhysics(system)
// 		case StageUpdate:
// 			a.scheduler.AppendRender(system)
// 		}
// 	}
// }

// // func AddSystems1[A ecs.Initializer](a *App, stage Stage, lambda func(time.Duration, A)) {
// // 	system := ecs.NewSystem1(a.world, lambda)
// // 	switch stage {
// // 	case StageStartup:
// // 		a.startupSystems = append(a.startupSystems, system)
// // 	case StageFixedUpdate:
// // 		a.scheduler.AppendPhysics(system)
// // 	case StageUpdate:
// // 		a.scheduler.AppendRender(system)
// // 	}
// // }

// func AddResource[T any](a *App, t *T) {
// 	ecs.PutResource(a.world, t)
// }

// func System1[A ecs.Initializer](sysFunc func(time.Duration, A)) func(*ecs.World) ecs.System {
// 	return func(world *ecs.World) ecs.System {
// 		ecs.NewSystem1(world, sysFunc)
// 	}
// }

// func System[A, B ecs.Initializer](sysFunc func(time.Duration, A, B)) func(*ecs.World) ecs.System {
// 	return func(world *ecs.World) ecs.System {
// 		ecs.NewSystem2(world, sysFunc)
// 	}
// }
