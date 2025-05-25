package types

import "golang.org/x/exp/constraints"

type Vector2[T constraints.Integer] struct {
	X, Y T
}

type Point[T constraints.Integer] = Vector2[T]

func (v Vector2[T]) DistanceSq(other Vector2[T]) T {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return dx*dx + dy*dy
}
