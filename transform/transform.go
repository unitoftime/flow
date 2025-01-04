package transform

import (
	"time"

	"github.com/unitoftime/ecs"
	"github.com/unitoftime/flow/ds"
	"github.com/unitoftime/flow/glm"
)

//go:generate go run ../../cod/cmd/cod

// Todo List
// 1. Migrate to 3D transforms
// 2. Revisit optimizations for heirarchy resolutions

type DefaultPlugin struct {
}

func (p DefaultPlugin) Initialize(world *ecs.World) {
	scheduler := ecs.GetResource[ecs.Scheduler](world)

	// TODO: This should be added to a better stage
	scheduler.AppendPhysics(ResolveHeirarchySystem(world))
}

func Default() Transform {
	return Transform{
		Scale: glm.Vec2{1, 1},
	}
}

func FromPos(pos glm.Vec2) Transform {
	return Transform{
		Pos:   pos,
		Rot:   0,
		Scale: glm.Vec2{1, 1},
	}
}

//cod:component
type Local struct {
	Transform
}

type Transform struct {
	Pos   glm.Vec2
	Rot   float64
	Scale glm.Vec2
}

func (t Transform) Mat4() glm.Mat4 {
	mat := glm.Mat4Ident
	mat.
		Scale(t.Scale.X, t.Scale.Y, 1).
		Rotate(t.Rot, glm.Vec3{0, 0, 1}).
		Translate(t.Pos.X, t.Pos.Y, 0)
	return mat
}

// Returns the transform, but with the vector space moved to the parent transform
func (t Transform) Globalize(parent Global) Global {
	childGlobal := t

	// OG
	// childGlobal.Pos = childLocal.Pos.Add(parentGlobal.Pos)

	// try 2 - manually calculate
	parentMat := parent.Mat4()
	dstPos := parentMat.Apply(t.Pos.Vec3())
	childGlobal.Pos = dstPos.Vec2()
	childGlobal.Rot = parent.Rot + t.Rot
	childGlobal.Scale = t.Scale.Mult(parent.Scale)

	// // try 3 - multiply mats and pull out values
	// Notes: https://math.stackexchange.com/questions/237369/given-this-transformation-matrix-how-do-i-decompose-it-into-translation-rotati
	// parentMat := parentGlobal.Mat4()
	// childMat := childLocal.Mat4()
	// parentMat.Mul(&childMat)

	// childGlobal.Pos = parentMat.GetTranslation().Vec2()
	// childGlobal.Rot = parentMat.GetRotation().Z
	// childGlobal.Scale = parentMat.GetScale().Vec2()
	return Global{childGlobal}
}

//cod:component
type Global struct {
	Transform
}

// func (t GlobalTransform) Mat4() glm.Mat4 {
// 	mat := glm.Mat4Ident
// 	mat.
// 		Scale(t.Scale.X, t.Scale.Y, 1).
// 		Rotate(t.Rot, glm.Vec3{0, 0, 1}).
// 		Translate(t.Pos.X, t.Pos.Y, 0)
// 	return mat
// }

//cod:component
type Parent struct {
	Id ecs.Id
}

// 3. You could reorganize so that parent's know their transform children, and recursively calculate each child's GlobalTransform
//
//cod:component
type Children struct {
	// List []ecs.Id
	MiniSlice ds.MiniSlice[[8]ecs.Id, ecs.Id]
}

// func (c *Children) Add(id ecs.Id) {
// 	c.List = append(c.List, id)
// }
// func (c *Children) Remove(id ecs.Id) {
// 	idx := slices.Index(c.List, id)
// 	if idx < 0 { return } // Skip: Doesnt exist

// 	c.List[idx] = c.List[len(c.List) - 1]
// 	c.List = c.List[:len(c.List) - 1]
// }
// func (c *Children) Clear() {
// 	if c.List == nil { return }
// 	c.List = c.List[:0]
// }

func (c *Children) Add(id ecs.Id) {
	c.MiniSlice.Append(id)
}
func (c *Children) Remove(id ecs.Id) {
	idx := c.MiniSlice.Find(id)
	if idx < 0 {
		return
	} // Skip: Doesnt exist
	c.MiniSlice.Delete(idx)
}
func (c *Children) Clear() {
	c.MiniSlice.Clear()
}

// Recursively goes through the transform heirarchy and calculates entity GlobalTransform based on their Parent's GlobalTransform and their local Transform.
func ResolveHeirarchySystem(world *ecs.World) ecs.System {
	queryTopLevel := ecs.Query3[Children, Local, Global](world,
		ecs.Without(Parent{}),
		ecs.Optional(Children{}, Local{}),
		// Note: Some entities are "top level" and dont have a local transform, so we will just use their global transform as is. I should move away from this model
		// ^^^^^
		// TODO: Fix hack where Local{} is optional
	)

	query := ecs.Query3[Children, Local, Global](world)
	return ecs.NewSystem(func(dt time.Duration) {
		queryTopLevel.MapId(func(id ecs.Id, children *Children, local *Local, global *Global) {
			// Resolve the top level transform (No movement, rotation, or scale)
			if local != nil {
				global.Transform = local.Transform
			}

			if children == nil {
				return
			}
			resolveTransform(query, children, global)
		})
	})
}

// TODO: This might loop forever if you have malformed heirarchies where there are cycles. I dont have any prevention for that right now.
func resolveTransform(
	query *ecs.View3[Children, Local, Global],
	children *Children, // TODO: Pointer here?
	parentGlobal *Global, // TODO: Pointer here?
) {
	// for i := range children.List {
	// 	nextChildren, childLocal, childGlobal := query.Read(children.List[i])
	// 	if childGlobal == nil { continue } // If child has no transform, skip
	// 	if childLocal == nil {
	// 		// If child has no local transform. Assume it is identity. And just copy parentGlobal
	// 		childGlobal.Pos = parentGlobal.Pos
	// 		childGlobal.Rot = parentGlobal.Rot
	// 		childGlobal.Scale = parentGlobal.Scale
	// 	} else {
	// 		*childGlobal = childLocal.Globalize(*parentGlobal)
	// 	}

	// 	if nextChildren == nil { continue } // Dont recurse: This child doesn't have a TransformChildren component
	// 	if len(nextChildren.List) > 0 {
	// 		resolveTransform(query, nextChildren, childGlobal)
	// 	}
	// }

	for _, childId := range children.MiniSlice.All() {
		nextChildren, childLocal, childGlobal := query.Read(childId)
		if childGlobal == nil {
			continue // If child has no transform, skip
		}
		if childLocal == nil {
			// If child has no local transform. Assume it is identity. And just copy parentGlobal
			childGlobal.Pos = parentGlobal.Pos
			childGlobal.Rot = parentGlobal.Rot
			childGlobal.Scale = parentGlobal.Scale
		} else {
			*childGlobal = childLocal.Globalize(*parentGlobal)
		}

		if nextChildren == nil {
			continue // Dont recurse: This child doesn't have a TransformChildren component
		}
		if nextChildren.MiniSlice.Len() > 0 {
			resolveTransform(query, nextChildren, childGlobal)
		}
	}

}

// Note: Originally you added this just because you wanted transform parenting with projectiles
// TODO: Potential optimization ideas:
// 1. You can sort the archetype storage by ECS ID, then make sure all children have higher IDs than parents
// 2. Read just the ECS IDs and do query.Read operations, rather than MapId
// Level Based idea: Iterate levels. I think I'd prefer to do it recursively through children
// type TransformParent struct {
// 	Level uint8
// 	Id ecs.Id
// }
// func ResolveTransformHeirarchySystem(game *Game) ecs.System {
// 	query := ecs.Query3[TransformParent, Transform, GlobalTransform](game.World)
// 	queryTopLevel := ecs.Query2[Transform, GlobalTransform](game.World,
// 		ecs.Without(TransformParent{}))
// 	return ecs.NewSystem(func(dt time.Duration) {
// 		queryTopLevel.MapId(func(id ecs.Id, local *Transform, global *GlobalTransform) {
// 			global.Position = local.Position
// 		})

// 		currentLevel := uint8(0)
// 		for {
// 			noneFound := true
// 			query.MapId(func(id ecs.Id, parent *TransformParent, local *Transform, global *GlobalTransform) {
// 				if parent.Level == currentLevel { return } // Skip everything not on our current level
// 				noneFound = false

// 				_, _, parentGlobal := query.Read(parent.Id)

// 				global.Position = local.Position.Add(parentGlobal.Position)
// 			})

// 			currentLevel++
// 			if noneFound {
// 				break
// 			}
// 			if currentLevel == math.MaxUint8 {
// 				break
// 			}
// 		}
// 	})
// }
