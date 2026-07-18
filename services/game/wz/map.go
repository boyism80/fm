package wz

import (
	"math"
	"sort"

	"github.com/boyism80/fm/types"
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

type NpcSpawn struct {
	*BaseSpawn
}

type MobSpawn struct {
	*BaseSpawn
}

const dropPointSearchOffset int16 = 50

type Map struct {
	ID                uint32
	Name              string
	Version           int
	Cloud             int
	ReturnMapId       int
	ForcedReturn      int
	FieldLimit        int
	VRTop             int
	VRLeft            int
	VRBottom          int
	VRRight           int
	HideMinimap       bool
	IsTown            bool
	Everlast          bool
	MobRate           float32
	RecoveryRate      float32
	BGM               string
	MapMark           string
	MapDesc           string
	MiniMapOnOff      bool
	Portals           map[uint8]Portal
	NpcSpawns         map[uint32]NpcSpawn
	MobSpawns         map[uint32]MobSpawn
	ReactorSpawns     map[uint32]ReactorSpawn
	Areas             []types.Rect[int16]
	Footholds         *types.QuadTreeNode[int16, Foothold]
	doorReturnPortals []Portal
}

func footholdSpansX(f Foothold, x int16) bool {
	xLo, xHi := f.X1, f.X2
	if xLo > xHi {
		xLo, xHi = xHi, xLo
	}
	return xLo <= x && x <= xHi
}

func footholdSurfaceYAtX(f Foothold, x int16) int16 {
	if f.X1 == f.X2 {
		return f.Y1
	}
	if f.Y1 == f.Y2 {
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

func (model *Map) findBelow(p types.Point[int16]) (*Foothold, bool) {
	if model == nil || model.Footholds == nil {
		return nil, false
	}
	rels := model.Footholds.Relations(p)
	xMatches := make([]Foothold, 0, len(rels))
	for _, fh := range rels {
		if footholdSpansX(fh, p.X) {
			xMatches = append(xMatches, fh)
		}
	}
	sort.Slice(xMatches, func(i, j int) bool {
		return xMatches[i].Compare(xMatches[j])
	})
	for i := range xMatches {
		fh := xMatches[i]
		if fh.IsWall() {
			continue
		}
		if fh.X1 != fh.X2 && fh.Y1 != fh.Y2 {
			calcY := footholdSurfaceYAtX(fh, p.X)
			if calcY >= p.Y {
				return &xMatches[i], true
			}
		} else {
			if fh.Y1 >= p.Y {
				return &xMatches[i], true
			}
		}
	}
	return nil, false
}

func (model *Map) PointBelow(point types.Point[int16]) *types.Point[int16] {
	if model == nil || model.Footholds == nil {
		return nil
	}
	fh, ok := model.findBelow(point)
	if !ok {
		return nil
	}
	var y int16
	if !fh.IsWall() && fh.Y1 != fh.Y2 {
		y = footholdSurfaceYAtX(*fh, point.X)
	} else {
		y = fh.Y1
	}
	return &types.Point[int16]{X: point.X, Y: y}
}

func (model *Map) DropPoint(initial types.Point[int16]) (types.Point[int16], bool) {
	search := types.Point[int16]{X: initial.X, Y: initial.Y - dropPointSearchOffset}
	if result := model.PointBelow(search); result != nil {
		return *result, true
	}
	return initial, false
}

func (model *Map) GetSpawnPosition(spawnPoint uint8) (types.Point[int16], bool) {
	if model == nil {
		return types.Point[int16]{}, false
	}
	portal, ok := model.Portals[spawnPoint]
	if !ok {
		return types.Point[int16]{}, false
	}
	position := portal.Position
	if snapped := model.PointBelow(position); snapped != nil {
		position = *snapped
	}
	return position, true
}

func (model *Map) FindClosestPortalSpawnID(pos types.Point[int16]) uint8 {
	if model == nil || len(model.Portals) == 0 {
		return 0
	}
	var best uint8
	var bestDist int64 = -1
	found := false
	for id, p := range model.Portals {
		if p.Type != 0 && p.Name != "sp" {
			continue
		}
		dx := int64(p.Position.X - pos.X)
		dy := int64(p.Position.Y - pos.Y)
		d := dx*dx + dy*dy
		if !found || d < bestDist {
			bestDist = d
			best = id
			found = true
		}
	}
	if !found {
		return 0
	}
	return best
}

func (model *Map) FindClosestDoorReturnPortalSpawnID(pos types.Point[int16]) (uint8, bool) {
	if model == nil || len(model.doorReturnPortals) == 0 {
		return 0, false
	}
	var best uint8
	var bestDist int64 = -1
	for i := range model.doorReturnPortals {
		p := model.doorReturnPortals[i]
		dx := int64(p.Position.X - pos.X)
		dy := int64(p.Position.Y - pos.Y)
		d := dx*dx + dy*dy
		if bestDist < 0 || d < bestDist {
			bestDist = d
			best = p.ID
		}
	}
	return best, true
}

func (model *Map) buildDoorReturnPortal() {
	if model == nil || len(model.Portals) == 0 {
		if model != nil {
			model.doorReturnPortals = nil
		}
		return
	}
	out := make([]Portal, 0, len(model.Portals))
	for _, p := range model.Portals {
		if p.Type == 6 {
			out = append(out, p)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	model.doorReturnPortals = out
}

func (model *Map) DoorReturnPortalSlots() []Portal {
	if model == nil || len(model.doorReturnPortals) == 0 {
		return nil
	}
	out := make([]Portal, len(model.doorReturnPortals))
	copy(out, model.doorReturnPortals)
	return out
}

func (model *Map) GetDoorReturnPosition(partyOwnerSlot int) (types.Point[int16], bool) {
	if model == nil {
		return types.Point[int16]{}, false
	}
	list := model.doorReturnPortals
	if len(list) == 0 {
		return types.Point[int16]{}, false
	}
	if partyOwnerSlot < 0 {
		partyOwnerSlot = 0
	}
	if partyOwnerSlot >= len(list) {
		partyOwnerSlot = len(list) - 1
	}
	return list[partyOwnerSlot].Position, true
}

func (model *Map) DoorReturnPortalSpawnID(partyOwnerSlot int) (uint8, bool) {
	if model == nil {
		return 0, false
	}
	list := model.doorReturnPortals
	if len(list) == 0 {
		return 0, false
	}
	if partyOwnerSlot < 0 {
		partyOwnerSlot = 0
	}
	if partyOwnerSlot >= len(list) {
		partyOwnerSlot = len(list) - 1
	}
	return list[partyOwnerSlot].ID, true
}
