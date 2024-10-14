package glm

import "github.com/go-gl/mathgl/mgl64"

type Vec3 struct {
	X, Y, Z float64
}

type Vec4 struct {
	X, Y, Z, W float64
}

type Box struct {
	Min, Max Vec3
}

type Mat2 [4]float64
type Mat3 [9]float64
type Mat4 [16]float64

// This is in column major order
var Mat3Ident Mat3 = Mat3{
	1.0, 0.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 0.0, 1.0,
}

// This is in column major order
var Mat4Ident Mat4 = Mat4{
	1.0, 0.0, 0.0, 0.0,
	0.0, 1.0, 0.0, 0.0,
	0.0, 0.0, 1.0, 0.0,
	0.0, 0.0, 0.0, 1.0,
}

func (m *Mat3) Translate(x, y float64) *Mat3 {
	m[i3_2_0] = m[i3_2_0] + x
	m[i3_2_1] = m[i3_2_1] + y
	return m
}

// Note: Scales around 0,0
func (m *Mat4) Scale(x, y, z float64) *Mat4 {
	m[i4_0_0] = m[i4_0_0] * x
	m[i4_1_0] = m[i4_1_0] * x
	m[i4_2_0] = m[i4_2_0] * x
	m[i4_3_0] = m[i4_3_0] * x

	m[i4_0_1] = m[i4_0_1] * y
	m[i4_1_1] = m[i4_1_1] * y
	m[i4_2_1] = m[i4_2_1] * y
	m[i4_3_1] = m[i4_3_1] * y

	m[i4_0_2] = m[i4_0_2] * z
	m[i4_1_2] = m[i4_1_2] * z
	m[i4_2_2] = m[i4_2_2] * z
	m[i4_3_2] = m[i4_3_2] * z

	return m
}

func (m *Mat4) Translate(x, y, z float64) *Mat4 {
	m[i4_3_0] = m[i4_3_0] + x
	m[i4_3_1] = m[i4_3_1] + y
	m[i4_3_2] = m[i4_3_2] + z
	return m
}

func (m *Mat4) GetTranslation() Vec3 {
	return Vec3{m[i4_3_0], m[i4_3_1], m[i4_3_2]}
}

// https://github.com/go-gl/mathgl/blob/v1.0.0/mgl32/transform.go#L159
func (m *Mat4) Rotate(angle float64, axis Vec3) *Mat4 {
	// quat := mgl32.Mat4ToQuat(mgl32.Mat4(*m))
	// return &retMat
	rotation := Mat4(mgl64.HomogRotate3D(angle, mgl64.Vec3{axis.X, axis.Y, axis.Z}))
	// retMat := Mat4(mgl32.Mat4(*m).)
	// return &retMat
	mNew := m.Mul(&rotation)
	*m = *mNew
	return m
}

// Note: This modifies in place
func (m *Mat4) Mul(n *Mat4) *Mat4 {
	// This is in column major order
	*m = Mat4{
		// return &Mat4{
		// Column 0
		m[i4_0_0]*n[i4_0_0] + m[i4_1_0]*n[i4_0_1] + m[i4_2_0]*n[i4_0_2] + m[i4_3_0]*n[i4_0_3],
		m[i4_0_1]*n[i4_0_0] + m[i4_1_1]*n[i4_0_1] + m[i4_2_1]*n[i4_0_2] + m[i4_3_1]*n[i4_0_3],
		m[i4_0_2]*n[i4_0_0] + m[i4_1_2]*n[i4_0_1] + m[i4_2_2]*n[i4_0_2] + m[i4_3_2]*n[i4_0_3],
		m[i4_0_3]*n[i4_0_0] + m[i4_1_3]*n[i4_0_1] + m[i4_2_3]*n[i4_0_2] + m[i4_3_3]*n[i4_0_3],

		// Column 1
		m[i4_0_0]*n[i4_1_0] + m[i4_1_0]*n[i4_1_1] + m[i4_2_0]*n[i4_1_2] + m[i4_3_0]*n[i4_1_3],
		m[i4_0_1]*n[i4_1_0] + m[i4_1_1]*n[i4_1_1] + m[i4_2_1]*n[i4_1_2] + m[i4_3_1]*n[i4_1_3],
		m[i4_0_2]*n[i4_1_0] + m[i4_1_2]*n[i4_1_1] + m[i4_2_2]*n[i4_1_2] + m[i4_3_2]*n[i4_1_3],
		m[i4_0_3]*n[i4_1_0] + m[i4_1_3]*n[i4_1_1] + m[i4_2_3]*n[i4_1_2] + m[i4_3_3]*n[i4_1_3],

		// Column 2
		m[i4_0_0]*n[i4_2_0] + m[i4_1_0]*n[i4_2_1] + m[i4_2_0]*n[i4_2_2] + m[i4_3_0]*n[i4_2_3],
		m[i4_0_1]*n[i4_2_0] + m[i4_1_1]*n[i4_2_1] + m[i4_2_1]*n[i4_2_2] + m[i4_3_1]*n[i4_2_3],
		m[i4_0_2]*n[i4_2_0] + m[i4_1_2]*n[i4_2_1] + m[i4_2_2]*n[i4_2_2] + m[i4_3_2]*n[i4_2_3],
		m[i4_0_3]*n[i4_2_0] + m[i4_1_3]*n[i4_2_1] + m[i4_2_3]*n[i4_2_2] + m[i4_3_3]*n[i4_2_3],

		// Column 3
		m[i4_0_0]*n[i4_3_0] + m[i4_1_0]*n[i4_3_1] + m[i4_2_0]*n[i4_3_2] + m[i4_3_0]*n[i4_3_3],
		m[i4_0_1]*n[i4_3_0] + m[i4_1_1]*n[i4_3_1] + m[i4_2_1]*n[i4_3_2] + m[i4_3_1]*n[i4_3_3],
		m[i4_0_2]*n[i4_3_0] + m[i4_1_2]*n[i4_3_1] + m[i4_2_2]*n[i4_3_2] + m[i4_3_2]*n[i4_3_3],
		m[i4_0_3]*n[i4_3_0] + m[i4_1_3]*n[i4_3_1] + m[i4_2_3]*n[i4_3_2] + m[i4_3_3]*n[i4_3_3],
	}
	return m
}

// Matrix Indices
const (
	// 4x4 - x_y
	i4_0_0 = 0
	i4_0_1 = 1
	i4_0_2 = 2
	i4_0_3 = 3
	i4_1_0 = 4
	i4_1_1 = 5
	i4_1_2 = 6
	i4_1_3 = 7
	i4_2_0 = 8
	i4_2_1 = 9
	i4_2_2 = 10
	i4_2_3 = 11
	i4_3_0 = 12
	i4_3_1 = 13

	i4_3_2 = 14
	i4_3_3 = 15

	// 3x3 - x_y
	i3_0_0 = 0
	i3_0_1 = 1
	i3_0_2 = 2
	i3_1_0 = 3
	i3_1_1 = 4
	i3_1_2 = 5
	i3_2_0 = 6
	i3_2_1 = 7
	i3_2_2 = 8
)

func (b Box) Rect() Rect {
	return Rect{
		Min: Vec2{b.Min.X, b.Min.Y},
		Max: Vec2{b.Max.X, b.Max.Y},
	}
}

func (a Box) Union(b Box) Box {
	x1, _ := minMax(a.Min.X, b.Min.X)
	_, x2 := minMax(a.Max.X, b.Max.X)
	y1, _ := minMax(a.Min.Y, b.Min.Y)
	_, y2 := minMax(a.Max.Y, b.Max.Y)
	z1, _ := minMax(a.Min.Z, b.Min.Z)
	_, z2 := minMax(a.Max.Z, b.Max.Z)
	return Box{
		Min: Vec3{x1, y1, z1},
		Max: Vec3{x2, y2, z2},
	}
}

// TODO: This is the wrong input matrix type
func (b Box) Apply(mat Mat4) Box {
	return Box{
		Min: mat.Apply(b.Min),
		Max: mat.Apply(b.Max),
	}
}

func (m *Mat4) Apply(v Vec3) Vec3 {
	return Vec3{
		m[i4_0_0]*v.X + m[i4_1_0]*v.Y + m[i4_2_0]*v.Z + m[i4_3_0], // w = 1.0
		m[i4_0_1]*v.X + m[i4_1_1]*v.Y + m[i4_2_1]*v.Z + m[i4_3_1], // w = 1.0
		m[i4_0_2]*v.X + m[i4_1_2]*v.Y + m[i4_2_2]*v.Z + m[i4_3_2], // w = 1.0
	}
}

func (m *Mat3) Apply(v Vec2) Vec2 {
	return Vec2{
		m[i3_0_0]*v.X + m[i3_1_0]*v.Y + m[i3_2_0],
		m[i3_0_1]*v.X + m[i3_1_1]*v.Y + m[i3_2_1],
	}
}

// Note: Returns a new Mat4
func (m *Mat4) Inv() *Mat4 {
	retMat := Mat4(mgl64.Mat4(*m).Inv())
	return &retMat
}

func (m *Mat4) Transpose() *Mat4 {
	retMat := Mat4(mgl64.Mat4(*m).Transpose())
	return &retMat
}
