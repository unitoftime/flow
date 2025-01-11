package glm

type IVec2 struct {
	X, Y int
}

func (v IVec2) Add(v2 IVec2) IVec2 {
	return IVec2{v.X+v2.X, v.Y+v2.Y}
}

func (v IVec2) Sub(v2 IVec2) IVec2 {
	return IVec2{v.X-v2.X, v.Y-v2.Y}
}
