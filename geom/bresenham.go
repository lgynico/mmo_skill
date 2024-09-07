package geom

const (
	MAX_FOR_CNT = 100
)

type PointCoordinate struct {
	X, Y int32
}

// 求解两点经过的格子坐标 限制最大循环次数
func Bresenham(x1, y1, x2, y2 int32) []PointCoordinate {
	var dx, dy, e, slope int32
	ret := make([]PointCoordinate, 0)
	if x1 > x2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}
	dx, dy = x2-x1, y2-y1
	if dy < 0 {
		dy = -dy
	}

	foreachCnt := 0

	switch {
	case x1 == x2 && y1 == y2:
		ret = append(ret, PointCoordinate{x1, y1})
	case y1 == y2:
		for ; dx != 0; dx-- {
			ret = append(ret, PointCoordinate{x1, y1})
			x1++

			foreachCnt++
			if foreachCnt > MAX_FOR_CNT {
				break
			}
		}
		ret = append(ret, PointCoordinate{x1, y1})
	case x1 == x2:
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for ; dy != 0; dy-- {
			ret = append(ret, PointCoordinate{x1, y1})
			y1++

			foreachCnt++
			if foreachCnt > MAX_FOR_CNT {
				break
			}
		}
		ret = append(ret, PointCoordinate{x1, y1})
	case dx == dy:
		if y1 < y2 {
			for ; dx != 0; dx-- {
				ret = append(ret, PointCoordinate{x1, y1})
				x1++
				y1++

				foreachCnt++
				if foreachCnt > MAX_FOR_CNT {
					break
				}
			}
		} else {
			for ; dx != 0; dx-- {
				ret = append(ret, PointCoordinate{x1, y1})
				x1++
				y1--

				foreachCnt++
				if foreachCnt > MAX_FOR_CNT {
					break
				}
			}
		}
		ret = append(ret, PointCoordinate{x1, y1})
	case dx > dy:
		if y1 < y2 {
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				ret = append(ret, PointCoordinate{x1, y1})
				x1++
				e -= dy
				if e < 0 {
					y1++
					e += slope
				}

				foreachCnt++
				if foreachCnt > MAX_FOR_CNT {
					break
				}
			}
		} else {
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				ret = append(ret, PointCoordinate{x1, y1})
				x1++
				e -= dy
				if e < 0 {
					y1--
					e += slope
				}

				foreachCnt++
				if foreachCnt > MAX_FOR_CNT {
					break
				}
			}
		}
		ret = append(ret, PointCoordinate{x2, y2})
	default:
		if y1 < y2 {
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				ret = append(ret, PointCoordinate{x1, y1})
				y1++
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}

				foreachCnt++
				if foreachCnt > MAX_FOR_CNT {
					break
				}
			}
		} else {
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				ret = append(ret, PointCoordinate{x1, y1})
				y1--
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}

				foreachCnt++
				if foreachCnt > MAX_FOR_CNT {
					break
				}
			}
		}
		ret = append(ret, PointCoordinate{x2, y2})
	}
	return ret
}
