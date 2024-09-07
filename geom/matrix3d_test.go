package geom

import (
	"fmt"
	"math"
	"testing"
)

func TestMatrixRotate(t *testing.T) {
	v := NewVector2d(1, 0)
	mat := MatrixRotate(math.Pi * 11 / 6)
	v1 := mat.MulV(v)
	fmt.Println(v1)

	mat = mat.Mul(MatrixTranslate(10, 10))
	r := createRect(0, 0, 10, 10, 1, 0, true)
	fmt.Println(r)
	r.Transform(mat)
	fmt.Println(r)
}

func TestTranslateV(t *testing.T) {
	speed := 1000.0
	dist := speed * float64(1000) / 1000

	heading := NewVector2d(1, 1).NormalizeN()
	pos := NewVector2d(1, 1)
	v := heading.MulN(dist)

	mat := MatrixTranslateV(v)
	pos = mat.MulV(pos)

	fmt.Println(pos)
}
