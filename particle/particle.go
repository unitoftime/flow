package particle

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	// "github.com/ungerik/go3d/float64/vec2"

	"github.com/unitoftime/ecs"

	// "github.com/unitoftime/flow/timer"
	"github.com/unitoftime/flow/glm"
	"github.com/unitoftime/flow/interp"
)

// + Position
// + Velocity
// - Acceleration

// + Size
// + Color
// + Type

// TODO
// - Rotation
// - Special Emitter

// func init() {
// 	gob.Register(ConstantPositioner{})
// 	gob.Register(CopyPositioner{})
// 	gob.Register(RingPositioner{})
// 	gob.Register(AnglePositioner{})
// 	gob.Register(PathPositioner{})
// 	gob.Register(PhysicsUpdater{})
// }

// TODO - Move to a math package once generics comes out
func Clamp(min, max, val float64) float64 {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

// --------------------------------------------------------------------------------------------------
// - Positioners
// --------------------------------------------------------------------------------------------------
type Lifetime struct {
	Total     time.Duration
	Remaining time.Duration
}

func NewLifetime(total time.Duration) Lifetime {
	return Lifetime{
		Total:     total,
		Remaining: total,
	}
}
func (l *Lifetime) Ratio() float64 {
	if l.Total == 0 {
		return 0
	}

	return 1 - Clamp(0, 1.0, l.Remaining.Seconds()/l.Total.Seconds())
}

type Color struct {
	Interp     interp.Interp
	Start, End color.NRGBA
}

func NewColor(interpolation interp.Interp, start, end color.NRGBA) Color {
	return Color{
		Interp: interpolation,
		Start:  start,
		End:    end,
	}
}

func (c *Color) Get(ratio float64) color.NRGBA {
	ratio = Clamp(0, 1.0, ratio)
	color := color.NRGBA{
		uint8(c.Interp.Uint8(c.Start.R, c.End.R, ratio)),
		uint8(c.Interp.Uint8(c.Start.G, c.End.G, ratio)),
		uint8(c.Interp.Uint8(c.Start.B, c.End.B, ratio)),
		uint8(c.Interp.Uint8(c.Start.A, c.End.A, ratio)),
	}
	return color
}

type Size struct {
	Interp     interp.Interp
	Start, End glm.Vec2
}

func NewSize(interpolation interp.Interp, start, end glm.Vec2) Size {
	return Size{
		Interp: interpolation,
		Start:  start,
		End:    end,
	}
}

func (s *Size) Get(ratio float64) glm.Vec2 {
	ratio = Clamp(0, 1.0, ratio)
	size := glm.Vec2{
		s.Interp.Float64(s.Start.X, s.End.X, ratio),
		s.Interp.Float64(s.Start.Y, s.End.Y, ratio),
	}

	return size
}

// --------------------------------------------------------------------------------------------------
// - ComponentFactory
// --------------------------------------------------------------------------------------------------
type PrefabBuilder interface {
	Add(*ecs.Entity)
}

// type RingBuilder struct {
// 	AngleRange glm.Vec2
// 	RadiusRange glm.Vec2
// }
// func (p *RingBuilder) Add(prefab *ecs.Entity) {
// 	angle := interp.Linear.Float64(p.AngleRange.X, p.AngleRange.Y, rand.Float64())
// 	radius := interp.Linear.Float64(p.RadiusRange.X, p.RadiusRange.Y, rand.Float64())

// 	// vec := vec2.UnitX
// 	vec := glm.Vec2{1, 0}
// 	vec.Scaled(radius).Rotated(angle)

// 	prefab.Add(ecs.C(glm.Pos{vec.X, vec.Y}))
// }

// // type AngleBuilder struct {
// // 	Scale float64
// // }
// // func (p *AngleBuilder) Add(prefab ecs.Entity) {
// // 	transform := prefab.Read(glm.Transform{}).(glm.Transform)
// // 	pos := vec2.T{transform.X, transform.Y}
// // 	prefab.Write(glm.Rigidbody{
// // 		Mass: 1,
// // 		Velocity: pos.Normalize().Scaled(p.Scale),
// // 	})
// // }

// // TODO - Should I just build this into the emitter?
// type LifetimeBuilder struct {
// 	Range glm.Vec2 // Specified in seconds
// }
// func (b *LifetimeBuilder) Add(prefab *ecs.Entity) {
// 	seconds := interp.Linear.Float64(b.Range.X, b.Range.Y, rand.Float64())
// 	prefab.Add(ecs.C(
// 		NewLifetime(time.Duration(seconds * 1000) * time.Millisecond),
// 	))
// }

// type TransformBuilder struct {
// 	PosPositioner Vec2Positioner
// }
// func (p *TransformBuilder) Add(prefab *ecs.Entity) {
// 	pos := p.PosPositioner.Vec2(glm.Vec2{})
// 	prefab.Add(ecs.C(glm.Pos{pos.X, pos.Y}))
// }

// type RigidbodyBuilder struct {
// 	Mass float64
// 	VelPositioner Vec2Positioner
// }
// func (b *RigidbodyBuilder) Add(prefab *ecs.Entity) {
// 	// transform := prefab.Read(glm.Transform{}).(glm.Transform)
// 	// transform, _ := ecs.ReadFromEntity[glm.Transform](prefab)
// 	// pos := glm.Vec2{transform.X, transform.Y}
// 	pos, _ := ecs.ReadFromEntity[glm.Pos](prefab)

// 	vel := b.VelPositioner.Vec2(glm.Vec2(pos))

// 	prefab.Add(ecs.C(glm.Rigidbody{
// 		Mass: b.Mass,
// 		Velocity: vel,
// 	}))
// }

// type ConstantBuilder struct {
// 	Component interface{}
// }
// func (b *RigidbodyBuilder) Add(prefab ecs.Entity) {
// }

// --------------------------------------------------------------------------------------------------
// - Positioners
// --------------------------------------------------------------------------------------------------
type Float64Positioner interface {
	Float64(count int, A float64) float64
}

type ConstFloat64Positioner struct {
	Val float64
}

func (p ConstFloat64Positioner) Float64(count int, v float64) float64 {
	return v + p.Val
}

type RandomFloat64Positioner struct {
}

func (p RandomFloat64Positioner) Float64(count int, v float64) float64 {
	return v + (2 * math.Pi * rand.Float64())
}

type Vec2Positioner interface {
	Vec2(count int, A glm.Vec2) glm.Vec2
}

type ConstantPositioner struct {
	X, Y float64
}

func (p ConstantPositioner) Vec2(count int, A glm.Vec2) glm.Vec2 {
	return glm.Vec2{p.X, p.Y}
}

type GeometricPositioner struct {
	Scale          glm.Vec2 // X is min, Y is max
	Offset         float64
	DistanceOffset float64
	DistanceMod    int
	DistanceRem    int
}

func (p *GeometricPositioner) Vec2(count int, A glm.Vec2) glm.Vec2 {
	theta := 2 * math.Pi * float64(count) * p.Offset

	// w := p.Scale.Y - p.Scale.X
	// scale := (w * rand.Float64()) + p.Scale.X
	modCount := count
	if p.DistanceMod != 0 {
		modCount = count % p.DistanceMod
	}
	if p.DistanceRem != 0 {
		modCount = count / p.DistanceRem
	}
	scale := float64(modCount) * p.DistanceOffset

	x := scale * math.Cos(theta)
	y := scale * math.Sin(theta)

	return glm.Vec2{x, y}.Add(A)
}

type CirclePositioner struct {
	Scale glm.Vec2 // X is min, Y is max
}

func (p *CirclePositioner) Vec2(count int, A glm.Vec2) glm.Vec2 {
	theta := 2 * math.Pi * rand.Float64()

	w := p.Scale.Y - p.Scale.X
	scale := (w * rand.Float64()) + p.Scale.X

	x := scale * math.Cos(theta)
	y := scale * math.Sin(theta)

	return glm.Vec2{x, y}.Add(A)
}

type RectPositioner struct {
	Min, Max glm.Vec2 // TODO - rectangle passed in
}

func (p RectPositioner) Vec2(count int, A glm.Vec2) glm.Vec2 {
	w := p.Max.X - p.Min.X
	h := p.Max.Y - p.Min.Y

	x := w*rand.Float64() + p.Min.X
	y := h*rand.Float64() + p.Min.Y

	return glm.Vec2{x, y}
}

type CopyPositioner struct {
	Scale float64
}

func (p *CopyPositioner) Vec2(count int, A glm.Vec2) glm.Vec2 {
	return A.Scaled(p.Scale)
}

type AnglePositioner struct {
	Scale float64
}

func (p *AnglePositioner) Vec2(count int, A glm.Vec2) glm.Vec2 {
	return A.Norm().Scaled(p.Scale)
}

type RingPositioner struct {
	AngleRange  glm.Vec2
	RadiusRange glm.Vec2
}

func (p *RingPositioner) Vec2(count int, A glm.Vec2) glm.Vec2 {
	angle := interp.Linear.Float64(p.AngleRange.X, p.AngleRange.Y, rand.Float64())
	radius := interp.Linear.Float64(p.RadiusRange.X, p.RadiusRange.Y, rand.Float64())

	vec := glm.Vec2{1, 0}
	vec = vec.Scaled(radius).Rotated(angle)
	return vec
}

// type FirePositioner struct {

// }
// func (p *FirePositioner) Vec2(A glm.Vec2) glm.Vec2 {
// 	return glm.Vec2{-A.X, 5}
// }

// An emitter is used to spawn particles in a certain way
type Emitter struct {
	// Max int
	Rate   float64 // Spawn how many per frame
	period int

	OneShot     bool
	Loop        bool
	Probability float64
	// Duration time.Duration
	// Type ParticleType
	Prefab *ecs.Entity

	// SizeCurve, ColorCurve interp.Interp

	// Timer timer.Timer

	Builders []PrefabBuilder

	// PosPositioner Vec2Positioner
	// VelPositioner Vec2Positioner
	// SizePositioner Vec2Positioner

	// PosBuilder PrefabBuilder
	// RbBuilder PrefabBuilder

	// RedPositioner Vec2Positioner
	// GreenPositioner Vec2Positioner
	// BluePositioner Vec2Positioner
	// AlphaPositioner Vec2Positioner
}

// func (e *Emitter) Update(world *ecs.World, position glm.Vec2, dt time.Duration) {
// 	count := 0
// 	// if e.OneShot {
// 	// 	count = e.Max
// 	// } else {
// 	// 	particles.EmissionCounter += dt
// 	// 	numParticles := math.Floor(particles.EmissionCounter.Seconds() * e.Rate)
// 	// 	particles.EmissionCounter -= time.Duration(math.Floor((numParticles / e.Rate) * 1e9)) * time.Nanosecond

// 	// 	count = int(numParticles)
// 	// }

// 	if e.Rate == 0 {
// 		return // Exit early if rate is set to 0
// 	}

// 	// 1/rate is the period, scaled to ms and then converted to duration
// 	// period := time.Duration(1000 * (1 / e.Rate)) * time.Millisecond
// 	period := int(1 / e.Rate)
// 	if period < 1.0 {
// 		count = int(e.Rate)
// 	} else {
// 		if e.period < 0 {
// 			e.period = period
// 			count = 1
// 		}

// 		e.period--
// 	}

// 	for i := 0; i < count; i++ {
// 		randP := rand.Float64()
// 		if randP < e.Probability {
// 			ok := e.Spawn(glm.Vec2{position.X, position.Y}, world)
// 			if !ok { break }
// 		}
// 	}

// 	// TODO - needs to be configurable
// 	// particles.Accel = ecs.Accelerator{pixel.V(position.X, position.Y)}
// }

// func (e *Emitter) Spawn(entPos glm.Vec2, world *ecs.World) bool {
// 	// If we don't loop, then only emit a Total equal to Max
// 	// if !e.Loop {
// 	// 	if p.Total > p.Max {
// 	// 		return false
// 	// 	}
// 	// }

// 	// Don't spawn if we're already full
// 	// if len(p.list) >= e.Max {
// 	// 	return false
// 	// }

// 	for i := range e.Builders {
// 		e.Builders[i].Add(e.Prefab)
// 	}

// 	// transform := e.Prefab.Read(glm.Transform{}).(glm.Transform)
// 	// transform, _ := ecs.ReadFromEntity[glm.Transform](e.Prefab)
// 	// transform.X += entPos.X
// 	// transform.Y += entPos.Y
// 	// e.Prefab.Add(ecs.C(transform))
// 	pos, _ := ecs.ReadFromEntity[glm.Pos](e.Prefab)
// 	pos.X += entPos.X
// 	pos.Y += entPos.Y
// 	e.Prefab.Add(ecs.C(pos))

// 	// sizes := e.SizePositioner.Vec2(vec2.Zero)

// 	id := world.NewId()
// 	e.Prefab.Write(world, id)

// 	return true
// }

// type ParticleType uint8

// type Particle struct {
// 	Position glm.Vec2
// 	Velocity glm.Vec2

// 	// Interpolation Values
// 	Size glm.Vec2
// 	Red glm.Vec2
// 	Green glm.Vec2
// 	Blue glm.Vec2
// 	Alpha glm.Vec2
// 	Type ParticleType
// 	MaxLife, Life time.Duration
// 	ratio float64 // Life ratio 0 = Full Life | 1 = No Life
// }

// func (p *Particle) GetSize(curve interp.Interp) float64 {
// 	// return Lerp(p.Size, p.ratio)

// 	// iValue := curve.Point(p.ratio)
// 	// return Lerp(p.Size, iValue.Y)

// 	return curve.Float64(p.Size.X, p.Size.Y, p.ratio)
// }

// func (p *Particle) GetColor(curve interp.Interp) color.NRGBA {
// 	// color := color.NRGBA{
// 	// 	uint8(Lerp(p.Red, p.ratio)),
// 	// 	uint8(Lerp(p.Green, p.ratio)),
// 	// 	uint8(Lerp(p.Blue, p.ratio)),
// 	// 	uint8(Lerp(p.Alpha, p.ratio)),
// 	// }

// 	// iValue := curve.Point(p.ratio)
// 	// color := color.NRGBA{
// 	// 	uint8(Lerp(p.Red, iValue.Y)),
// 	// 	uint8(Lerp(p.Green, iValue.Y)),
// 	// 	uint8(Lerp(p.Blue, iValue.Y)),
// 	// 	uint8(Lerp(p.Alpha, iValue.Y)),
// 	// }

// 	color := color.NRGBA{
// 		uint8(curve.Float64(p.Red.X, p.Red.Y, p.ratio)),
// 		uint8(curve.Float64(p.Green.X, p.Green.Y, p.ratio)),
// 		uint8(curve.Float64(p.Blue.X, p.Blue.Y, p.ratio)),
// 		uint8(curve.Float64(p.Alpha.X, p.Alpha.Y, p.ratio)),
// 	}
// 	return color
// }

// type Accelerator struct {
// 	Position glm.Vec2
// }

// func (a *Accelerator) GetAcceleration(p *Particle) glm.Vec2 {
// 	vec := vec2.Sub(&a.Position, &p.Position)
// 	return vec.Scaled(0.01)
// }

// type ParticleUpdater interface {
// 	Update(*Particle, time.Duration)
// }

// type PhysicsUpdater struct {
// 	// Acceleration Type (if we want)
// }

// func (u PhysicsUpdater) Update(p *Particle, dt time.Duration) {
// 	delta := p.Velocity.Scaled(dt.Seconds())
// 	p.Position.Add(&delta)
// }

// type PathUpdater struct {
// 	Path []glm.Vec2
// }

// func (u PathUpdater) Update(p *Particle, dt time.Duration) {
// 	// delta := p.Velocity.Scaled(dt.Seconds())
// 	// p.Position.Add(&delta)
// }

// type Group struct {
// 	Max int
// 	Total int
// 	EmissionCounter time.Duration

// 	Updater ParticleUpdater

// 	PosPositioner Vec2Positioner
// 	SizeCurve interp.Interp
// 	ColorCurve interp.Interp

// 	list []Particle
// }

// func NewGroup(initSize int, updater ParticleUpdater, sizeCurve, colorCurve interp.Interp) Group {
// 	return Group{
// 		Max: initSize,
// 		Updater: updater,
// 		SizeCurve: sizeCurve,
// 		ColorCurve: colorCurve,
// 		list: make([]Particle, initSize),
// 	}
// }

// func (g *Group) List() []Particle {
// 	return g.list
// }

// func (g *Group) Update(dt time.Duration) {
// 	for i := range g.list {
// 		g.list[i].Life -= dt

// 		g.Updater.Update(&g.list[i], dt)
// 	}

// 	// This loop removes p whose life has expired
// 	i := 0
// 	for {
// 		if i >= len(g.list) { break }
// 		if g.list[i].Life <= 0 {
// 			// If our life is over then we get removed
// 			// Move last element to this position
// 			g.list[i] = g.list[len(g.list)-1]
// 			g.list = g.list[:len(g.list)-1]
// 		} else {
// 			i++
// 		}
// 	}
// }

// type PathPositioner struct {
// 	Path []glm.Vec2
// 	Lengths []float64
// 	TotalLength float64
// 	Variation glm.Vec2
// }

// func RandomPath(start, end glm.Vec2, n int, pathVariation float64, variation glm.Vec2) *PathPositioner {
// 	path := make([]glm.Vec2, n)

// 	// iValues := make(float64, n)
// 	// iValues.X = 0 // Interpolate to start point
// 	// iValues[len(iValues)-1] = 1.0 // Interpolate to end point

// 	path.X = start
// 	path[len(path)-1] = end

// 	nVec := vec2.Sub(&end, &start)

// 	latVec := nVec.Normalize().Rotate90DegLeft()

// 	for i := 1; i < n-2; i++ {
// 		interpVec := vec2.Interpolate(&start, &end, rand.Float64())

// 		rnd := 2 * (rand.Float64() - 0.5) * pathVariation
// 		lateral := latVec.Scaled(rnd)

// 		path[i] = vec2.Add(&interpVec, &lateral)
// 	}

// 	sort.Slice(path, func(i, j int) bool {
// 		ii := vec2.Sub(&start, &path[i])
// 		jj := vec2.Sub(&start, &path[j])
// 		return ii.LengthSqr() < jj.LengthSqr()
// 	})

// 	return StraightPath(path, variation)
// }

// func StraightPath(path []glm.Vec2, variation glm.Vec2) *PathPositioner {
// 	lengths := make([]float64, 0, len(path)-1)
// 	totalLength := 0.0
// 	for i := 0; i < len(path)-1; i++ {
// 		v := vec2.Sub(&path[i+1], &path[i])
// 		d := v.Length()
// 		lengths = append(lengths, d)
// 		totalLength += d
// 	}

// 	return &PathPositioner{
// 		Path: path,
// 		Lengths: lengths,
// 		TotalLength: totalLength,
// 		Variation: variation,
// 	}
// }

// func (p *PathPositioner) Vec2(A glm.Vec2) glm.Vec2 {
// 	rnd := rand.Float64() * p.TotalLength

// 	// Find a random interpolation value along the TotalLength
// 	// index -> index of the path point before our interp value
// 	// rnd -> Becomes the inner path's interp value
// 	index := 0
// 	for i,l := range p.Lengths {
// 		if rnd < l {
// 			index = i
// 			break
// 		} else {
// 			rnd -= l
// 		}
// 	}

// 	v := vec2.Sub(&p.Path[index+1], &p.Path[index])
// 	rndVal := v.Normalize().Scale(rnd)
// 	point := vec2.Add(&p.Path[index], rndVal)

// 	variation := glm.Vec2{
// 		p.Variation.X * (rand.Float64() * 2 - 1),
// 		p.Variation.Y * (rand.Float64() * 2 - 1),
// 	}
// 	point.Add(&variation)
// 	return point
// }
