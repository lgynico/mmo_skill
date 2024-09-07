package geom

import "math"

type Matrix3d struct {
	mat [3][3]float64
}

func IdentityMatrix() *Matrix3d {
	return &Matrix3d{
		mat: [3][3]float64{
			{1, 0, 0},
			{0, 1, 0},
			{0, 0, 1},
		},
	}
}

func ZeroMatrix() *Matrix3d {
	return &Matrix3d{
		mat: [3][3]float64{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
	}
}

func (m *Matrix3d) Add(m1 Matrix3d) *Matrix3d {
	return &Matrix3d{
		mat: [3][3]float64{
			{m.mat[0][0] + m1.mat[0][0], m.mat[0][1] + m1.mat[0][1], m.mat[0][2] + m1.mat[0][2]},
			{m.mat[1][0] + m1.mat[1][0], m.mat[1][1] + m1.mat[1][1], m.mat[1][2] + m1.mat[1][2]},
			{m.mat[2][0] + m1.mat[2][0], m.mat[2][1] + m1.mat[2][1], m.mat[2][2] + m1.mat[2][2]},
		},
	}
}

func (m *Matrix3d) Sub(m1 Matrix3d) *Matrix3d {
	return &Matrix3d{
		mat: [3][3]float64{
			{m.mat[0][0] - m1.mat[0][0], m.mat[0][1] - m1.mat[0][1], m.mat[0][2] - m1.mat[0][2]},
			{m.mat[1][0] - m1.mat[1][0], m.mat[1][1] - m1.mat[1][1], m.mat[1][2] - m1.mat[1][2]},
			{m.mat[2][0] - m1.mat[2][0], m.mat[2][1] - m1.mat[2][1], m.mat[2][2] - m1.mat[2][2]},
		},
	}
}

func (m *Matrix3d) Mul(m1 *Matrix3d) *Matrix3d {
	return &Matrix3d{
		mat: [3][3]float64{
			{
				m.mat[0][0]*m1.mat[0][0] + m.mat[0][1]*m1.mat[1][0] + m.mat[0][2]*m1.mat[2][0],
				m.mat[0][0]*m1.mat[0][1] + m.mat[0][1]*m1.mat[1][1] + m.mat[0][2]*m1.mat[2][1],
				m.mat[0][0]*m1.mat[0][2] + m.mat[0][1]*m1.mat[1][2] + m.mat[0][2]*m1.mat[2][2]},
			{
				m.mat[1][0]*m1.mat[0][0] + m.mat[1][1]*m1.mat[1][0] + m.mat[1][2]*m1.mat[2][0],
				m.mat[1][0]*m1.mat[0][1] + m.mat[1][1]*m1.mat[1][1] + m.mat[1][2]*m1.mat[2][1],
				m.mat[1][0]*m1.mat[0][2] + m.mat[1][1]*m1.mat[1][2] + m.mat[1][2]*m1.mat[2][2]},
			{
				m.mat[2][0]*m1.mat[0][0] + m.mat[2][1]*m1.mat[1][0] + m.mat[2][2]*m1.mat[2][0],
				m.mat[2][0]*m1.mat[0][1] + m.mat[2][1]*m1.mat[1][1] + m.mat[2][2]*m1.mat[2][1],
				m.mat[2][0]*m1.mat[0][2] + m.mat[2][1]*m1.mat[1][2] + m.mat[2][2]*m1.mat[2][2]},
		},
	}
}

func MatrixTranslateV(v *Vector2d) *Matrix3d {
	return MatrixTranslate(v.x, v.y)
}

func MatrixTranslate(x, y float64) *Matrix3d {
	return &Matrix3d{
		mat: [3][3]float64{
			{1, 0, 0},
			{0, 1, 0},
			{x, y, 1},
		},
	}
}

func MatrixScaleV(v *Vector2d) *Matrix3d {
	return MatrixScale(v.x, v.y)
}

func MatrixScale(x, y float64) *Matrix3d {
	return &Matrix3d{
		mat: [3][3]float64{
			{x, 0, 0},
			{0, y, 0},
			{0, 0, 1},
		},
	}
}

func MatrixShearV(v *Vector2d) *Matrix3d {
	return MatrixShear(v.x, v.y)
}

func MatrixShear(x, y float64) *Matrix3d {
	return &Matrix3d{
		mat: [3][3]float64{
			{x, 0, 0},
			{0, y, 0},
			{0, 0, 1},
		},
	}
}

func MatrixRotate(rad float64) *Matrix3d {
	sin := math.Sin(rad)
	cos := math.Cos(rad)
	return &Matrix3d{
		mat: [3][3]float64{
			{cos, sin, 0},
			{-sin, cos, 0},
			{0, 0, 1},
		},
	}
}

func (m *Matrix3d) MulV(vec *Vector2d) *Vector2d {
	// return NewVector3d(
	// 	vec.x*m.mat[0][0]+vec.y*m.mat[1][0]+vec.w*m.mat[2][0],
	// 	vec.x*m.mat[0][1]+vec.y*m.mat[1][1]+vec.w*m.mat[2][1],
	// 	vec.x*m.mat[0][2]+vec.y*m.mat[1][2]+vec.w*m.mat[2][2])
	return NewVector2d(
		vec.x*m.mat[0][0]+vec.y*m.mat[1][0]+vec.w*m.mat[2][0],
		vec.x*m.mat[0][1]+vec.y*m.mat[1][1]+vec.w*m.mat[2][1])
}
