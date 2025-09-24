package spatial

import (
	"fmt"
	"math"
	"testing"

	"github.com/unitoftime/flow/glm"
)

type intersectTest struct {
	a, b     Shape
	expected bool
}

func TestIntersectAABB(t *testing.T) {
	baseRect5 := glm.CR(5)
	baseRect10 := glm.CR(10)

	tests := []intersectTest{
		{
			AABB(glm.R(10, 10, 20, 20)),
			AABB(glm.R(15, 15, 25, 25)),
			true,
		},
		{
			AABB(baseRect5.WithCenter(glm.Vec2{})),
			AABB(baseRect5.WithCenter(glm.Vec2{})),
			true,
		},
		{
			AABB(baseRect5.WithCenter(glm.Vec2{})),
			AABB(baseRect5.WithCenter(glm.Vec2{5, 2})),
			true,
		},
		{
			AABB(baseRect5.WithCenter(glm.Vec2{})),
			AABB(baseRect5.WithCenter(glm.Vec2{50, 20})),
			false,
		},
		{
			AABB(baseRect5.WithCenter(glm.Vec2{})),
			AABB(baseRect5.WithCenter(glm.Vec2{10, 10})),
			true,
		},
		{
			AABB(baseRect5.WithCenter(glm.Vec2{})),
			AABB(baseRect5.WithCenter(glm.Vec2{0, 10})),
			true,
		},
		{
			AABB(baseRect5.WithCenter(glm.Vec2{})),
			AABB(baseRect5.WithCenter(glm.Vec2{0, 11})),
			false,
		},
		{
			AABB(baseRect5.WithCenter(glm.Vec2{})),
			AABB(baseRect10.WithCenter(glm.Vec2{})),
			true,
		},
		{
			AABB(baseRect5.WithCenter(glm.Vec2{0, 5})),
			AABB(baseRect10.WithCenter(glm.Vec2{})),
			true,
		},
		{
			AABB(baseRect5.WithCenter(glm.Vec2{0, 10})),
			AABB(baseRect10.WithCenter(glm.Vec2{})),
			true,
		},
		{
			AABB(baseRect5.WithCenter(glm.Vec2{0, 15})),
			AABB(baseRect10.WithCenter(glm.Vec2{})),
			true,
		},
		{
			// AABB(glm.R(4, -4, 12, 4.0)),
			AABB(glm.R(3.6799, -3.9999, 11.6799, 4.0)),
			AABB(glm.R(-12, -28, 12, 4)),
			true,
		},
		{
			//Min:{X:-4.000000000000001 Y:29.791999999999994} Max:{X:3.999999999999999 Y:37.791999999999994}
			//Min:{X:4 Y:20} Max:{X:28 Y:44}
			AABB(glm.R(-4.000000000000001, 29.791999999999994, 3.999999999999999, 37.791999999999994)),
			AABB(glm.R(4, 20, 28, 44)),
			false,
		},
	}

	for i := range tests {
		fmt.Println("Test:", i)
		a := tests[i].a
		b := tests[i].b

		expected := tests[i].expected
		actual := a.Intersects(b)
		if expected != actual {
			t.Errorf("Failed. Actual(%v) Expected(%v).\n  A: %+v\n  B: %+v\n", actual, expected, a, b)
		}
	}
}

func TestIntersectRect(t *testing.T) {
	baseRect5 := glm.CR(5)
	// baseRect10 := glm.CR(10)

	tests := []intersectTest{
		{
			Rect(baseRect5, *glm.IM4().Rotate(math.Pi/4, glm.Vec3{0, 0, 1})),
			Rect(baseRect5, *glm.IM4().Translate(0, 0, 0)),
			true,
		},
		{
			Rect(baseRect5, *glm.IM4().Rotate(math.Pi/4, glm.Vec3{0, 0, 1})),
			Rect(baseRect5, *glm.IM4().Translate(5, 0, 0)),
			true,
		},
		{
			// Rotated 45 degrees and nudged inward
			Rect(baseRect5, *glm.IM4().Rotate(math.Pi/4, glm.Vec3{0, 0, 1})),
			Rect(baseRect5, *glm.IM4().Translate(5+7.07106781187-0.001, 0, 0)),
			true,
		},

		{
			// Rotated 45 degrees and nudged outward
			Rect(baseRect5, *glm.IM4().Rotate(math.Pi/4, glm.Vec3{0, 0, 1})),
			Rect(baseRect5, *glm.IM4().Translate(5+7.07106781187+0.001, 0, 0)),
			false,
		},

		{
			// Rotated 45 degrees and nudged inward
			Rect(baseRect5, *glm.IM4().Rotate(math.Pi/4, glm.Vec3{0, 0, 1}).Translate(10, 10, 0)),
			Rect(baseRect5, *glm.IM4().Translate(5+7.07106781187-0.001, 0, 0).Translate(10, 10, 0)),
			true,
		},
		{
			// Rotated 45 degrees and both translated
			Rect(baseRect5, *glm.IM4().Rotate(math.Pi/4, glm.Vec3{0, 0, 1}).Translate(10, 10, 0)),
			Rect(baseRect5, *glm.IM4().Translate(5+7.07106781187+0.001, 0, 0).Translate(10, 10, 0)),
			false,
		},
	}

	for i := range tests {
		fmt.Println("Test:", i)
		a := tests[i].a
		b := tests[i].b

		expected := tests[i].expected
		actual := a.Intersects(b)
		if expected != actual {
			t.Errorf("Failed. Actual(%v) Expected(%v).\n  A: %+v\n  B: %+v\n", actual, expected, a, b)
		}
	}
}

func TestIntersectPositionToIndex(t *testing.T) {
	h := NewPositionHasher([2]int{16, 16})
	idx := h.PositionToIndex(glm.Vec2{-1, -1})
	r := h.IndexToRect(idx)
	fmt.Println("idx:", idx)
	fmt.Println("rect:", r)
}
