package geom

import (
	"fmt"
	"testing"
)

func TestPointInRect(t *testing.T) {
	rect := NewShape(SHAPE_RECT, 100, 100, 0)
	p := NewVector2d(50, 50)

	ret := rect.(*Rect).IsPointInside(p)
	fmt.Println(rect.(*Rect).String())
	fmt.Println(ret)
}
