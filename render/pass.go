package render

import (
	"time"

	"github.com/unitoftime/ecs"
	"github.com/unitoftime/flow/glm"
	"github.com/unitoftime/flow/spatial"
	"github.com/unitoftime/flow/transform"
	"github.com/unitoftime/glitch"
)

// Systems:
// - Update cameras
// - Calculate vision list for each camera
//   - (Optimization) - If a lot of render passes, you might build a spatial map

// Clear, Draw, Next
// - Calculate all of the render passes that need to be executed
// - Calculate the visible entities for each render pass
//   - (Optimization) - If a lot of render passes, you might build a spatial map
// - Loop through render passes and execute them

// Clear a target
// Add commands to something
// Flush to that target
// Go next

// Render Pass Structure
// - RenderPass
// |-- Render Target
// |-- Camera
//   |-- Visible Entity
//   |-- Visible Entity
// |-- Camera
//   |-- Visible Entity
//   |-- Visible Entity

func UpdateCameraSystem(dt time.Duration, query *ecs.View2[Camera, Target]) {
	query.MapId(func(id ecs.Id, camera *Camera, target *Target) {
		camera.Update(target.draw.Bounds())
	})
}

func CalculateVisibilitySystem(dt time.Duration, camQuery *ecs.View2[Camera, VisionList], query *ecs.View3[transform.Global, Sprite, Visibility]) {

	// For each camera, calculate it's frustrum, clear its vision list, and add every visible sprite
	camQuery.MapId(func(_ ecs.Id, camera *Camera, visionList *VisionList) {
		winShape := camera.WorldSpaceShape()
		visionList.Clear()

		query.MapId(func(id ecs.Id, gt *transform.Global, sprite *Sprite, vis *Visibility) {
			vis.Calculated = false
			if vis.Hide { return }

			mat := gt.Mat4()
			shape := spatial.Rect(sprite.Bounds(), mat)

			vis.Calculated = shape.Intersects(winShape)

			if vis.Calculated {
				visionList.Add(id)
			}
		})
	})


	// TODO: Optimization if there are a lot of cameras: Preconfigure a spatial hash
	// chunkSize := [2]int{1024/8, 1024/8}
	// visionMap := spatial.NewHashmap[ecs.Id](chunkSize, 8)
	// query.MapId(func(id ecs.Id, gt *transform.Global, sprite *Sprite, vis *Visibility, calcVis *CalculatedVisibility) {
	// 	mat := gt.Mat4()
	// 	shape := spatial.Rect(sprite.sprite.Bounds(), mat)

	// 	visionMap.Add(shape, id)
	// })
}

type RenderPassList struct {
	List []RenderPass
}

func NewRenderPassList() RenderPassList {
	return RenderPassList{
		List: make([]RenderPass, 0),
	}
}
func (l *RenderPassList) Add(p RenderPass) {
	l.List = append(l.List, p)
}
func (l *RenderPassList) Clear() {
	l.List = l.List[:0]
}

type RenderTarget interface{
	glitch.Target
	glitch.BatchTarget
	Bounds() glm.Rect
}

type RenderPass struct {
	batchTarget glitch.BatchTarget
	drawTarget RenderTarget
	clearColor glm.RGBA
	camera Camera // TODO: I feel like there could/should support multiple?
	visionList VisionList
}

func CalculateRenderPasses(dt time.Duration, passes *RenderPassList, query *ecs.View3[Camera, Target, VisionList]) {
	passes.Clear()
	query.MapId(func(id ecs.Id, camera *Camera, target *Target, visionList *VisionList) {
		passes.Add(RenderPass{
			drawTarget: target.draw,
			batchTarget: target.batch,
			camera: *camera,
			visionList: *visionList,
		})
	})

	// TODO: Sort somehow? Priority? Order added?
}


func ExecuteRenderPass(dt time.Duration, passes *RenderPassList, query *ecs.View2[transform.Global, Sprite]) {
	for _, pass := range passes.List {
		glitch.Clear(pass.drawTarget, pass.clearColor)

		camera := &pass.camera.Camera
		glitch.SetCamera(camera)

		for _, id := range pass.visionList.List {
			gt, sprite := query.Read(id)
			if gt == nil { continue }
			if sprite == nil { continue }

			mat := gt.Mat4()
			sprite.Sprite.DrawColorMask(pass.batchTarget, mat, glm.White)
		}
	}
}


func RenderingFlushSystem(dt time.Duration, query *ecs.View1[Window]) {
	query.MapId(func(_ ecs.Id, win *Window) {
		win.Update()
	})
}
