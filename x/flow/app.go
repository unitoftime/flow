package flow

import (
	// "time"
	"github.com/unitoftime/ecs"
)

type Stage uint8
const (
	StageStartup Stage = iota
	StageFixedUpdate
	StageUpdate
)

type App struct {
	world *ecs.World
	scheduler *ecs.Scheduler

	startupSystems []ecs.System
}

func NewApp() *App {
	scheduler := ecs.NewScheduler()
	app := &App{
		world: ecs.NewWorld(),
		scheduler: scheduler,
		startupSystems: make([]ecs.System, 0),
	}

	AddResource(app, scheduler)

	return app
}

func (a *App) Run() {
	for _, sys := range a.startupSystems {
		sys.Run(0)
	}

	a.scheduler.Run()
}

func (a *App) AddSystems(stage Stage, systems ...func(*ecs.World) ecs.System) {
	for _, sys := range systems {
		system := sys(a.world)
		switch stage {
		case StageStartup:
			a.startupSystems = append(a.startupSystems, system)
		case StageFixedUpdate:
			a.scheduler.AppendPhysics(system)
		case StageUpdate:
			a.scheduler.AppendRender(system)
		}
	}
}

// func AddSystems1[A ecs.Initializer](a *App, stage Stage, lambda func(time.Duration, A)) {
// 	system := ecs.NewSystem1(a.world, lambda)
// 	switch stage {
// 	case StageStartup:
// 		a.startupSystems = append(a.startupSystems, system)
// 	case StageFixedUpdate:
// 		a.scheduler.AppendPhysics(system)
// 	case StageUpdate:
// 		a.scheduler.AppendRender(system)
// 	}
// }

func AddResource[T any](a *App, t *T) {
	ecs.PutResource(a.world, t)
}

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
