package render

import (
	// "math"

	"github.com/unitoftime/flow/glm"
	"github.com/unitoftime/glitch"
)

type Camera struct {
	Camera   *glitch.CameraOrtho
	Position glm.Vec2
	Zoom     float64
	bounds   glm.Rect
}

func NewCamera(bounds glm.Rect, x, y float64) *Camera {
	return &Camera{
		Camera:   glitch.NewCameraOrtho(),
		Position: glm.Vec2{x, y},
		Zoom:     1.0,
		bounds:   bounds,
	}
}

func (c *Camera) Update(bounds glm.Rect) {
	// // Snap camera
	// c.Position[0] = float32(math.Round(float64(c.Position[0])))
	// c.Position[1] = float32(math.Round(float64(c.Position[1])))

	c.bounds = bounds

	// TODO - Note: This is just to center point (0, 0), this should be selected some other way
	screenCenter := bounds.Center()

	c.Camera.SetOrtho2D(bounds)
	movePos := glm.Vec2{c.Position.X, c.Position.Y}.Sub(screenCenter)

	c.Camera.SetView2D(movePos.X, movePos.Y, c.Zoom, c.Zoom)
}

func (c *Camera) Project(point glitch.Vec3) glitch.Vec3 {
	return c.Camera.Project(point)
}

func (c *Camera) Unproject(point glitch.Vec3) glitch.Vec3 {
	return c.Camera.Unproject(point)
}

func (c *Camera) WorldSpaceRect() glm.Rect {
	box := c.bounds.ToBox()
	min := c.Unproject(box.Min)
	max := c.Unproject(box.Max)

	return glm.R(min.X, min.Y, max.X, max.Y)
}

// type Camera struct {
// 	win *pixelgl.Window
// 	Position pixel.Vec
// 	Zoom float64
// 	mat pixel.Matrix
// }

// func NewCamera(win *pixelgl.Window, x, y float64) *Camera {
// 	return &Camera{
// 		win: win,
// 		Position: pixel.V(x, y),
// 		Zoom: 1.0,
// 		mat: pixel.IM,
// 	}
// }

// func (c *Camera) Update() {
// 	screenCenter := c.win.Bounds().Center()

// 	movePos := pixel.V(-c.Position.X, -c.Position.Y).Add(screenCenter)
// 	c.mat = pixel.IM.Moved(movePos).Scaled(screenCenter, c.Zoom)
// }

// func (c *Camera) Mat() pixel.Matrix {
// 	return c.mat
// }
