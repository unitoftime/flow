package transform

import (
	"github.com/unitoftime/ecs"
)

var ChildrenComp = ecs.NewComp[Children]()

func (c Children) CompId() ecs.CompId {
	return ChildrenComp.CompId()
}

func (c Children) CompWrite(w ecs.W) {
	ChildrenComp.WriteVal(w, c)
}

var GlobalComp = ecs.NewComp[Global]()

func (c Global) CompId() ecs.CompId {
	return GlobalComp.CompId()
}

func (c Global) CompWrite(w ecs.W) {
	GlobalComp.WriteVal(w, c)
}

var LocalComp = ecs.NewComp[Local]()

func (c Local) CompId() ecs.CompId {
	return LocalComp.CompId()
}

func (c Local) CompWrite(w ecs.W) {
	LocalComp.WriteVal(w, c)
}

var ParentComp = ecs.NewComp[Parent]()

func (c Parent) CompId() ecs.CompId {
	return ParentComp.CompId()
}

func (c Parent) CompWrite(w ecs.W) {
	ParentComp.WriteVal(w, c)
}
