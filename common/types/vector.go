package types

type Vec2 struct {
	X, Y int
}

func (v Vec2) DistanceSq(other Vec2) int {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return dx*dx + dy*dy
}

type Rect struct {
	Left, Top, Right, Bottom int
}

func (r Rect) Contains(pos Vec2) bool {
	return pos.X >= r.Left && pos.X <= r.Right &&
		pos.Y >= r.Top && pos.Y <= r.Bottom
}
