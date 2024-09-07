package geom

import (
	"fmt"
	"math"
	"testing"
)

func TestDot(t *testing.T) {
	v1 := NewVector2d(1, 0)
	v2 := NewVector2d(-1, 1)
	dot := v1.Dot(v2)
	fmt.Println(dot)
	fmt.Println(math.Acos(dot))

	f := math.Atan2(v2.y, v2.x)
	fmt.Println(f)
	fmt.Println(math.Pi * 3 / 4)
}
