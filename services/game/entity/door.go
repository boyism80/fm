package entity

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

type DoorEndpoint struct {
	Map      *Map
	PortalID uint8
	Position types.Vector2[int16]
}

type DoorKey string

func NewDoor(key DoorKey, ownerID uint32, skillID constant.SkillID, field, returnEp DoorEndpoint, partyID *uint32) *Door {
	door := &Door{
		ObjectCore: ObjectCore{
			Map: nil,
		},
		Key:     key,
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
	Key     DoorKey
	OwnerID uint32
	SkillID constant.SkillID
	Field   DoorEndpoint
	Return  DoorEndpoint
	PartyID *uint32
	// fieldRef keeps an instance field open while its door stands, so the owner can come back from town.
	fieldRef *MapRef
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
	return ch.Party.Is(d.PartyID)
}

func (d *Door) Destination(from *Map) (DoorEndpoint, bool) {
	switch from {
	case d.Field.Map:
		return d.Return, true
	case d.Return.Map:
		return d.Field, true
	}
	return DoorEndpoint{}, false
}

func (ch *Character) UseDoor(actx actor.Context, ownerID uint32) bool {
	m := ch.GetMap()
	if m == nil || m.Wz == nil {
		return false
	}
	door := m.FindDoorOwner(ownerID)
	if door == nil || door.Map != m || !door.UsableBy(ch) {
		return false
	}
	dest, ok := door.Destination(m)
	if !ok || dest.Map == nil || dest.Map.Wz == nil {
		return false
	}

	if err := ch.Warp(actx, dest.Map, dest.PortalID); err != nil {
		log.Printf("use door warp: %v", err)
		return false
	}
	return true
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
	onFieldMap := vm == d.Field.Map
	isOwner := viewer.GetID() == d.OwnerID
	sameParty := viewer.Party.Is(d.PartyID)
	if !onFieldMap && !isOwner && !sameParty {
		return
	}
	var portalPoint types.Vector2[int16]
	if onFieldMap {
		portalPoint = d.Field.Position
	} else {
		portalPoint = d.Return.Position
	}
	if dm == d.Field.Map {
		viewer.Send(&response.SpawnDoor{
			OwnerID:  d.OwnerID,
			Position: portalPoint,
			Animated: false,
		}, types.SEND_POLICY_ENCRYPT)
	}
	usePartyPortal := d.PartyID != nil && viewer.Party.ID() != nil && (isOwner || sameParty)
	if usePartyPortal {
		viewer.Send(&response.PartyPortal{
			TownMapID:   d.Return.Map.TemplateID(),
			TargetMapID: d.Field.Map.TemplateID(),
			SkillID:     uint32(d.SkillID),
			Position:    portalPoint,
			Animated:    false,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}
	viewer.Send(&response.SpawnPortal{
		DestMapID:   d.Return.Map.TemplateID(),
		SourceMapID: d.Field.Map.TemplateID(),
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
	if vm == d.Field.Map {
		portalPoint = d.Field.Position
	} else {
		portalPoint = d.Return.Position
	}
	_ = viewer.Send(&response.SpawnPortal{
		DestMapID:   d.Return.Map.TemplateID(),
		SourceMapID: d.Field.Map.TemplateID(),
		Position:    &portalPoint,
	}, types.SEND_POLICY_ENCRYPT)
}

func (d *Door) BuildRPC() *internal.PartyDoor {
	return &internal.PartyDoor{
		Town:   d.Return.Map.TemplateID(),
		Target: d.Field.Map.TemplateID(),
		X:      int32(d.Field.Position.X),
		Y:      int32(d.Field.Position.Y),
	}
}
