package types

import "golang.org/x/exp/constraints"

type Rect[T constraints.Integer] struct {
	Left, Top, Right, Bottom T
}

func NewRect[T constraints.Integer](lt, rb Point[T]) Rect[T] {
	return Rect[T]{
		Left:   min(lt.X, rb.X),
		Top:    min(lt.Y, rb.Y),
		Right:  max(lt.X, rb.X),
		Bottom: max(lt.Y, rb.Y),
	}
}

func (r Rect[T]) Valid() bool {
	return r.Left != 0 || r.Top != 0 || r.Right != 0 || r.Bottom != 0
}

func (r Rect[T]) AtOrigin(origin Point[T]) Rect[T] {
	return Rect[T]{
		Left:   origin.X + r.Left,
		Top:    origin.Y + r.Top,
		Right:  origin.X + r.Right,
		Bottom: origin.Y + r.Bottom,
	}
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
