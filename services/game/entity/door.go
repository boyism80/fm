package entity

import (
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

func NewDoor(
	position types.Vector2[int16],
	ownerID uint32,
	skillID constant.SkillID,
	returnMapID, fieldMapID uint32,
	returnPortalID, fieldPortalID uint8,
	partyID *uint32,
	fieldPosition, townPortalPosition types.Vector2[int16],
) *Door {
	door := &Door{
		ObjectCore: ObjectCore{
			Position: position,
			Map:      nil,
		},
		OwnerID:            ownerID,
		SkillID:            skillID,
		ReturnMapID:        returnMapID,
		FieldMapID:         fieldMapID,
		ReturnPortalID:     returnPortalID,
		FieldPortalID:      fieldPortalID,
		PartyID:            partyID,
		FieldPosition:      fieldPosition,
		TownPortalPosition: townPortalPosition,
	}
	door.ObjectCore.self = door
	return door
}

type Door struct {
	ObjectCore
	OwnerID            uint32
	SkillID            constant.SkillID
	ReturnMapID        uint32
	FieldMapID         uint32
	ReturnPortalID     uint8
	FieldPortalID      uint8
	PartyID            *uint32
	FieldPosition      types.Vector2[int16]
	TownPortalPosition types.Vector2[int16]
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
	vm := viewer.GetMap()
	if vm == nil || vm.Wz == nil {
		return
	}
	dm := d.GetMap()
	if dm == nil || dm.Wz == nil {
		return
	}
	viewerMapID := vm.Wz.ID
	onFieldMap := viewerMapID == d.FieldMapID
	isOwner := viewer.GetID() == d.OwnerID
	sameParty := false
	if d.PartyID != nil {
		if vid := viewer.GetPartyID(); vid != nil && *vid == *d.PartyID {
			sameParty = true
		}
	}
	if !onFieldMap && !isOwner && !sameParty {
		return
	}
	var portalPoint types.Vector2[int16]
	if onFieldMap {
		portalPoint = d.FieldPosition
	} else {
		portalPoint = d.TownPortalPosition
	}
	if dm.Wz.ID == d.FieldMapID {
		viewer.Send(&response.SpawnDoor{
			OwnerID:  d.OwnerID,
			Position: portalPoint,
			Animated: false,
		}, types.SEND_POLICY_ENCRYPT)
	}
	vParty := viewer.GetPartyID()
	usePartyPortal := d.PartyID != nil && vParty != nil && (isOwner || *vParty == *d.PartyID)
	if usePartyPortal {
		viewer.Send(&response.PartyPortal{
			TownMapID:   d.ReturnMapID,
			TargetMapID: d.FieldMapID,
			SkillID:     uint32(d.SkillID),
			Position:    portalPoint,
			Animated:    false,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}
	viewer.Send(&response.SpawnPortal{
		DestMapID:   d.ReturnMapID,
		SourceMapID: d.FieldMapID,
		SkillID:     uint32(d.SkillID),
		Position:    &portalPoint,
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
