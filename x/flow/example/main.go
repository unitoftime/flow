package main

import (
	"fmt"
	"time"

	"github.com/unitoftime/ecs"

	"github.com/unitoftime/flow/x/flow"
	"github.com/unitoftime/flow/phy2"
	"github.com/unitoftime/flow/timer"

	// "github.com/unitoftime/glitch"
)

func main() {
	app := flow.NewApp()
	// app.AddSystems(flow.StageFixedUpdate, ecs.NewSystem1(world, movement))
	app.AddSystems(flow.StageStartup, setupSystem)

	app.AddSystems(flow.StageFixedUpdate,
		movementSystem,
		(timerSystem[CustomTimer, Activator]), // Weird bug?: without parenthesis it doesnt work
	)
	app.Run()
}

func movementSystem(world *ecs.World) ecs.System {
	query := ecs.Query2[phy2.Pos, phy2.Vel](world)
	return ecs.NewSystem(func(dt time.Duration) {
		query.MapId(func(id ecs.Id, pos *phy2.Pos, vel *phy2.Vel) {
			fmt.Println("map move")
			pos.X += vel.X * dt.Seconds()
			pos.Y += vel.Y * dt.Seconds()
		})
	})
}

func setupSystem(world *ecs.World) ecs.System {
	return ecs.NewSystem(func(dt time.Duration) {
		fmt.Println("setupSystem")
		id := world.NewId()
		ecs.Write(world, id,
			ecs.C(phy2.Pos{}),
			ecs.C(phy2.Vel{}),
			ecs.C(timer.New(1 * time.Second, 1.0)),
			ecs.C(CustomTimer{timer.New(1 * time.Second, 1.0)}),
			ecs.C(Activator{}),
		)
	})
}

type PTimer[T any] interface {
	*T
	Update(time.Duration) bool
}

type PActivator[A any] interface {
	*A
	Activate()
}

type CustomTimer struct {
	timer.Timer
}
type Activator struct {
}
func (a *Activator) Activate(/*commandbuffer*/) {
	fmt.Println("DOTHING")
	// do thing
}

func timerSystem[T any, A any, PT PTimer[T], PA PActivator[A]](world *ecs.World) ecs.System {
	query := ecs.Query2[T, A](world)
	return ecs.NewSystem(func(dt time.Duration) {
		query.MapId(func(id ecs.Id, t *T, activator *A) {
			fmt.Println("map timer", t, dt.Seconds())
			if PT(t).Update(dt) {
				PA(activator).Activate(/*world.commandbuffer()*/)
			}
		})
	})
}

// func movement2(dt time.Duration, query *flow.Query2[*phy2.Pos, *phy2.Vel, ecs.Filter1[ecs.With[phy2.Rot]]]) {
// // func movement2(dt time.Duration, query *ecs.Query2[phy2.Pos, phy2.Vel, ecs.Filter1[ecs.With[phy2.Rot]]]) {
// // func movement2(dt time.Duration, query *ecs.Query2[phy2.Pos, phy2.Vel, ecs.Filter[ecs.All]]) {
// // func movement2(dt time.Duration, query *ecs.Query2[phy2.Pos, phy2.Vel, ecs.With[phy2.Rot]]) {
// // func movement(dt time.Duration, query *ecs.View2[phy2.Pos, phy2.Vel]) {
// func movement(dt time.Duration, query *ecs.View2[phy2.Pos, phy2.Vel]) {
// 	query.Map2(func(pos *phy2.Pos, vel *phy2.Vel) {
// 		pos.X += vel.X * dt
// 		pos.Y += vel.Y * dt
// 	})
// }

// func movement(world *ecs.World) ecs.System {
// 	var query ecs.Query2[phy2.Pos, phy2.Vel, ecs.Filter1[ecs.With[phy2.Rot]]]
// 	query.Init(world)
// 	return ecs.NewSystem(func(dt time.Duration) {
// 		movement2(dt, query)
// 	})
// }

// func movement2(dt time.Duration, query *ecs.View2[phy2.Pos, phy2.Vel, ecs.Filter1[ecs.With[phy2.Rot]]]) {
// 	query.Map2(func(pos *phy2.Pos, vel *phy2.Vel) {
// 		pos.X += vel.X * dt
// 		pos.Y += vel.Y * dt
// 	})
// }
