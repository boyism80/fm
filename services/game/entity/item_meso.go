package entity

import (
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

type Meso struct {
	*FieldPlacement
	Count int32
}

func (meso *Meso) GetCount32() int32 { return meso.Count }
func (meso *Meso) GetFieldPlacement() *FieldPlacement {
	return meso.FieldPlacement
}
func (meso *Meso) BindFieldPlacement(placement *FieldPlacement) {
	meso.FieldPlacement = placement
}
func (meso *Meso) IsMeso() bool { return true }

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
		ID:           fp.OID,
		Animation:    constant.DROP_ITEM_ANIMATION_TYPE_NONE,
		DropType:     fp.DropType,
		Count:        meso.Count,
		OwnerID:      fp.Owner,
		Position:     fp.Position,
		SpawnedPoint: fp.SpawnedPoint,
		IsPlayerDrop: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func NewMeso(count int32, position types.Point[int16], ownerID uint32, dropType constant.DropType, sequence uint32, gw GameWorld, mapInstance *Map) *Meso {
	m := &Meso{
		FieldPlacement: &FieldPlacement{
			ObjectCore: &ObjectCore{
				OID:       sequence,
				Position:  position,
				GameWorld: gw,
				Map:       mapInstance,
			},
			SpawnedPoint: position,
			DropType:     dropType,
			Owner:        ownerID,
		},
		Count: count,
	}
	if m.FieldPlacement != nil && m.FieldPlacement.ObjectCore != nil {
		m.FieldPlacement.ObjectCore.self = m
	}
	return m
}
