package data

import "github.com/boyism80/fm/common/types"

type Foothold struct {
	X1, Y1, X2, Y2 int16
	ID             int16
	Prev, Next     int16
}

func (f Foothold) Bounds() types.Rect[int16] {
	left := min(f.X1, f.X2)
	right := max(f.X1, f.X2)
	top := min(f.Y1, f.Y2)
	bottom := max(f.Y1, f.Y2)
	return types.Rect[int16]{Left: left, Top: top, Right: right, Bottom: bottom}
}

func (f Foothold) Compare(o types.AnySpatial[int16]) bool {
	other, ok := o.(Foothold)
	if !ok {
		return false
	}

	if f.Y2 < other.Y1 {
		return true
	}
	if f.Y1 > other.Y2 {
		return false
	}
	fTop := min(f.Y1, f.Y2)
	oTop := min(other.Y1, other.Y2)
	if fTop != oTop {
		return fTop < oTop
	}
	return f.ID < other.ID
}
