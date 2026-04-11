package types

import "golang.org/x/exp/constraints"

type Vector2[T constraints.Integer] struct {
	X, Y T
}

type Point[T constraints.Integer] = Vector2[T]

func (v Vector2[T]) DistanceSq(other Vector2[T]) int64 {
	dx := int64(v.X) - int64(other.X)
	dy := int64(v.Y) - int64(other.Y)
	return dx*dx + dy*dy
}
