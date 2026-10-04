package entity

import (
	"fmt"

	"github.com/boyism80/fm/services/game/constant"
)

type DoorContainer struct {
	owner    *Character
	sequence uint32
	casts    map[constant.SkillID]*doorCast
}

type doorCast struct {
	buff *SkillBuff
	key  DoorKey
	town *Map
	door *Door
}

func NewDoorContainer(owner *Character) *DoorContainer {
	return &DoorContainer{
		owner: owner,
		casts: make(map[constant.SkillID]*doorCast),
	}
}

func (dc *DoorContainer) Spawn(buff *SkillBuff) {
	m := dc.owner.GetMap()
	if m == nil || m.Wz == nil {
		return
	}

	skillID := constant.SkillID(buff.Wz.ID)
	// The old town half must be removed before the new one is requested, so its portal slot is free again.
	if cast := dc.casts[skillID]; cast != nil {
		dc.despawn(skillID, cast, true)
	}
	if skillID != constant.SkillMysticDoor {
		return
	}

	dc.sequence++
	cast := &doorCast{
		buff: buff,
		key:  DoorKey(fmt.Sprintf("%d:%d", dc.owner.GetID(), dc.sequence)),
	}
	dc.casts[skillID] = cast
	cast.town = dc.owner.GameWorld.GetMapSystem().CreateReturnDoor(dc.owner, cast.key, skillID)
}

func (dc *DoorContainer) SpawnField(key DoorKey, skillID constant.SkillID, returnEp, fieldEp DoorEndpoint) *Door {
	cast := dc.casts[skillID]
	if cast == nil || cast.key != key || cast.door != nil {
		return nil
	}
	m := dc.owner.GetMap()
	if m == nil || m.Wz == nil || m != fieldEp.Map {
		return nil
	}
	ref, err := m.Reserve()
	if err != nil {
		return nil
	}

	door := NewDoor(key, dc.owner.GetID(), skillID, fieldEp, returnEp, dc.owner.partyID)
	door.fieldRef = ref
	cast.door = door
	m.AddDoor(door)
	dc.owner.Listener.OnPartyMemberFieldsChanged(dc.owner)
	return door
}

func (dc *DoorContainer) Remove(buff *SkillBuff, animated bool) {
	skillID := constant.SkillID(buff.Wz.ID)
	cast := dc.casts[skillID]
	if cast == nil || cast.buff != buff {
		return
	}
	dc.despawn(skillID, cast, animated)
}

func (dc *DoorContainer) despawn(skillID constant.SkillID, cast *doorCast, animated bool) {
	delete(dc.casts, skillID)
	dc.owner.GameWorld.GetMapSystem().DespawnDoor(cast.town, cast.key, animated, false)
	if cast.door == nil {
		return
	}

	dc.owner.GameWorld.GetMapSystem().DespawnDoor(cast.door.Field.Map, cast.key, animated, false)
	dc.owner.Listener.OnPartyMemberFieldsChanged(dc.owner)
}

func (dc *DoorContainer) forget(door *Door) {
	cast := dc.casts[door.SkillID]
	if cast == nil || cast.door != door {
		return
	}
	delete(dc.casts, door.SkillID)
}

func (dc *DoorContainer) Find(skillID constant.SkillID) *Door {
	cast := dc.casts[skillID]
	if cast == nil {
		return nil
	}
	return cast.door
}

func (dc *DoorContainer) SetPartyID(partyID *uint32) {
	for _, cast := range dc.casts {
		if cast.door != nil {
			cast.door.PartyID = partyID
		}
		dc.owner.GameWorld.GetMapSystem().SetDoorPartyID(cast.town, cast.key, partyID)
	}
}
