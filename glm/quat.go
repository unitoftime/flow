package glm

import (
	"github.com/go-gl/mathgl/mgl64"
)

type Quat struct {
	// x, y, z, w float64
	inner mgl64.Quat
}

func IQuat() Quat {
	return Quat{mgl64.QuatIdent()}
}

func QuatRotate(angle float64, axis Vec3) Quat {
	quat := mgl64.QuatRotate(angle, mgl64.Vec3{axis.X, axis.Y, axis.Z})

	return Quat{quat}
}

func QuatZ(angle float64) Quat {
	return QuatRotate(angle, Vec3{0, 0, 1})
}

func (q *Quat) Equals(q2 Quat) bool {
	return q.inner.OrientationEqual(q2.inner)
}

func (q *Quat) RotateQuat(rQuat Quat) *Quat {
	q.inner = q.inner.Mul(rQuat.inner)
	return q
}
func (q *Quat) RotateX(angle float64) *Quat {
	return q.Rotate(angle, Vec3{1, 0, 0})
}
func (q *Quat) RotateY(angle float64) *Quat {
	return q.Rotate(angle, Vec3{0, 1, 0})
}
func (q *Quat) RotateZ(angle float64) *Quat {
	return q.Rotate(angle, Vec3{0, 0, 1})
}

func (q *Quat) Rotate(angle float64, axis Vec3) *Quat {
	rQuat := QuatRotate(angle, axis)

	return q.RotateQuat(rQuat)
}

func (q *Quat) Mat4() Mat4 {
	return Mat4(q.inner.Mat4())
}
