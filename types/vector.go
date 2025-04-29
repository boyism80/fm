package types

type Vec2 struct {
	X, Y int
}

func (v Vec2) DistanceSq(other Vec2) int {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return dx*dx + dy*dy
}
