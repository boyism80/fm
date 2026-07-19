package entity

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

type DoorEndpoint struct {
	MapID    uint32
	PortalID uint8
	Position types.Vector2[int16]
}

func NewDoor(ownerID uint32, skillID constant.SkillID, field, returnEp DoorEndpoint, partyID *uint32) *Door {
	door := &Door{
		ObjectCore: ObjectCore{
			Map: nil,
		},
		OwnerID: ownerID,
		SkillID: skillID,
		Field:   field,
		Return:  returnEp,
		PartyID: partyID,
	}
	door.ObjectCore.self = door
	return door
}

type Door struct {
	ObjectCore
	OwnerID uint32
	SkillID constant.SkillID
	Field   DoorEndpoint
	Return  DoorEndpoint
	PartyID *uint32
}

func (d *Door) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeDoor
}

func (d *Door) Is(typ constant.ObjectType) bool {
	return d.GetObjectType().Has(typ)
}

func (d *Door) UsableBy(ch *Character) bool {
	if d == nil || ch == nil {
		return false
	}
	if ch.GetID() == d.OwnerID {
		return true
	}
	if d.SkillID != constant.SkillMysticDoor {
		return false
	}
	if d.PartyID == nil {
		return false
	}
	vParty := ch.GetPartyID()
	if vParty == nil {
		return false
	}
	return *vParty == *d.PartyID
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
	onFieldMap := viewerMapID == d.Field.MapID
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
		portalPoint = d.Field.Position
	} else {
		portalPoint = d.Return.Position
	}
	if dm.Wz.ID == d.Field.MapID {
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
			TownMapID:   d.Return.MapID,
			TargetMapID: d.Field.MapID,
			SkillID:     uint32(d.SkillID),
			Position:    portalPoint,
			Animated:    false,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}
	viewer.Send(&response.SpawnPortal{
		DestMapID:   d.Return.MapID,
		SourceMapID: d.Field.MapID,
		SkillID:     uint32(d.SkillID),
		Position:    &portalPoint,
	}, types.SEND_POLICY_ENCRYPT)
}

func (d *Door) SendDestroySyncToViewer(viewer *Character) {
	if d == nil || viewer == nil {
		return
	}
	viewer.Send(&response.RemoveDoor{
		OwnerID:  d.OwnerID,
		Animated: false,
	}, types.SEND_POLICY_ENCRYPT)
	viewer.Send(&response.SpawnPortal{
		DestMapID:   response.DisabledPortalMapID,
		SourceMapID: response.DisabledPortalMapID,
		SkillID:     0,
		Position:    nil,
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

func (d *Door) SendOwnerPortalResync(viewer *Character) {
	if d == nil || viewer == nil {
		return
	}
	if viewer.GetID() != d.OwnerID {
		return
	}
	vm := viewer.GetMap()
	if vm == nil || vm.Wz == nil {
		return
	}
	var portalPoint types.Vector2[int16]
	if vm.Wz.ID == d.Field.MapID {
		portalPoint = d.Field.Position
	} else {
		portalPoint = d.Return.Position
	}
	_ = viewer.Send(&response.SpawnPortal{
		DestMapID:   d.Return.MapID,
		SourceMapID: d.Field.MapID,
		SkillID:     uint32(d.SkillID),
		Position:    &portalPoint,
	}, types.SEND_POLICY_ENCRYPT)
}

func (d *Door) ToProto() *internal.PartyDoor {
	return &internal.PartyDoor{
		Town:   d.Return.MapID,
		Target: d.Field.MapID,
		X:      int32(d.Field.Position.X),
		Y:      int32(d.Field.Position.Y),
	}
}
