package main

import (
	"github.com/unitoftime/ecs"
	"github.com/unitoftime/flow/x/flow/example/hot/assets/code"
)

func GetFixedSystems(world *ecs.World) []ecs.System {
	return code.GetFixedSystems(world)
}

func GetRenderSystems(world *ecs.World) []ecs.System {
	return code.GetRenderSystems(world)
}
