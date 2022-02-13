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
		camera.Position[0] += (cursor.DragStart[0] - mX) / float32(camera.Zoom)
		camera.Position[1] += (cursor.DragStart[1] - mY) / float32(camera.Zoom)

		cursor.DragStart = glitch.Vec2{mX, mY}
		cursor.Dragging = true
	}

	if !win.Pressed(dragButton) {
		cursor.Dragging = false
	}
}
