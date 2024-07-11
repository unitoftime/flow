package render

import (
	"github.com/unitoftime/glitch"
)

type Cursor struct {
	Dragging bool
	DragStart glitch.Vec2
}

var cursor Cursor

func MouseDrag(win *glitch.Window, camera *Camera, dragButton glitch.Key) {
	mX, mY := win.MousePosition()

	if win.JustPressed(dragButton) {
		cursor.DragStart = glitch.Vec2{mX, mY}
	}

	if win.Pressed(dragButton) {
		// this is the right ratio, but you have to pass this in someway maybe camera knows about FB scaling? Idk
		// camera.Position[0] += (cursor.DragStart[0] - mX) / float32(camera.Zoom) * (camera.bounds.W() / 1920.0)
		// camera.Position[1] += (cursor.DragStart[1] - mY) / float32(camera.Zoom) * (camera.bounds.W() / 1920.0)
		camera.Position.X += (cursor.DragStart.X - mX) / camera.Zoom
		camera.Position.Y += (cursor.DragStart.Y - mY) / camera.Zoom

		cursor.DragStart = glitch.Vec2{mX, mY}
		cursor.Dragging = true
	}

	if !win.Pressed(dragButton) {
		cursor.Dragging = false
	}
}
