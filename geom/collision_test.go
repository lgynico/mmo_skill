package geom

import (
	"math"
	"testing"
)

func TestPointDistToLine(t *testing.T) {
	type testCase struct {
		a, b *Vector2d
		p    *Vector2d
		want float64
	}
	tcs := map[string]testCase{
		"锐角":  {NewVector2d(0, 0), NewVector2d(2, 0), NewVector2d(1, 1), 1},
		"钝角1": {NewVector2d(0, 0), NewVector2d(2, 0), NewVector2d(3, 3), math.Sqrt(10)},
		"钝角2": {NewVector2d(0, 0), NewVector2d(2, 0), NewVector2d(-1, 1), math.Sqrt(2)},
		"直角":  {NewVector2d(0, 0), NewVector2d(2, 0), NewVector2d(0, 1), 1},
		"同向1": {NewVector2d(0, 0), NewVector2d(2, 0), NewVector2d(1, 0), 0},
		"同向2": {NewVector2d(0, 0), NewVector2d(2, 0), NewVector2d(3, 0), 1},
		"反向":  {NewVector2d(0, 0), NewVector2d(2, 0), NewVector2d(-1, 0), 1},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			got := PointDistToLine(tc.p, tc.a, tc.b)
			if math.Dim(tc.want, got) > 0.00001 {
				t.Fatal(tc, got)
			}
		})
	}
}

func TestFanCircleCollision(t *testing.T) {
	type testCase struct {
		fan  *Fan
		cir  *Circle
		want bool
	}
	tcs := map[string]testCase{
		"outside":            {createFan(1, 90), createCircle(3, 0, 1), false},
		"tangency":           {createFan(1, 90), createCircle(2, 0, 1), true},
		"insideCollision":    {createFan(1, 90), createCircle(1, 0, 0.5), true},
		"insideNotCollision": {createFan(1, 90), createCircle(-0.5, 0, 0.25), false},
		"tangenPoint":        {createFan(math.Sqrt(2), 90), createCircle(1, 2, 1), true},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			got := FanCicleCollision(tc.fan, tc.cir)
			if tc.want != got {
				t.Fatal(tc, got)
			}
		})
	}
}

func TestCircleRectCollision(t *testing.T) {
	type testCase struct {
		cir  *Circle
		rect *Rect
		want bool
	}
	tcs := map[string]testCase{
		"1": {createCircle(0, 0, 20), createRect(0, 0, 20, 10, 1, 0, true), true},
		"2": {createCircle(0, 0, 1), createRect(2, 0, 2, 2, 1, 0, true), false},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			got := CircleRectCollision(tc.cir, tc.rect)
			if tc.want != got {
				t.Fatal(tc, got)
			}
		})
	}
}

func TestTwoRectCollision(t *testing.T) {
	type testCase struct {
		rect0 *Rect
		rect1 *Rect
		want  bool
	}
	tcs := map[string]testCase{
		"1": {createRect(3150, 2050, 500, 500, 1, 0, false), createRect(2550, 2850, 500, 500, 1, 0, false), false},
		"2": {createRect(0, 0, 10, 10, 1, 0, false), createRect(5, 5, 10, 10, 1, 1, false), true},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			got := TwoRectCollision(tc.rect0, tc.rect1)
			if tc.want != got {
				t.Fatal(tc, got)
			}
		})
	}
}

func createFan(r, deg float64) *Fan {
	base := &BaseShape{
		pos:     NewVector2d(0, 0),
		heading: NewVector2d(1, 0),
	}

	return newFan(base, r, deg)
}

func createCircle(x, y, r float64) *Circle {
	base := BaseShape{
		pos:     NewVector2d(x, y),
		heading: NewVector2d(0, 0),
	}
	return &Circle{
		BaseShape: &base,
		radius:    r,
	}
}

func createRect(x, y float64, width, height float64, hx, hy float64, offset bool) *Rect {
	base := &BaseShape{
		pos:     NewVector2d(0, 0),
		heading: NewVector2d(1, 0),
	}
	rect := newRect(base, width, height, offset)
	rect.Translate(NewVector2d(x, y))
	rect.SetHeading(NewVector2d(hx, hy))
	return rect
}
