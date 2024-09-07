package geom

import (
	"fmt"
	"math"
)

// Line is a 2D line segment, between points A and B.
type Line struct {
	A, B *Vector2d
}

// L creates and returns a new Line.
func L(from, to *Vector2d) Line {
	return Line{
		A: NewVector2d(from.x, from.y),
		B: NewVector2d(to.x, to.y),
	}
}

// Bounds returns the lines bounding box.  This is in the form of a normalized Rect.
// func (l Line) Bounds() Rect {
// 	return R(l.A.X, l.A.Y, l.B.X, l.B.Y).Norm()
// }

// 返回线的中心点
func (l Line) Center() *Vector2d {
	temp := l.A.To(l.B)
	temp.Mul(0.5)
	return l.A.AddN(temp)
}

// 返回v点和线上最接近的点
func (l Line) Closest(v *Vector2d) *Vector2d {
	// between is a helper function which determines whether x is greater than min(a, b) and less than max(a, b)
	between := func(a, b, x float64) bool {
		min := math.Min(a, b)
		max := math.Max(a, b)
		return min < x && x < max
	}

	// Closest point will be on a line which perpendicular to this line.
	// If and only if the infinite perpendicular line intersects the segment.
	m, b := l.Formula()

	// Account for horizontal lines
	if m == 0 {
		x := v.x
		y := l.A.y

		// check if the X coordinate of v is on the line
		if between(l.A.x, l.B.x, v.x) {
			return NewVector2d(x, y)
		}

		// Otherwise get the closest endpoint
		if l.A.To(v).LengthSq() < l.B.To(v).LengthSq() {
			return l.A
		}
		return l.B
	}

	// Account for vertical lines
	if math.IsInf(math.Abs(m), 1) {
		x := l.A.x
		y := v.y

		// check if the Y coordinate of v is on the line
		if between(l.A.y, l.B.y, v.y) {
			return NewVector2d(x, y)
		}

		// Otherwise get the closest endpoint
		if l.A.To(v).LengthSq() < l.B.To(v).LengthSq() {
			return l.A
		}
		return l.B
	}

	perpendicularM := -1 / m
	perpendicularB := v.y - (perpendicularM * v.x)

	// Coordinates of intersect (of infinite lines)
	x := (perpendicularB - b) / (m - perpendicularM)
	y := m*x + b

	// Check if the point lies between the x and y bounds of the segment
	if !between(l.A.x, l.B.x, x) && !between(l.A.y, l.B.y, y) {
		// Not within bounding box
		toStart := v.To(l.A)
		toEnd := v.To(l.B)

		if toStart.Length() < toEnd.Length() {
			return l.A
		}
		return l.B
	}

	return NewVector2d(x, y)
}

// v点是否在线上
func (l Line) Contains(v *Vector2d) bool {
	return l.Closest(v).Eq(v)
}

// 返回直线方程y= mx + b 的m 和 b参数
func (l Line) Formula() (m, b float64) {
	// Account for horizontal lines
	if l.B.y == l.A.y {
		return 0, l.A.y
	}

	m = (l.B.y - l.A.y) / (l.B.x - l.A.x)
	b = l.A.y - (m * l.A.x)

	return m, b
}

// 返回两线段相交的交点 不相交则返回 false
func (l Line) Intersect(k Line) (*Vector2d, bool) {
	// Check if the lines are parallel
	lDir := l.A.To(l.B)
	kDir := k.A.To(k.B)
	if lDir.x == kDir.x && lDir.y == kDir.y {
		return NewVector2d(0, 0), false
	}

	// The lines intersect - but potentially not within the line segments.
	// Get the intersection point for the lines if they were infinitely long, check if the point exists on both of the
	// segments
	lm, lb := l.Formula()
	km, kb := k.Formula()

	// Account for vertical lines
	if math.IsInf(math.Abs(lm), 1) && math.IsInf(math.Abs(km), 1) {
		// Both vertical, therefore parallel
		return NewVector2d(0, 0), false
	}

	var x, y float64

	if math.IsInf(math.Abs(lm), 1) || math.IsInf(math.Abs(km), 1) {
		// One line is vertical
		intersectM := lm
		intersectB := lb
		verticalLine := k

		if math.IsInf(math.Abs(lm), 1) {
			intersectM = km
			intersectB = kb
			verticalLine = l
		}

		y = intersectM*verticalLine.A.x + intersectB
		x = verticalLine.A.x
	} else {
		// Coordinates of intersect
		x = (kb - lb) / (lm - km)
		y = lm*x + lb
	}

	if l.Contains(NewVector2d(x, y)) && k.Contains(NewVector2d(x, y)) {
		// The intersect point is on both line segments, they intersect.
		return NewVector2d(x, y), true
	}

	return NewVector2d(0, 0), false
}

// IntersectCircle will return the shortest Vector2d such that moving the Line by that Vector2d will cause the Line and Circle
// to no longer intesect.  If they do not intersect at all, this function will return a zero-vector.
// func (l Line) IntersectCircle(c Circle) Vector2d {
// 	// Get the point on the line closest to the center of the circle.
// 	closest := l.Closest(c.Center)
// 	cirToClosest := c.Center.To(closest)

// 	if cirToClosest.Len() >= c.Radius {
// 		return ZV
// 	}

// 	return cirToClosest.Scaled(cirToClosest.Len() - c.Radius)
// }

// IntersectRect will return the shortest Vector2d such that moving the Line by that Vector2d will cause  the Line and Rect to
// no longer intesect.  If they do not intersect at all, this function will return a zero-vector.
// func (l Line) IntersectRect(r Rect) Vector2d {
// 	// Check if either end of the line segment are within the rectangle
// 	if r.Contains(l.A) || r.Contains(l.B) {
// 		// Use the Rect.Intersect to get minimal return value
// 		rIntersect := l.Bounds().Intersect(r)
// 		if rIntersect.H() > rIntersect.W() {
// 			// Go vertical
// 			return V(0, rIntersect.H())
// 		}
// 		return V(rIntersect.W(), 0)
// 	}

// 	// Check if any of the rectangles' edges intersect with this line.
// 	for _, edge := range r.Edges() {
// 		if _, ok := l.Intersect(edge); ok {
// 			// Get the closest points on the line to each corner, where:
// 			//  - the point is contained by the rectangle
// 			//  - the point is not the corner itself
// 			corners := r.Vertices()
// 			var closest *Vector2d
// 			closestCorner := corners[0]
// 			for _, c := range corners {
// 				cc := l.Closest(c)
// 				if closest == nil || (closest.Len() > cc.Len() && r.Contains(cc)) {
// 					closest = &cc
// 					closestCorner = c
// 				}
// 			}

// 			return closest.To(closestCorner)
// 		}
// 	}

// 	// No intersect
// 	return ZV
// }

// 返回线段长度
func (l Line) Len() float64 {
	return l.A.To(l.B).Length()
}

// Moved will return a line moved by the delta Vector2d provided.
func (l Line) Moved(delta *Vector2d) Line {
	return Line{
		A: l.A.AddN(delta),
		B: l.B.AddN(delta),
	}
}

// 把线段旋转
func (l Line) Rotated(around *Vector2d, angle float64) Line {
	// Move the line so we can use `Vector2d.Rotated`
	lineShifted := l.Moved(around.MulN(-1))

	lineRotated := Line{
		A: lineShifted.A.RotateN(angle),
		B: lineShifted.B.RotateN(angle),
	}

	return lineRotated.Moved(around)
}

// Scaled will return the line scaled around the center point.
func (l Line) Scaled(scale float64) Line {
	return l.ScaledXY(l.Center(), scale)
}

// ScaledXY will return the line scaled around the Vector2d provided.
func (l Line) ScaledXY(around *Vector2d, scale float64) Line {
	toA := around.To(l.A)
	toA.Mul(scale)
	toB := around.To(l.B)
	toB.Mul(scale)

	return Line{
		A: around.AddN(toA),
		B: around.AddN(toB),
	}
}

func (l Line) String() string {
	return fmt.Sprintf("Line(%v, %v)", l.A, l.B)
}
