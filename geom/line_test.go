package geom

import (
	"fmt"
	"testing"
)

func TestLineClosest(t *testing.T) {
	line := L(NewVector2d(1, 1), NewVector2d(5, 5))
	point := line.Closest(NewVector2d(5, 2))
	fmt.Println(point)

	// 3.5 3.5
}

func TestLineFormula(t *testing.T) {
	line := L(NewVector2d(1, 1), NewVector2d(5, 2))
	m, b := line.Formula()
	fmt.Println(m)
	fmt.Println(b)
}

func TestLineIntersect(t *testing.T) {
	line := L(NewVector2d(1, 1), NewVector2d(5, 5))
	line1 := L(NewVector2d(2, 3), NewVector2d(4, 1))
	point, ok := line.Intersect(line1)
	if ok {
		fmt.Println(point)
	}
}

func TestBresenham(t *testing.T) {
	fmt.Println(Bresenham(1, 1, 5000, 2000))

}
