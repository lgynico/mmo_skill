package geom

import (
	"math"

	"github.com/lgynico/mmo_skill/utils"
)

func TwoShapeCollision(s1, s2 Shape) bool {
	if cir, ok := s1.(*Circle); ok {
		if cir2, ok := s2.(*Circle); ok {
			return TwoCircleCollision(cir, cir2)
		}
		if rect, ok := s2.(*Rect); ok {
			return CircleRectCollision(cir, rect)
		}
		if fan, ok := s2.(*Fan); ok {
			return FanCicleCollision(fan, cir)
		}
		return PointInCircle(s2.Pos(), cir)
	}
	if rect, ok := s1.(*Rect); ok {
		if rect2, ok := s2.(*Rect); ok {
			return TwoRectCollision(rect, rect2)
		}
		if cir, ok := s2.(*Circle); ok {
			return CircleRectCollision(cir, rect)
		}
		if fan, ok := s2.(*Fan); ok {
			return FanRectCollision(fan, rect)
		}
		return PointInRect(s2.Pos(), rect)
	}
	if fan, ok := s1.(*Fan); ok {
		if rect, ok := s2.(*Rect); ok {
			return FanRectCollision(fan, rect)
		}
		if cir, ok := s2.(*Circle); ok {
			return FanCicleCollision(fan, cir)
		}
	}
	return false
}

func TwoCircleCollision(c1, c2 *Circle) bool {
	return TwoCircleCollision2(c1.Pos(), c1.Radius(), c2.Pos(), c2.Radius())
}

// TODO 为啥不行
// func TwoRectCollision(r1, r2 *Rect) bool {
// 	a0, a1, a2 := r1.LeftTop(), r1.RightTop(), r1.RightBottom()
// 	n0 := a0.SubN(a1).DivN(2)
// 	n1 := a1.SubN(a2).DivN(2)
// 	ca := a0.AddN(a2).DivN(2)

// 	d0 := n0.Length()
// 	d1 := n1.Length()
// 	n1.Div(d1)
// 	n0.Div(d0)

// 	b0, b1, b2 := r2.LeftTop(), r2.RightTop(), r2.RightBottom()
// 	n2 := b0.SubN(b1).DivN(2)
// 	n3 := b1.SubN(b2).DivN(2)
// 	cb := b0.AddN(b2).DivN(2)

// 	d2 := n2.Length()
// 	d3 := n3.Length()
// 	n2.Div(d2)
// 	n3.Div(d3)

// 	c := ca.SubN(cb)

// 	da := d0
// 	db := d2 * math.Abs(n2.Dot(n0))
// 	db += d3 * math.Abs(n3.Dot(n0))

// 	if da+db < math.Abs(c.Dot(n0)) {
// 		return false

// 	}

// 	da = d1
// 	db = d2 * math.Abs(n2.Dot(n1))
// 	db += d3 * math.Abs(n3.Dot(n1))

// 	if da+db < math.Abs(c.Dot(n1)) {
// 		return false

// 	}

// 	da = d2
// 	db = d0 * math.Abs(n0.Dot(n2))
// 	db += d1 * math.Abs(n1.Dot(n2))

// 	if da+db < math.Abs(c.Dot(n2)) {
// 		return false

// 	}

// 	da = d3
// 	db = d0 * math.Abs(n0.Dot(n3))
// 	db += d1 * math.Abs(n1.Dot(n3))

// 	return da+db < math.Abs(c.Dot(n3))
// }

func TwoRectCollision(r1, r2 *Rect) bool {
	axes := make([]*Vector2d, 0, 4)

	a1 := r1.leftTop.SubN(r1.rightTop).Perp()
	a1.Normalize()
	axes = append(axes, a1)

	a2 := r1.rightTop.SubN(r2.rightBottom).Perp()
	a2.Normalize()
	axes = append(axes, a2)

	a3 := r2.leftTop.SubN(r2.rightTop).Perp()
	a3.Normalize()
	axes = append(axes, a3)

	a4 := r2.rightTop.SubN(r2.rightBottom).Perp()
	a4.Normalize()
	axes = append(axes, a4)

	for _, axi := range axes {
		max1, min1 := r1.Projection(axi)
		max2, min2 := r2.Projection(axi)
		if min1 > max2 || min2 > max1 {
			return false
		}
	}

	return true
}

func CircleRectCollision(cir *Circle, rect *Rect) bool {
	axes := make([]*Vector2d, 0, 3)

	a1 := rect.leftTop.SubN(rect.rightTop).Perp()
	a1.Normalize()
	axes = append(axes, a1)

	a2 := rect.rightTop.SubN(rect.rightBottom).Perp()
	a2.Normalize()
	axes = append(axes, a2)

	c := cir.Pos()
	a3 := rect.NearestVertexToPoint(c).SubN(c).Perp()
	a3.Normalize()
	axes = append(axes, a3)

	for _, axi := range axes {
		max1, min1 := rect.Projection(axi)
		max2, min2 := cir.Projection(axi)
		if min1 > max2 || min2 > max1 {
			return false
		}
	}

	return true
}

func FanCicleCollision(fan *Fan, cir *Circle) bool {
	// 两个圆相交
	bCircleCollision := TwoCircleCollision2(fan.Pos(), fan.Radius(), cir.Pos(), cir.Radius())
	if !bCircleCollision {
		return false
	}

	// 圆心在扇形夹角范围内
	a := fan.Pos()
	b := fan.GetP1()
	c := fan.GetP2()
	p := cir.Pos()
	ab := b.SubN(a)
	ac := c.SubN(a)
	ap := p.SubN(a)
	ab.Normalize()
	ac.Normalize()
	ap.Normalize()

	sitha1 := math.Acos(ap.Dot(ab))
	sitha2 := math.Acos(ap.Dot(ac))
	sitha := math.Acos(ab.Dot(ac))

	if utils.IsFloat64Eq(sitha1+sitha2, sitha) {
		return true
	}

	// 圆心到扇形两边的距离
	dist := PointDistToLine(p, a, b)
	if dist <= cir.Radius() {
		return true
	}

	dist = PointDistToLine(p, a, c)
	return dist <= cir.Radius()

}

func FanRectCollision(fan *Fan, rect *Rect) bool {
	return false
}

func PointInCircle(point *Vector2d, cir *Circle) bool {
	distSq := cir.Pos().DistSqTo(point)
	return distSq <= cir.Radius()*cir.Radius()

}

func PointInRect(point *Vector2d, rect *Rect) bool {

	return false
}

func TwoCircleCollision2(c1 *Vector2d, r1 float64, c2 *Vector2d, r2 float64) bool {
	if utils.IsFloat64Eq(r1, 0) || utils.IsFloat64Eq(r2, 0) {
		return false
	}
	distSq := c1.DistSqTo(c2)
	r := r1 + r2
	return distSq <= (r * r)
}

func PointDistToLine(p *Vector2d, a, b *Vector2d) float64 {
	ab := b.SubN(a)
	ap := p.SubN(a)
	dot := ap.Dot(ab)
	if dot <= 0 {
		return ap.Length()
	}

	abLen := ab.LengthSq()
	if dot >= abLen {
		bp := p.SubN(b)
		return bp.Length()
	}

	r := dot / abLen
	px := a.GetX() + (b.GetX()-a.GetX())*r
	py := a.GetY() + (b.GetY()-a.GetY())*r
	shadow := NewVector2d(px, py)

	return p.SubN(shadow).Length()
}
