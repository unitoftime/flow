package render

import (
	"github.com/jstewart7/glitch"
	// "github.com/faiface/pixel"
	// "github.com/faiface/pixel/pixelgl"
)

type Camera struct {
	win *glitch.Window
	Camera *glitch.CameraOrtho
	Position glitch.Vec2
	Zoom float32
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
	// screenCenter := c.win.Bounds().Center()

	// movePos := glitch.Vec2{c.Position[0], c.Position[1]}.Add(screenCenter)
	// c.mat = pixel.IM.Moved(movePos).Scaled(screenCenter, c.Zoom)
	c.Camera.SetOrtho2D(c.win)
	c.Camera.SetView2D(0, 0, 1.0, 1.0)
	// c.Camera.SetView2D(movePos[0], movePos[1], c.Zoom, c.Zoom)
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
