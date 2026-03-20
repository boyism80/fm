package entity

import (
	"github.com/boyism80/fm/game/constant"
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

func NewMeso(count int32, position types.Point[int16], ownerID uint32, dropType constant.DropType, sequence uint32, context GameContext, mapInstance *Map) *Meso {
	return &Meso{
		Drop: &Drop{
			ObjectCore: &ObjectCore{
				OID:      sequence,
				Position: position,
				Context:  context,
				Map:      mapInstance,
			},
			SpawnedPoint: position,
			DropType:     dropType,
			Owner:        ownerID,
		},
		Count: count,
	}
}
