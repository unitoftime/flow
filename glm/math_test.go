package glm

import (
	"fmt"
	"math"
	"testing"
)

var rad2Deg = 180.0 / math.Pi

func angleTest(a, b Vec2) {
	fmt.Println(a, b, Angle(a, b) * rad2Deg)
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
