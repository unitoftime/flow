package render

import (
	"github.com/unitoftime/glitch"
)

type Camera struct {
	win *glitch.Window
	Camera *glitch.CameraOrtho
	Position glitch.Vec2
	Zoom float64
}

func NewCamera(win *glitch.Window, x, y float32) *Camera {
	return &Camera{
		win: win,
		Camera: glitch.NewCameraOrtho(),
		Position: glitch.Vec2{x, y},
		Zoom: 1.0,
	}
}

func (c *Camera) Update() {
	// TODO - Note: This is just to center point (0, 0), this should be selected some other way
	screenCenter := c.win.Bounds().Center()

	c.Camera.SetOrtho2D(c.win)
	movePos := glitch.Vec2{c.Position[0], c.Position[1]}.Sub(screenCenter)

	c.Camera.SetView2D(movePos[0], movePos[1], float32(c.Zoom), float32(c.Zoom))
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
