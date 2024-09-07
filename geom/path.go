package geom

type Path interface {
	IsFinished() bool
	CurrentWaypoint() *Vector2d
	NextWaypoint()
}

type BasePath struct {
	waypoints []*Vector2d // 不包括起点
	idx       int
	size      int
}

func NewPath(waypoints ...*Vector2d) Path {
	return &BasePath{
		waypoints: waypoints,
		size:      len(waypoints),
		idx:       0,
	}
}

func (p *BasePath) IsFinished() bool {
	return p.idx >= p.size-1
}

func (p *BasePath) CurrentWaypoint() *Vector2d {
	if p.idx >= p.size {
		return nil
	}
	return p.waypoints[p.idx]
}

func (p *BasePath) NextWaypoint() {
	p.idx++
}
