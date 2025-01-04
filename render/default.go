package render

import (
	"time"

	"github.com/unitoftime/ecs"
	"github.com/unitoftime/flow/glm"
	"github.com/unitoftime/glitch"
)

//go:generate go run ../../cod/cmd/cod

type DefaultPlugin struct {
}

func (p DefaultPlugin) Initialize(world *ecs.World) {
	scheduler := ecs.GetResource[ecs.Scheduler](world)

	rl := NewRenderPassList()
	ecs.PutResource(world, &rl)

	scheduler.AddSystems(ecs.StageStartup,
		ecs.NewSystem1(SetupRenderingSystem),
	)
	scheduler.AddSystems(ecs.StageUpdate,
		ecs.NewSystem1(UpdateCameraSystem),
		ecs.NewSystem2(CalculateVisibilitySystem),

		ecs.NewSystem2(CalculateRenderPasses),
		ecs.NewSystem2(ExecuteRenderPass),
		ecs.NewSystem1(RenderingFlushSystem),
	)
}

//cod:component
type VisionList struct {
	List []ecs.Id
}
func NewVisionList() VisionList {
	return VisionList{
		List: make([]ecs.Id, 0),
	}
}
func (vl *VisionList) Add(id ecs.Id) {
	vl.List = append(vl.List, id)
}
func (vl *VisionList) Clear() {
	vl.List = vl.List[:0]
}


//cod:component
type Window struct {
	*glitch.Window
}

//cod:component
type Sprite struct {
	Sprite *glitch.Sprite
}
func (s Sprite) Bounds() glm.Rect {
	return s.Sprite.Bounds()
}

//cod:component
type Target struct {
	draw RenderTarget
	batch glitch.BatchTarget // TODO: Indexed?
}

//cod:component
type Visibility struct {
	Hide bool // If set true, the entity will always be calculated as invisible
	Calculated bool
}

//cod:component
type CalculatedVisibility struct {
	Visible bool // If set true, the entity is visible
}

func SetupRenderingSystem(dt time.Duration, commands *ecs.CommandQueue) {
	defer commands.Execute() // TODO: Remove

	win, err := glitch.NewWindow(1920, 1080, "TODO", glitch.WindowConfig{
		Vsync: true,
	})
	if err != nil {
		panic(err)
	}

	commands.SpawnEmpty().
		Insert(Window{win})

	commands.SpawnEmpty().
		Insert(Camera2DBundle{win})
}
