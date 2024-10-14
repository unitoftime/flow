package glm

import (
	"fmt"
	"math"
	"testing"
)

var rad2Deg = 180.0 / math.Pi

func angleTest(a, b Vec2) {
	fmt.Println(a, b, Angle(a, b)*rad2Deg)
}

func TestAngle(t *testing.T) {
	angleTest(Vec2{0, 1}, Vec2{1, 0})
	angleTest(Vec2{1, 0}, Vec2{0, 1})
	angleTest(Vec2{1, 1}, Vec2{1, 0})
	angleTest(Vec2{1, 1}, Vec2{1, 1})
	angleTest(Vec2{2, 1}, Vec2{1, 1})
	angleTest(Vec2{1, 1}, Vec2{2, 1})
	angleTest(Vec2{1, 1}, Vec2{-2, 1})
	angleTest(Vec2{1, 1}, Vec2{-2, -1})
}

func TestDotProduct(t *testing.T) {
	a := Vec3{1, 2, 3}
	b := Vec3{1, 5, 7}
	dot := a.Dot(b)
	if dot != 32 {
		panic("Should be 32")
	}
}

// func TestAngle(t *testing.T) {
// 	a := Vec3{2, -4, -1}
// 	b := Vec3{1, 5, 2}
// 	angle := a.Angle(b)
// 	// Should be: 2.4928086
// 	fmt.Println(angle)
// }

func TestRotate(t *testing.T) {
	a := Vec3{1, 0, 0}
	b := a.Rotate2D(3.14159 / 2)
	fmt.Println(b)
}

func TestProjectUnproject(t *testing.T) {
	data := []Vec3{
		Vec3{},
		Vec3{1, 2, 3},
	}

	// Arbitrary matrix
	mat := Mat4Ident
	mat.Translate(-100.2, -200.8, 0).
		Scale(0.7, 0.4, 1.0).
		Translate(200.9, 300.1, 0)

	for i := range data {
		intermediate := mat.Apply(data[i])
		result := mat.Inv().Apply(intermediate)
		fmt.Println("In:  ", data[i])
		fmt.Println("Out: ", result)
	}
}

// func TestSettingGlobalVariable(t *testing.T) {
// 	original := Mat4Ident
// 	original[2] = 123
// 	fmt.Println(Mat4Ident)
// 	Mat4Ident[0] += 10
// 	fmt.Println(Mat4Ident)
// 	Mat4Ident[0] -= 10
// 	fmt.Println(Mat4Ident)
// 	Mat4Ident = Mat4Ident
// 	fmt.Println(Mat4Ident)
// }
