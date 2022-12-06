package render

import (
	"time"
	"image/color"

	"github.com/unitoftime/glitch"
	// "github.com/unitoftime/ecs"
	"github.com/unitoftime/flow/phy2"
)

// This is an animation frame
type Frame struct {
	sprite *glitch.Sprite
	Dur time.Duration
	mount map[string]glitch.Vec2
	MirrorY bool
	// MirrorX bool
}
func NewFrame(sprite *glitch.Sprite, dur time.Duration) Frame {
	return Frame{
		sprite: sprite,
		Dur: dur,
		mount: make(map[string]glitch.Vec2),
	}
}

func (f *Frame) SetMount(name string, point glitch.Vec2) {
	f.mount[name] = point
}

func (f *Frame) Mount(name string) glitch.Vec2 {
	pos, ok := f.mount[name]
	if !ok {
		return glitch.Vec2{}
	}
	if f.MirrorY {
		pos[0] = -pos[0]
	}
	return pos
}

type Animation struct {
	frameIdx int
	remainingDur time.Duration
	frames map[string][]Frame // This is the map of all animations and their associated frames
	animName string
	curAnim []Frame // This is the current animation frames that we are operating on
	Color color.NRGBA // TODO - performance on interfaces vs structs?
	Scale glitch.Vec2
	// translation glitch.Vec3
}

func NewAnimation(startingAnim string, frames map[string][]Frame) Animation {
	anim := Animation{
		frames: frames,
		Color: color.NRGBA{255, 255, 255, 255},
		Scale: glitch.Vec2{1, 1},
	}
	anim.SetAnimation(startingAnim)
	return anim
}

// func (a *Animation) SetTranslation(pos glitch.Vec3) {
// 	a.translation = pos
// }

func (a *Animation) SetAnimation(name string) {
	if name == a.animName { return } // Skip if we aren't actually changing the animation

	newAnim, ok := a.frames[name]
	if !ok { return }

	a.animName = name
	a.curAnim = newAnim
	a.SetFrame(0)
}

func (a *Animation) NextFrame() {
	a.SetFrame(a.frameIdx + 1)
}

func (a *Animation) SetFrame(idx int) {
	a.frameIdx = idx % len(a.curAnim)
	frame := a.curAnim[a.frameIdx]
	a.remainingDur = frame.Dur
}

func (a *Animation) GetFrame() Frame {
	idx := a.frameIdx % len(a.curAnim)
	return a.curAnim[idx]
}

// Steps the animation forward by dt amount of time
func (anim *Animation) Update(dt time.Duration) {
	anim.remainingDur -= dt
	if anim.remainingDur < 0 {
		// Change to a new animation frame
		anim.NextFrame()
	}
}

// Draws the animation to the render pass
func (anim *Animation) Draw(target glitch.BatchTarget, pos *phy2.Pos) {
	frame := anim.curAnim[anim.frameIdx]

	// frame.sprite.SetTranslation(anim.translation)

	mat := glitch.Mat4Ident
	mat.Scale(anim.Scale[0], anim.Scale[1], 1.0)
	if frame.MirrorY {
		mat.Scale(-1.0, 1.0, 1.0)
	}

	mat.Translate(float32(pos.X), float32(pos.Y), 0)
	// TODO - I think there's some mistakes here with premultiplied vs non premultiplied alpha
	col := glitch.RGBA{float32(anim.Color.R)/255.0, float32(anim.Color.G)/255.0, float32(anim.Color.B)/255.0, float32(anim.Color.A)/255.0}

	frame.sprite.DrawColorMask(target, mat, col)
}
