package wz

import (
	"sort"

	"github.com/boyism80/fm/services/game/constant"
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

const forcedReturnNone = 999999999

const (
	freeMarketFirstMapID    = 910000000
	freeMarketLastMapID     = 910000022
	messageBoxBlockedRegion = 91
)

type FieldType int

const (
	FieldTypeCoconut               FieldType = 4
	FieldTypeMonsterCarnival       FieldType = 10
	FieldTypeMonsterCarnivalRevive FieldType = 11
	FieldTypeDojo                  FieldType = 14
)

type Map struct {
	ID           uint32
	Name         string
	Version      int
	Cloud        int
	ReturnMapId  int
	ForcedReturn uint32
	FieldLimit   constant.FieldLimit
	FieldType    FieldType

	Continent         uint32
	FreeMarket        bool
	MessageBoxBlocked bool

	VRTop             int
	VRLeft            int
	VRBottom          int
	VRRight           int
	HideMinimap       bool
	IsTown            bool
	Everlast          bool
	HasClock          bool
	MobRate           float32
	RecoveryRate      float32
	BGM               string
	MapMark           string
	MapDesc           string
	MiniMapOnOff      bool
	PersonalShop      bool
	EntrustedShop     bool
	Portals           map[uint8]Portal
	NpcSpawns         map[uint32]NpcSpawn
	MobSpawns         map[uint32]MobSpawn
	ReactorSpawns     map[uint32]ReactorSpawn
	Areas             []types.Rect[int16]
	Footholds         *types.QuadTreeNode[int16, Foothold]
	Carnival          *CarnivalField
	doorReturnPortals []Portal
}

func (m *Map) Limits(limit constant.FieldLimit) bool {
	return m != nil && m.FieldLimit&limit != 0
}

// The client field of these types reads one extra team byte at the end of the spawn player packet.
func (m *Map) HasTeam() bool {
	if m == nil {
		return false
	}
	switch m.FieldType {
	case FieldTypeCoconut, FieldTypeMonsterCarnival, FieldTypeMonsterCarnivalRevive:
		return true
	default:
		return false
	}
}

func (model *Map) FootholdBelow(p types.Point[int16]) (*Foothold, int16, bool) {
	if model == nil || model.Footholds == nil {
		return nil, 0, false
	}
	bottom := model.Footholds.Bounds.Bottom
	if p.Y > bottom {
		return nil, 0, false
	}
	var candidates []Foothold
	model.Footholds.QueryRange(types.Rect[int16]{
		Left:   p.X,
		Top:    p.Y,
		Right:  p.X,
		Bottom: bottom,
	}, &candidates)
	xMatches := make([]Foothold, 0, len(candidates))
	for _, fh := range candidates {
		if fh.X1 <= p.X && fh.X2 >= p.X {
			xMatches = append(xMatches, fh)
		}
	}
	sort.Slice(xMatches, func(i, j int) bool {
		return xMatches[i].Y2 < xMatches[j].Y1
	})
	for i := range xMatches {
		if xMatches[i].Wall {
			continue
		}
		if y := xMatches[i].SurfaceY(p.X); y >= p.Y {
			return &xMatches[i], y, true
		}
	}
	return nil, 0, false
}

func (model *Map) PointBelow(point types.Point[int16]) *types.Point[int16] {
	_, y, ok := model.FootholdBelow(point)
	if !ok {
		return nil
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
	return portal.Position, true
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
