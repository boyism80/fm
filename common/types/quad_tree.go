package types

import (
	"math"
	"sort"

	"golang.org/x/exp/constraints"
)

type AnySpatial[T constraints.Integer] interface {
	Bounds() Rect[T]
	Compare(other AnySpatial[T]) bool
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
	maxItems = 4
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
	if len(n.Items) < maxItems || n.Depth >= maxDepth {
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

func (n *QuadTreeNode[T, U]) relations(p Vector2[T]) []U {
	list := []U{}
	current := n
	for current != nil {
		list = append(list, current.Items...)
		mid := Point[T]{
			X: (current.Bounds.Left + current.Bounds.Right) / 2,
			Y: (current.Bounds.Top + current.Bounds.Bottom) / 2,
		}
		var idx int
		if p.X <= mid.X {
			if p.Y <= mid.Y {
				idx = NW
			} else {
				idx = SW
			}
		} else {
			if p.Y <= mid.Y {
				idx = NE
			} else {
				idx = SE
			}
		}
		current = current.Children[idx]
	}
	return list
}

func (n *QuadTreeNode[T, U]) Find(p Vector2[T]) (*U, bool) {
	relations := n.relations(p)
	items := []U{}
	for _, item := range relations {
		b := item.Bounds()
		if b.Left <= p.X && p.X <= b.Right {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Compare(items[j])
	})
	for _, item := range items {
		bound := item.Bounds()
		if bound.Left != bound.Right && bound.Top != bound.Bottom {
			s1 := float64(math.Abs(float64(bound.Bottom - bound.Top)))
			s2 := float64(math.Abs(float64(bound.Right - bound.Left)))
			s4 := float64(math.Abs(float64(bound.Left - p.X)))
			alpha := math.Atan(s2 / s1)
			beta := math.Atan(s1 / s2)
			s5 := math.Cos(alpha) * (s4 / math.Cos(beta))
			var calcY T
			if bound.Bottom < bound.Top {
				calcY = T(bound.Top) - T(int64(s5))
			} else {
				calcY = T(bound.Top) + T(int64(s5))
			}
			if calcY >= p.Y {
				return &item, true
			}
		} else {
			if T(bound.Top) >= p.Y {
				return &item, true
			}
		}
	}
	return nil, false
}
