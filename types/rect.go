package types

import "golang.org/x/exp/constraints"

type Rect[T constraints.Integer] struct {
	Left, Top, Right, Bottom T
}

func (r Rect[T]) ContainsPoint(pos Point[T]) bool {
	return pos.X >= r.Left && pos.X <= r.Right &&
		pos.Y >= r.Top && pos.Y <= r.Bottom
}

func (r Rect[T]) ContainsRect(other Rect[T]) bool {
	return other.Left >= r.Left && other.Right <= r.Right &&
		other.Top >= r.Top && other.Bottom <= r.Bottom
}

func (r Rect[T]) Intersects(other Rect[T]) bool {
	return !(other.Left > r.Right || other.Right < r.Left ||
		other.Top > r.Bottom || other.Bottom < r.Top)
}
