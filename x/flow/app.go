package flow

import (
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
	return &App{
		world: ecs.NewWorld(),
		scheduler: ecs.NewScheduler(),
		startupSystems: make([]ecs.System, 0),
	}
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
