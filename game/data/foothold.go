// Package data provides MapleStory game data specifications and types.
// This file contains foothold collision data structures for map navigation.
package data

import "github.com/boyism80/fm/core/types"

// Foothold represents a walkable platform or surface in a MapleStory map.
// Footholds define where characters can stand and walk, providing collision detection.
type Foothold struct {
	X1, Y1, X2, Y2 int16 // Line segment coordinates
	ID             int16 // Foothold identifier
	Prev, Next     int16 // Connected foothold IDs
}

// Bounds returns the rectangular boundary of this foothold.
// Used by the QuadTree for spatial partitioning and collision detection.
func (f Foothold) Bounds() types.Rect[int16] {
	left := min(f.X1, f.X2)
	right := max(f.X1, f.X2)
	top := min(f.Y1, f.Y2)
	bottom := max(f.Y1, f.Y2)
	return types.Rect[int16]{Left: left, Top: top, Right: right, Bottom: bottom}
}

// Compare determines the rendering order of overlapping footholds.
// Returns true if this foothold should be processed before the other.
// Used for proper collision detection when multiple footholds overlap.
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
