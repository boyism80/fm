package game

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type Meso struct {
	*fieldPlacement
	count int32
}

func (meso *Meso) GetCount32() int32 { return meso.count }
func (meso *Meso) GetFieldPlacement() *FieldPlacement {
	return meso.fieldPlacement
}
func (meso *Meso) BindFieldPlacement(placement *FieldPlacement) {
	meso.fieldPlacement = placement
}

func (meso *Meso) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeItem
}

func (meso *Meso) SendSpawnSyncToViewer(viewer *Character) {
	if meso == nil || viewer == nil {
		return
	}
	fp := meso.GetFieldPlacement()
	if fp == nil {
		return
	}
	viewer.Send(&response.SpawnMeso{
		ID:           fp.oid,
		Animation:    constant.DropItemAnimationTypeNone,
		DropType:     fp.dropType,
		Count:        meso.count,
		OwnerID:      fp.owner,
		Position:     fp.position,
		SpawnedPoint: fp.spawnedPoint,
		IsPlayerDrop: fp.playerDrop,
	}, types.SEND_POLICY_ENCRYPT)
}

func (meso *Meso) SendDestroySyncToViewer(viewer *Character) {
	if meso == nil || viewer == nil {
		return
	}
	fp := meso.GetFieldPlacement()
	if fp == nil {
		return
	}
	viewer.Send(&response.RemoveItem{
		Mode:        constant.RemoveItemTypeNoAnimated,
		OID:         fp.oid,
		CharacterId: 0,
	}, types.SEND_POLICY_ENCRYPT)
}

func NewMeso(count int32, position types.Point[int16], ownerID uint32, dropType constant.DropType, sequence uint32, gw GameWorld, mapInstance *Map) *Meso {
	m := &Meso{
		fieldPlacement: &FieldPlacement{
			ObjectCore: &ObjectCore{
				oid:       sequence,
				position:  position,
				gameWorld: gw,
				onMap:     mapInstance,
			},
			spawnedPoint: position,
			dropType:     dropType,
			owner:        ownerID,
		},
		count: count,
	}
	if m.fieldPlacement != nil && m.fieldPlacement.ObjectCore != nil {
		m.fieldPlacement.ObjectCore.self = m
		m.fieldPlacement.initTimers()
	}
	return m
}

func (meso *Meso) inArea(area types.Rect[int32]) bool {
	if meso.GetFieldPlacement() == nil {
		return false
	}
	pos := meso.GetPosition()
	return area.ContainsPoint(types.Point[int32]{X: int32(pos.X), Y: int32(pos.Y)})
}
