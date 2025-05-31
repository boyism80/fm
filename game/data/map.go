package data

import (
	"math"

	"github.com/boyism80/fm/common/types"
)

type Portal struct {
	ID          uint8
	Name        string
	TargetMapId int32
	Target      string
	Position    types.Point[int16]
	ScriptName  string
	Type        uint8
}

type MapSpec struct {
	ID           uint32
	Name         string
	Version      int
	Cloud        int
	ReturnMapId  int
	ForcedReturn int
	FieldLimit   int
	VRTop        int
	VRLeft       int
	VRBottom     int
	VRRight      int
	HideMinimap  bool
	IsTown       bool
	MobRate      float32
	BGM          string
	MapMark      string
	MapDesc      string
	MiniMapOnOff bool
	Portals      map[uint8]Portal
	Footholds    *types.QuadTreeNode[int16, Foothold]
}

func (ms *MapSpec) FootholdPoint(point types.Point[int16]) *types.Point[int16] {
	foothold, ok := ms.Footholds.Find(point)
	if !ok {
		return nil
	}

	top := foothold.Y1
	if foothold.X1 != foothold.X2 && foothold.Y1 != foothold.Y2 {
		s1 := float64(math.Abs(float64(foothold.Y2 - foothold.Y1)))
		s2 := float64(math.Abs(float64(foothold.X2 - foothold.X1)))
		dx := float64(math.Abs(float64(point.X - foothold.X1)))

		alpha := math.Atan(s2 / s1)
		beta := math.Atan(s1 / s2)
		offset := math.Cos(alpha) * (dx / math.Cos(beta))

		if foothold.Y2 < foothold.Y1 {
			top = foothold.Y1 - int16(offset)
		} else {
			top = foothold.Y1 + int16(offset)
		}
	}

	pt := types.Point[int16]{X: point.X, Y: top}
	return &pt
}

func (ms *MapSpec) DropPoint(initial types.Point[int16]) types.Point[int16] {
	highest := types.Point[int16]{X: initial.X, Y: initial.Y - int16(100)}
	if result := ms.FootholdPoint(highest); result != nil {
		return *result
	}
	return initial
}

func (spec *MapSpec) FindPortal(name string) (*Portal, bool) {
	for _, portal := range spec.Portals {
		if portal.Name == name {
			return &portal, true
		}
	}

	return nil, false
}
