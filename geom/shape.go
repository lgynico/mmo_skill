package geom

import (
	"fmt"
	"math"

	"github.com/lgynico/mmo_skill/utils"
)

type ShapeType float64

const (
	SHAPE_CIRCLE ShapeType = iota + 1
	SHAPE_RECT
	SHAPE_FAN
)

type Shape interface {
	Translate(*Vector2d)
	Pos() *Vector2d
	SetHeading(*Vector2d)
	Radius() float64
	Projection(*Vector2d) (float64, float64)
	Transform(matrix *Matrix3d)
	IsPointInside(*Vector2d) bool
}

type BaseShape struct {
	pos     *Vector2d
	heading *Vector2d
}

func (s *BaseShape) Translate(v *Vector2d) {
	s.pos.Translate(v.x, v.y)
}

func (s *BaseShape) Pos() *Vector2d {
	return s.pos
}

func (s *BaseShape) SetHeading(heading *Vector2d) {
	s.heading = heading.NormalizeN()
}

func (s *BaseShape) Projection(axi *Vector2d) (max, min float64) {
	return 0, 0
}

func (s *BaseShape) Transform(mat *Matrix3d) {
	s.pos = mat.MulV(NewVector2d(0, 0))
}

func (s *BaseShape) Radius() float64 {
	return 0
}

func (s *BaseShape) IsPointInside(p *Vector2d) bool {
	return s.Pos().Eq(p)
}

// ===========================================================
// 矩形
// ===========================================================
type Rect struct {
	*BaseShape
	vertexes    [4]*Vector2d
	leftTop     *Vector2d
	rightTop    *Vector2d
	rightBottom *Vector2d
	leftBottom  *Vector2d
}

func (r *Rect) Translate(v *Vector2d) {
	r.pos.Translate(v.x, v.y)
	r.leftTop.Translate(v.x, v.y)
	r.rightTop.Translate(v.x, v.y)
	r.rightBottom.Translate(v.x, v.y)
	r.leftBottom.Translate(v.x, v.y)
}

func (r *Rect) SetHeading(heading *Vector2d) {
	oldHeading := r.heading
	r.BaseShape.SetHeading(heading)

	// log.Printf("oldHeading: %v, newHeading: %v\n", oldHeading, heading)

	dot := heading.Dot(oldHeading)
	radian := math.Acos(dot)
	// if utils.IsFloat64Lt(heading.Y, 0) {
	// 	radian += math.Pi
	// }

	r.leftTop.Rotate(radian)
	r.rightTop.Rotate(radian)
	r.rightBottom.Rotate(radian)
	r.leftBottom.Rotate(radian)
}

func (r *Rect) Transform(mat *Matrix3d) {
	r.BaseShape.Transform(mat)
	r.leftTop = mat.MulV(r.vertexes[0])
	r.rightTop = mat.MulV(r.vertexes[1])
	r.rightBottom = mat.MulV(r.vertexes[2])
	r.leftBottom = mat.MulV(r.vertexes[3])
}

func (s *Rect) LeftTop() *Vector2d {
	return s.leftTop
}

func (s *Rect) LeftBottom() *Vector2d {
	return s.leftBottom
}

func (s *Rect) RightBottom() *Vector2d {
	return s.rightBottom
}

func (s *Rect) RightTop() *Vector2d {
	return s.rightTop
}

func (s *Rect) NearestVertexToPoint(p *Vector2d) *Vector2d {
	v := s.leftTop
	minDist := s.leftTop.DistSqTo(p)

	t := s.rightTop
	dist := t.DistSqTo(p)
	if dist < minDist {
		minDist = dist
		v = t
	}

	t = s.rightBottom
	dist = t.DistSqTo(p)
	if dist < minDist {
		minDist = dist
		v = t
	}

	t = s.leftBottom
	dist = t.DistSqTo(p)
	if dist < minDist {
		v = t
	}

	return v
}

func (r *Rect) Projection(axi *Vector2d) (max, min float64) {
	min = axi.Dot(r.leftTop)
	max = min

	dot := axi.Dot(r.rightTop)
	max = math.Max(max, dot)
	min = math.Min(min, dot)

	dot = axi.Dot(r.rightBottom)
	max = math.Max(max, dot)
	min = math.Min(min, dot)

	dot = axi.Dot(r.leftBottom)
	max = math.Max(max, dot)
	min = math.Min(min, dot)

	return
}

func getCross(a *Vector2d, b *Vector2d, p *Vector2d) float64 {
	v1 := NewVector2d(b.x-a.x, b.y-a.y)
	v2 := NewVector2d(p.x-a.x, p.y-a.y)
	return v1.Cross(v2)
}

func (r *Rect) IsPointInside(p *Vector2d) bool {
	cross1 := getCross(r.leftTop, r.rightTop, p)
	cross2 := getCross(r.rightBottom, r.leftBottom, p)

	cross3 := getCross(r.leftBottom, r.leftTop, p)
	cross4 := getCross(r.rightTop, r.rightBottom, p)
	return cross1*cross2 >= 0 && cross3*cross4 >= 0
}

func (r *Rect) Radius() float64 {
	v := r.leftBottom.AddN(r.rightTop)
	v.Div(2)
	return v.Length()
}

func (r *Rect) String() string {
	return fmt.Sprintf("rect[pos: %v, leftTop: %v, rightTop: %v, rightBottom: %v, leftBottom: %v]",
		r.pos, r.leftTop, r.rightTop, r.rightBottom, r.leftBottom)
	// return fmt.Sprintf("rect[\n%v\n%v\n%v\n%v\n%v\n]", r.pos, r.leftTop, r.rightTop, r.rightBottom, r.leftBottom)
}

// ===========================================================
// 圆形
// ===========================================================
type Circle struct {
	*BaseShape
	radius float64
}

func (c *Circle) Translate(v *Vector2d) {
	c.pos.Translate(v.x, v.y)
}

func (c *Circle) Radius() float64 {
	return c.radius
}

func (c *Circle) Projection(axi *Vector2d) (max, min float64) {
	dot := c.pos.Dot(axi)
	max = dot + c.radius
	min = dot - c.radius
	return
}

func (c *Circle) IsPointInside(p *Vector2d) bool {
	distSq := c.Pos().DistSqTo(p)
	return distSq <= c.Radius()*c.Radius()
}

func (c *Circle) String() string {
	return fmt.Sprintf("circle[pos:%v, radius:%.2f]", c.pos, c.radius)
}

// ===========================================================
// 扇形
// ===========================================================
type Fan struct {
	*BaseShape
	radius   float64
	degree   float64
	p1       *Vector2d
	p2       *Vector2d
	vertexes [2]*Vector2d
}

func (f *Fan) Translate(v *Vector2d) {
	f.pos.Translate(v.x, v.y)
	f.p1.Translate(v.x, v.y)
	f.p2.Translate(v.x, v.y)
}

func (f *Fan) SetHeading(heading *Vector2d) {
	oldHeading := f.heading
	f.BaseShape.SetHeading(heading)

	dot := heading.Dot(oldHeading)
	radian := math.Acos(dot)
	// if utils.IsFloat64Lt(heading.Y, 0) {
	// 	radian += math.Pi
	// }

	f.p1.Rotate(radian)
	f.p2.Rotate(radian)
}

func (f *Fan) Transform(mat *Matrix3d) {
	f.BaseShape.Transform(mat)
	f.p1 = mat.MulV(f.vertexes[0])
	f.p2 = mat.MulV(f.vertexes[1])
}

func (f *Fan) Radius() float64 {
	return f.radius
}

func (f *Fan) Degree() float64 {
	return f.degree
}

func (f *Fan) Radian() float64 {
	deg := float32(f.degree)
	return float64(deg / math.Pi * 180)
}

func (f *Fan) GetP1() *Vector2d {
	return f.p1
}

func (f *Fan) GetP2() *Vector2d {
	return f.p2
}

func (f *Fan) IsPointInside(p *Vector2d) bool {
	//判断距离
	distSq := f.Pos().DistSqTo(p)
	flag := distSq <= f.Radius()*f.Radius()
	if !flag {
		return false
	}

	// 圆心在扇形夹角范围内
	a := f.Pos()
	b := f.GetP1()
	c := f.GetP2()

	ab := b.SubN(a)
	ac := c.SubN(a)
	ap := p.SubN(a)
	ab.Normalize()
	ac.Normalize()
	ap.Normalize()

	sitha1 := math.Acos(ap.Dot(ab))
	sitha2 := math.Acos(ap.Dot(ac))
	sitha := math.Acos(ab.Dot(ac))

	return utils.IsFloat64Eq(sitha1+sitha2, sitha)
}

func (f *Fan) String() string {
	return fmt.Sprintf("radius: %.2f, degree: %.2f, pos: %v, p1:%v, p2: %v\n", f.radius, f.degree, f.pos, f.p1, f.p2)
}

func NewCircle(radius float64) *Circle {
	cir := NewShape(SHAPE_CIRCLE, radius)
	return cir.(*Circle)
}

func NewShapeInt(shapeType ShapeType, params ...int) Shape {
	ps := make([]float64, 0, len(params))
	for _, v := range params {
		ps = append(ps, float64(v))
	}
	return NewShape(shapeType, ps...)
}

func NewShape(shapeType ShapeType, params ...float64) Shape {
	baseShape := &BaseShape{
		pos:     NewVector2d(0, 0),
		heading: NewVector2d(1, 0),
	}
	var shape Shape
	switch shapeType {
	case SHAPE_RECT:
		width := params[0]
		height := params[1]
		bOffset := false
		if len(params) >= 3 {
			bOffset = params[2] == 1
		}
		shape = newRect(baseShape, width, height, bOffset)
	case SHAPE_CIRCLE:
		shape = &Circle{
			BaseShape: baseShape,
			radius:    params[0],
		}
	// case SHAPE_POINT:
	// 	shape = &Circle{
	// 		BaseShape: baseShape,
	// 		radius:    0,
	// 	}
	case SHAPE_FAN:
		radius := params[0]
		degree := params[1]
		shape = newFan(baseShape, radius, degree)
	default:
		shape = baseShape
	}

	return shape
}

func newRect(base *BaseShape, width, height float64, offset bool) *Rect {
	var (
		left, right float64
		top         = height / 2
		bottom      = -top
	)

	if offset {
		left = 0
		right = width
	} else {
		left = -width / 2
		right = -left
	}

	leftTop := NewVector2d(left, top)
	rightTop := NewVector2d(right, top)
	rightBottom := NewVector2d(right, bottom)
	leftBottom := NewVector2d(left, bottom)

	return &Rect{
		BaseShape: base,
		vertexes: [4]*Vector2d{
			leftTop.Copy(), rightTop.Copy(), rightBottom.Copy(), leftBottom.Copy(),
		},
		leftTop:     leftTop,
		rightTop:    rightTop,
		rightBottom: rightBottom,
		leftBottom:  leftBottom,
	}
}

func newFan(base *BaseShape, radius, degree float64) *Fan {
	c := base.Pos()
	deg := float64(degree)

	a := deg / 2.0 * math.Pi / 180
	x := c.x + radius*float64(math.Cos(a))
	y := c.y + radius*float64(math.Sin(a))

	p1 := NewVector2d(x, y)
	p2 := NewVector2d(x, -y)

	return &Fan{
		BaseShape: base,
		radius:    radius,
		degree:    degree,
		p1:        p1,
		p2:        p2,
		vertexes: [2]*Vector2d{
			p1.Copy(), p2.Copy(),
		},
	}
}
