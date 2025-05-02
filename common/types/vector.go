package types

import "golang.org/x/exp/constraints"

type Vec2[T constraints.Integer] struct {
	X, Y T
}

func (v Vec2[T]) DistanceSq(other Vec2[T]) T {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return dx*dx + dy*dy
}

type Rect[T constraints.Integer] struct {
	Left, Top, Right, Bottom T
}

func (r Rect[T]) Contains(pos Vec2[T]) bool {
	return pos.X >= r.Left && pos.X <= r.Right &&
		pos.Y >= r.Top && pos.Y <= r.Bottom
}
