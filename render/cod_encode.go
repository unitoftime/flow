package render

import (
	"github.com/unitoftime/ecs"
)

var CalculatedVisibilityComp = ecs.NewComp[CalculatedVisibility]()

func (c CalculatedVisibility) CompId() ecs.CompId {
	return CalculatedVisibilityComp.CompId()
}

func (c CalculatedVisibility) CompWrite(w ecs.W) {
	CalculatedVisibilityComp.WriteVal(w, c)
}

var CameraComp = ecs.NewComp[Camera]()

func (c Camera) CompId() ecs.CompId {
	return CameraComp.CompId()
}

func (c Camera) CompWrite(w ecs.W) {
	CameraComp.WriteVal(w, c)
}

var SpriteComp = ecs.NewComp[Sprite]()

func (c Sprite) CompId() ecs.CompId {
	return SpriteComp.CompId()
}

func (c Sprite) CompWrite(w ecs.W) {
	SpriteComp.WriteVal(w, c)
}

var TargetComp = ecs.NewComp[Target]()

func (c Target) CompId() ecs.CompId {
	return TargetComp.CompId()
}

func (c Target) CompWrite(w ecs.W) {
	TargetComp.WriteVal(w, c)
}

var VisibilityComp = ecs.NewComp[Visibility]()

func (c Visibility) CompId() ecs.CompId {
	return VisibilityComp.CompId()
}

func (c Visibility) CompWrite(w ecs.W) {
	VisibilityComp.WriteVal(w, c)
}

var VisionListComp = ecs.NewComp[VisionList]()

func (c VisionList) CompId() ecs.CompId {
	return VisionListComp.CompId()
}

func (c VisionList) CompWrite(w ecs.W) {
	VisionListComp.WriteVal(w, c)
}

var WindowComp = ecs.NewComp[Window]()

func (c Window) CompId() ecs.CompId {
	return WindowComp.CompId()
}

func (c Window) CompWrite(w ecs.W) {
	WindowComp.WriteVal(w, c)
}
