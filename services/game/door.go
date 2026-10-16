package game

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
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
		objectCore: ObjectCore{
			onMap: nil,
		},
		key:      key,
		ownerID:  ownerID,
		skillID:  skillID,
		field:    field,
		onReturn: returnEp,
		partyID:  partyID,
	}
	door.objectCore.self = door
	return door
}

type Door struct {
	objectCore
	key      DoorKey
	ownerID  uint32
	skillID  constant.SkillID
	field    DoorEndpoint
	onReturn DoorEndpoint
	partyID  *uint32
	// fieldRef keeps an instance field open while its door stands, so the owner can come back from town.
	fieldRef *MapRef
}

func (d *Door) SetPartyID(partyID *uint32) {
	d.partyID = partyID
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
	if ch.GetID() == d.ownerID {
		return true
	}
	if d.skillID != constant.SkillMysticDoor {
		return false
	}
	return ch.party.Is(d.partyID)
}

func (d *Door) Destination(from *Map) (DoorEndpoint, bool) {
	switch from {
	case d.field.Map:
		return d.onReturn, true
	case d.onReturn.Map:
		return d.field, true
	}
	return DoorEndpoint{}, false
}

func (d *Door) Use(actx actor.Context, ch *Character) bool {
	m := ch.GetMap()
	if m == nil || m.wz == nil || d.onMap != m || !d.UsableBy(ch) {
		return false
	}
	dest, ok := d.Destination(m)
	if !ok || dest.Map == nil || dest.Map.wz == nil {
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
	if vm == nil || vm.wz == nil {
		return
	}
	dm := d.GetMap()
	if dm == nil || dm.wz == nil {
		return
	}
	onFieldMap := vm == d.field.Map
	isOwner := viewer.GetID() == d.ownerID
	sameParty := viewer.party.Is(d.partyID)
	if !onFieldMap && !isOwner && !sameParty {
		return
	}
	var portalPoint types.Vector2[int16]
	if onFieldMap {
		portalPoint = d.field.Position
	} else {
		portalPoint = d.onReturn.Position
	}
	if dm == d.field.Map {
		viewer.Send(&response.SpawnDoor{
			OwnerID:  d.ownerID,
			Position: portalPoint,
			Animated: false,
		}, types.SEND_POLICY_ENCRYPT)
	}
	usePartyPortal := d.partyID != nil && viewer.party.ID() != nil && (isOwner || sameParty)
	if usePartyPortal {
		viewer.Send(&response.PartyPortal{
			TownMapID:   d.onReturn.Map.TemplateID(),
			TargetMapID: d.field.Map.TemplateID(),
			SkillID:     uint32(d.skillID),
			Position:    portalPoint,
			Animated:    false,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}
	viewer.Send(&response.SpawnPortal{
		DestMapID:   d.onReturn.Map.TemplateID(),
		SourceMapID: d.field.Map.TemplateID(),
		Position:    &portalPoint,
	}, types.SEND_POLICY_ENCRYPT)
}

func (d *Door) SendDestroySyncToViewer(viewer *Character) {
	if d == nil || viewer == nil {
		return
	}
	viewer.Send(&response.RemoveDoor{
		OwnerID:  d.ownerID,
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
	if m == nil || d.oid == 0 {
		return
	}
	m.RemoveDoor(d.oid, animated)
}

func (d *Door) SendOwnerPortalResync(viewer *Character) {
	if d == nil || viewer == nil {
		return
	}
	if viewer.GetID() != d.ownerID {
		return
	}
	vm := viewer.GetMap()
	if vm == nil || vm.wz == nil {
		return
	}
	var portalPoint types.Vector2[int16]
	if vm == d.field.Map {
		portalPoint = d.field.Position
	} else {
		portalPoint = d.onReturn.Position
	}
	_ = viewer.Send(&response.SpawnPortal{
		DestMapID:   d.onReturn.Map.TemplateID(),
		SourceMapID: d.field.Map.TemplateID(),
		Position:    &portalPoint,
	}, types.SEND_POLICY_ENCRYPT)
}

func (d *Door) BuildRPC() *internal.PartyDoor {
	return &internal.PartyDoor{
		Town:   d.onReturn.Map.TemplateID(),
		Target: d.field.Map.TemplateID(),
		X:      int32(d.field.Position.X),
		Y:      int32(d.field.Position.Y),
	}
}
