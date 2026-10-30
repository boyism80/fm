package wz

import (
	"math"

	"github.com/boyism80/fm/types"
)

// Foothold is one segment of a map's ground. ZMass is its group in Map.wz (foothold/page/zmass/id); in the air the
// client hits walls and ceilings only of the group it last stood on. Force is the conveyor push of the foothold
// (Map.wz force / 100), and ForbidFallDown refuses a jump down through it.
type Foothold struct {
	X1, Y1, X2, Y2 int16
	ID             int16
	Prev, Next     int16
	Wall           bool
	ZMass          int16
	Force          float64
	ForbidFallDown bool
}

// LadderRope is a ladder or rope at X from the top Y1 to the bottom Y2; UpperFloor lets a climber step off its top.
type LadderRope struct {
	ID         int16
	X, Y1, Y2  int16
	Ladder     bool
	UpperFloor bool
}

func (f Foothold) Bounds() types.Rect[int16] {
	left := min(f.X1, f.X2)
	right := max(f.X1, f.X2)
	top := min(f.Y1, f.Y2)
	bottom := max(f.Y1, f.Y2)
	return types.Rect[int16]{Left: left, Top: top, Right: right, Bottom: bottom}
}

func (f Foothold) SurfaceY(x int16) int16 {
	if f.X1 == f.X2 || f.Y1 == f.Y2 {
		return f.Y1
	}
	s1 := math.Abs(float64(f.Y2 - f.Y1))
	s2 := math.Abs(float64(f.X2 - f.X1))
	dx := math.Abs(float64(x - f.X1))
	alpha := math.Atan(s2 / s1)
	beta := math.Atan(s1 / s2)
	offset := math.Cos(alpha) * (dx / math.Cos(beta))
	if f.Y2 < f.Y1 {
		return f.Y1 - int16(offset)
	}
	return f.Y1 + int16(offset)
}
