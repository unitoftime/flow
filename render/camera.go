package render

import (
	"github.com/unitoftime/ecs"
	"github.com/unitoftime/flow/glm"
	"github.com/unitoftime/flow/spatial"
	"github.com/unitoftime/glitch"
)

// Bevy Reference
// https://docs.rs/bevy_core_pipeline/0.15.1/src/bevy_core_pipeline/core_2d/camera_2d.rs.html#34

// TODO: Make a cod macro to generate stuff like this? ie use required components instead of manually setting up a bundle
type Camera2DBundle struct {
	target RenderTarget
}

func (c Camera2DBundle) CompWrite(w ecs.W) {
	cam := NewCamera2D(glm.CR(0), 0, 0)
	target := Target{
		draw:  c.target,
		batch: c.target,
	}
	visionList := NewVisionList()

	cam.CompWrite(w)
	target.CompWrite(w)
	visionList.CompWrite(w)
}

//cod:component
type Camera struct {
	Camera   glitch.CameraOrtho
	Position glm.Vec2
	Zoom     float64
	bounds   glm.Rect
}

func NewCamera2D(bounds glm.Rect, x, y float64) Camera {
	orthoCam := glitch.NewCameraOrtho()
	return Camera{
		Camera:   *orthoCam,
		Position: glm.Vec2{x, y},
		Zoom:     1.0,
		bounds:   bounds,
	}
}

func NewCamera(bounds glm.Rect, x, y float64) *Camera {
	cam := NewCamera2D(bounds, x, y)
	return &cam
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

func (c *Camera) WorldSpaceShape() spatial.Shape {
	cameraMatrix := c.Camera.GetInverseMat4()
	return spatial.Rect(c.bounds, cameraMatrix)
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
