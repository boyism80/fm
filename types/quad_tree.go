package types

import (
	"golang.org/x/exp/constraints"
)

type AnySpatial[T constraints.Integer] interface {
	Bounds() Rect[T]
}

const (
	NW = iota
	NE
	SW
	SE
)

type QuadTreeNode[T constraints.Integer, U AnySpatial[T]] struct {
	Bounds   Rect[T]
	Items    []U
	Children [4]*QuadTreeNode[T, U]
	Depth    int
}

const (
	maxDepth = 8
)

func NewQuadTreeNode[T constraints.Integer, U AnySpatial[T]](bounds Rect[T], depth int) *QuadTreeNode[T, U] {
	return &QuadTreeNode[T, U]{Bounds: bounds, Depth: depth, Items: []U{}}
}

func (n *QuadTreeNode[T, U]) subdivide() {
	mid := Point[T]{
		X: (n.Bounds.Left + n.Bounds.Right) / 2,
		Y: (n.Bounds.Top + n.Bounds.Bottom) / 2,
	}

	n.Children[NW] = NewQuadTreeNode[T, U](Rect[T]{n.Bounds.Left, n.Bounds.Top, mid.X, mid.Y}, n.Depth+1)
	n.Children[NE] = NewQuadTreeNode[T, U](Rect[T]{mid.X + 1, n.Bounds.Top, n.Bounds.Right, mid.Y}, n.Depth+1)
	n.Children[SW] = NewQuadTreeNode[T, U](Rect[T]{n.Bounds.Left, mid.Y + 1, mid.X, n.Bounds.Bottom}, n.Depth+1)
	n.Children[SE] = NewQuadTreeNode[T, U](Rect[T]{mid.X + 1, mid.Y + 1, n.Bounds.Right, n.Bounds.Bottom}, n.Depth+1)
}

func (n *QuadTreeNode[T, U]) Insert(item U) bool {
	b := item.Bounds()
	if !n.Bounds.Intersects(b) {
		return false
	}
	if n.Depth >= maxDepth {
		n.Items = append(n.Items, item)
		return true
	}
	if n.Children[NW] == nil {
		n.subdivide()
	}
	for _, idx := range []int{NW, NE, SW, SE} {
		c := n.Children[idx]
		if c.Bounds.ContainsRect(b) {
			return c.Insert(item)
		}
	}
	n.Items = append(n.Items, item)
	return true
}

func (n *QuadTreeNode[T, U]) QueryRange(r Rect[T], out *[]U) {
	if !n.Bounds.Intersects(r) {
		return
	}
	for _, item := range n.Items {
		if r.Intersects(item.Bounds()) {
			*out = append(*out, item)
		}
	}
	if n.Children[NW] != nil {
		for _, idx := range []int{NW, NE, SW, SE} {
			n.Children[idx].QueryRange(r, out)
		}
	}
}
