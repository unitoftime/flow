package render

import (
	"time"
	// "image/color"

	"github.com/unitoftime/glitch"
	"github.com/unitoftime/flow/phy2"
)

// TODO - it might make more sense to make this like an aseprite wrapper object that has layers, frames, tags, etc

// This is an animation frame
type Frame struct {
	Sprite *glitch.Sprite
	// Origin phy2.Pos
	Dur time.Duration
	mount map[string]phy2.Pos // TODO - this is just kind of arbitrary data for my mountpoint system
}
func NewFrame(sprite *glitch.Sprite, dur time.Duration) Frame {
	return Frame{
		Sprite: sprite,
		Dur: dur,
		mount: make(map[string]phy2.Pos),
	}
}

func (f *Frame) Bounds() glitch.Rect {
	return f.Sprite.Bounds()
}

func (f *Frame) SetMount(name string, point phy2.Pos) {
	f.mount[name] = point
}

func (f *Frame) Mount(name string) phy2.Pos {
	pos, ok := f.mount[name]
	if !ok {
		return phy2.Pos{}
	}
	return pos
}

type Animation struct {
	frameIdx int
	remainingDur time.Duration
	frames map[string][]Frame // This is the map of all animations and their associated frames
	animName string
	curAnim []Frame // This is the current animation frames that we are operating on
	// Color color.NRGBA // TODO - performance on interfaces vs structs?
	// Rotation float64
	// Scale glitch.Vec2
	Loop bool
	speed float64 // This is used to scale the duration of the animation evenly so that the animation can fit a certain time duration
	// translation glitch.Vec3

	// MirrorX bool // TODO
	MirrorY bool // Mirror around the Y axis
}

func NewAnimation(startingAnim string, frames map[string][]Frame) Animation {
	anim := Animation{
		frames: frames,
		// Color: color.NRGBA{255, 255, 255, 255},
		// Scale: glitch.Vec2{1, 1},
		Loop: true,
		speed: 1.0,
	}
	if startingAnim == "" {
		// Just set some random animation if unset
		for name := range frames {
			anim.SetAnimation(name)
			break
		}
	} else {
		anim.SetAnimation(startingAnim)
	}

	return anim
}

// func (a *Animation) SetTranslation(pos glitch.Vec3) {
// 	a.translation = pos
// }

func (a *Animation) SetAnimationWithDuration(name string, dur time.Duration) {
	a.SetAnimation(name)
	totalAnimTime := 0 * time.Second
	for _, frame := range a.curAnim {
		totalAnimTime += frame.Dur
	}
	a.speed = totalAnimTime.Seconds() / dur.Seconds()
}

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
	if a.Loop {
		// If the idx is passed the animation, then loop it
		a.frameIdx = idx % len(a.curAnim)
		frame := a.curAnim[a.frameIdx]
		a.remainingDur = frame.Dur
	} else {
		// If the idx is passed the animation, snap to the last frame
		if idx >= len(a.curAnim) {
			idx = len(a.curAnim) - 1
		}
		a.frameIdx = idx
		frame := a.curAnim[a.frameIdx]
		a.remainingDur = frame.Dur
	}
}

func (a *Animation) GetFrame() Frame {
	idx := a.frameIdx % len(a.curAnim)
	return a.curAnim[idx]
}

// Steps the animation forward by dt amount of time
// Returns true if the animation frame has changed, else returns false
func (anim *Animation) Update(dt time.Duration) bool {
	adjustedDt := time.Duration(1_000_000_000 * anim.speed * dt.Seconds())

	anim.remainingDur -= adjustedDt
	if anim.remainingDur < 0 {
		// Change to a new animation frame
		anim.NextFrame()
		return true
	}
	return false
}

// // Draws the animation to the render pass
// func (anim *Animation) Draw(target glitch.BatchTarget, pos Pos) {
// 	frame := anim.curAnim[anim.frameIdx]

// 	// frame.sprite.SetTranslation(anim.translation)

// 	mat := glitch.Mat4Ident
// 	// mat.Translate(float32(frame.Origin.X), float32(frame.Origin.Y), 0)
// 	mat.Scale(anim.Scale[0], anim.Scale[1], 1.0)
// 	if anim.MirrorY {
// 		mat.Scale(-1.0, 1.0, 1.0)
// 	}
// 	mat.Rotate(anim.Rotation, glitch.Vec3{0, 0, 1})
// 	mat.Translate(pos.X, pos.Y, 0)

// 	// TODO - I think there's some mistakes here with premultiplied vs non premultiplied alpha
// 	col := glitch.RGBA{anim.Color.R/255.0, (anim.Color.G)/255.0, (anim.Color.B)/255.0, (anim.Color.A)/255.0}

// 	frame.Sprite.DrawColorMask(target, mat, col)
// }
