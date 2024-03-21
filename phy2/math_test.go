package phy2

import (
	"fmt"
	"math"
	"testing"
)

var rad2Deg = 180.0 / math.Pi

func angleTest(a, b Vec) {
	fmt.Println(a, b, Angle(a, b) * rad2Deg)
}

func TestAngle(t *testing.T) {
	angleTest(Vec{0, 1}, Vec{1, 0})
	angleTest(Vec{1, 0}, Vec{0, 1})
	angleTest(Vec{1, 1}, Vec{1, 0})
	angleTest(Vec{1, 1}, Vec{1, 1})
	angleTest(Vec{2, 1}, Vec{1, 1})
	angleTest(Vec{1, 1}, Vec{2, 1})
	angleTest(Vec{1, 1}, Vec{-2, 1})
	angleTest(Vec{1, 1}, Vec{-2, -1})
}
