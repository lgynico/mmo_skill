package geom

import (
	"fmt"
	"math"

	"github.com/lgynico/mmo_skill/utils"
)

type Vector2d struct {
	x, y, w float64
}

func NewVector2d(x, y float64) *Vector2d {
	return &Vector2d{
		x: x,
		y: y,
		w: 1.0,
	}
}

func NewVector2dByCell(cellX, cellY int) *Vector2d {
	return &Vector2d{
		x: float64(cellX) / 2.0,
		y: float64(cellY) / 2.0,
		w: 1.0,
	}
}

// func NewVector3d(x, y, w float64) *Vector2d {
// 	return &Vector2d{
// 		x: x, y: y,
// 		w: w,
// 	}
// }

func (v *Vector2d) GetX() float64 {
	return v.x
}

func (v *Vector2d) GetY() float64 {
	return v.y
}

func (v *Vector2d) SetX(x float64) {
	v.x = x
}

func (v *Vector2d) SetY(y float64) {
	v.y = y
}

func (v *Vector2d) Copy() *Vector2d {
	return NewVector2d(v.x, v.y)
}

func (v *Vector2d) Length() float64 {
	lenSq := v.LengthSq()
	return math.Sqrt(lenSq)
}

func (v *Vector2d) LengthSq() float64 {
	sqX := v.x * v.x
	sqY := v.y * v.y
	return sqX + sqY
}

func (v *Vector2d) Translate(tx, ty float64) {
	v.x += tx
	v.y += ty
}

func (v *Vector2d) TranslateN(tx, ty float64) *Vector2d {
	vec := NewVector2d(v.x, v.y)
	vec.Translate(tx, ty)
	return vec
}

func (v *Vector2d) Scale(sx, sy float64) {
	v.x *= sx
	v.y *= sy
}

func (v *Vector2d) ScaleN(sx, sy float64) *Vector2d {
	vec := NewVector2d(v.x, v.y)
	vec.Scale(sx, sy)
	return vec
}

// V1(x0, y0), V2(x1, y1)
// x0 = PcosA, y0 = PsinA
// x1 = Pcos(A + B) = PcosAcosB - PsinAsinB = x0cosB - y0sinB
// y1 = Psin(A + B) = PsinAsinB + PcosAcosB = y0sinB + x0cosB
func (v *Vector2d) Rotate(radian float64) {
	if math.IsNaN(radian) {
		return
	}
	sin := math.Sin(radian)
	cos := math.Cos(radian)

	x := (v.x * cos) - (v.y * sin)
	y := (v.x * sin) + (v.y * cos)

	v.x = x
	v.y = y
}

func (v *Vector2d) RotateN(radian float64) *Vector2d {
	vec := NewVector2d(v.x, v.y)
	vec.Rotate(radian)
	return vec
}

func (v *Vector2d) Shear(sx, sy float64) {
	x := v.x + (sx * v.y)
	y := v.y + (sy * v.x)

	v.x = x
	v.y = y
}

func (v *Vector2d) ShearN(sx, sy float64) *Vector2d {
	vec := NewVector2d(v.x, v.y)
	vec.Shear(sx, sy)
	return vec
}

func (v *Vector2d) Dot(v2 *Vector2d) float64 {
	return (v.x * v2.x) + (v.y * v2.y)
}

func (v *Vector2d) Cross(v2 *Vector2d) float64 {
	return (v.x * v2.y) - (v2.x * v.y)
}

func (v *Vector2d) Add(v2 *Vector2d) {
	v.x += v2.x
	v.y += v2.y
}

func (v *Vector2d) AddN(v2 *Vector2d) *Vector2d {
	vec := NewVector2d(v.x, v.y)
	vec.Add(v2)
	return vec
}

func (v *Vector2d) Sub(v2 *Vector2d) {
	v.x -= v2.x
	v.y -= v2.y
}

func (v *Vector2d) SubN(v2 *Vector2d) *Vector2d {
	vec := NewVector2d(v.x, v.y)
	vec.Sub(v2)
	return vec
}

func (v *Vector2d) Mul(scalar float64) {
	v.x *= scalar
	v.y *= scalar
}

func (v *Vector2d) MulN(scalar float64) *Vector2d {
	vec := NewVector2d(v.x, v.y)
	vec.Mul(scalar)
	return vec
}

func (v *Vector2d) Div(scalar float64) {
	if !utils.IsFloat64Eq(scalar, 0) {
		v.x /= scalar
		v.y /= scalar
	}
}

func (v *Vector2d) DivN(scalar float64) *Vector2d {
	vec := NewVector2d(v.x, v.y)
	vec.Div(scalar)
	return vec
}

func (v *Vector2d) Normalize() {
	len := v.Length()
	v.Div(len)
}

func (v *Vector2d) NormalizeN() *Vector2d {
	vec := NewVector2d(v.x, v.y)
	vec.Normalize()
	return vec
}

func (v *Vector2d) Inverse() *Vector2d {
	return NewVector2d(-v.x, -v.y)
}

func (v *Vector2d) Perp() *Vector2d {
	return NewVector2d(-v.y, v.x)
}

func (v *Vector2d) DistTo(v2 *Vector2d) float64 {
	distSq := v.DistSqTo(v2)
	return math.Sqrt(distSq)
}

func (v *Vector2d) DistSqTo(v2 *Vector2d) float64 {
	dx := v2.x - v.x
	dy := v2.y - v.y
	return (dx * dx) + (dy * dy)
}

func (v *Vector2d) Angle() float64 {
	return math.Atan2(v.y, v.x)
}

func (v *Vector2d) Project(v2 *Vector2d) *Vector2d {
	scalar := v.Dot(v2) / v2.LengthSq()
	x := v2.x * scalar
	y := v2.y * scalar
	return NewVector2d(x, y)
}

func (v *Vector2d) Trim(minX, minY, maxX, maxY float64) {
	if v.x > maxX {
		v.x = maxX
	}

	if v.y > maxY {
		v.y = maxY
	}

	if v.x < minX {
		v.x = minX
	}

	if v.y < minY {
		v.y = minY
	}

}

func (v *Vector2d) GetCellX() int {
	return int(v.x * 2)
}

func (v *Vector2d) GetCellY() int {
	return int(v.y * 2)
}

func (v *Vector2d) To(v2 *Vector2d) *Vector2d {
	return NewVector2d(
		v2.x-v.x,
		v2.y-v.y,
	)
}
func (v *Vector2d) Eq(v2 *Vector2d) bool {
	return utils.IsFloat64Eq(v.x, v2.x) && utils.IsFloat64Eq(v.y, v2.y)
}

func (v *Vector2d) String() string {
	return fmt.Sprintf("(%.2f, %.2f)", v.x, v.y)
}
