package entity

import (
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

func NewDoor(position types.Vector2[int16], ownerID uint32, skillID constant.SkillID, returnMapID, fieldMapID uint32, returnPortalID, fieldPortalID uint8) *Door {
	door := &Door{
		ObjectCore: ObjectCore{
			Position: position,
			Map:      nil,
		},
		OwnerID:        ownerID,
		SkillID:        skillID,
		ReturnMapID:    returnMapID,
		FieldMapID:     fieldMapID,
		ReturnPortalID: returnPortalID,
		FieldPortalID:  fieldPortalID,
	}
	door.ObjectCore.self = door
	return door
}

type Door struct {
	ObjectCore
	OwnerID        uint32
	SkillID        constant.SkillID
	ReturnMapID    uint32
	FieldMapID     uint32
	ReturnPortalID uint8
	FieldPortalID  uint8
}

func (d *Door) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeDoor
}

func (d *Door) Is(typ constant.ObjectType) bool {
	return d.GetObjectType().Has(typ)
}

func (d *Door) SendSpawnSyncToViewer(viewer *Character) {
	if d == nil || viewer == nil {
		return
	}
	viewer.Send(&response.SpawnDoor{
		OwnerID:  d.OwnerID,
		Position: d.Position,
		Animated: false,
	}, types.SEND_POLICY_ENCRYPT)
	pos := d.Position
	viewer.Send(&response.SpawnPortal{
		DestMapID:   d.ReturnMapID,
		SourceMapID: d.FieldMapID,
		SkillID:     uint32(d.SkillID),
		Position:    &pos,
	}, types.SEND_POLICY_ENCRYPT)
}

func (d *Door) Spawn(animated bool) {
	m := d.GetMap()
	if m == nil {
		return
	}
	m.AddDoor(d)
}

func (d *Door) Remove(animated bool) {
	m := d.GetMap()
	if m == nil || d.OID == 0 {
		return
	}
	m.RemoveDoor(d.OID, animated)
}
