package entity

import (
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

type Meso struct {
	*Drop
	Count int32
}

func (meso *Meso) GetCount32() int32   { return meso.Count }
func (meso *Meso) GetDrop() *Drop      { return meso.Drop }
func (meso *Meso) BindDrop(drop *Drop) { meso.Drop = drop }
func (meso *Meso) IsMeso() bool        { return true }

func (meso *Meso) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeItem
}

func (meso *Meso) SendSpawnSyncToViewer(viewer *Character) {
	if meso == nil || viewer == nil {
		return
	}
	drop := meso.GetDrop()
	if drop == nil {
		return
	}
	viewer.Send(&response.SpawnMeso{
		ID:           drop.OID,
		Animation:    constant.DROP_ITEM_ANIMATION_TYPE_NONE,
		DropType:     drop.DropType,
		Count:        meso.Count,
		OwnerID:      drop.Owner,
		Position:     drop.Position,
		SpawnedPoint: drop.SpawnedPoint,
		IsPlayerDrop: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func NewMeso(count int32, position types.Point[int16], ownerID uint32, dropType constant.DropType, sequence uint32, gw GameWorld, mapInstance *Map) *Meso {
	m := &Meso{
		Drop: &Drop{
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
	if m.Drop != nil && m.Drop.ObjectCore != nil {
		m.Drop.ObjectCore.self = m
	}
	return m
}
