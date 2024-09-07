package utils

import "math"

const FLOAT_EQUALS_RANGE = 0.00001

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func IsFloat64Eq(x, y float64) bool {
	return math.Abs(x-y) <= FLOAT_EQUALS_RANGE
}

func IsFloat64Eg(x, y float64) bool {
	return IsFloat64Eq(x, y) || x > y
}

func IsFloat64El(x, y float64) bool {
	return IsFloat64Eq(x, y) || x < y
}

func IsFloat64Gt(x, y float64) bool {
	return !IsFloat64Eq(x, y) && x > y
}

func IsFloat64Lt(x, y float64) bool {
	return !IsFloat64Eq(x, y) && x < y
}
