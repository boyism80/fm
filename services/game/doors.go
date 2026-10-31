package game

import (
	"fmt"

	"github.com/boyism80/fm/maple/constant"
)

type Doors struct {
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

func NewDoors(owner *Character) *Doors {
	return &Doors{
		owner: owner,
		casts: make(map[constant.SkillID]*doorCast),
	}
}

func (dc *Doors) Spawn(buff *SkillBuff) {
	m := dc.owner.GetMap()
	if m == nil || m.wz == nil {
		return
	}

	skillID := constant.SkillID(buff.wz.ID)
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
	cast.town = dc.owner.gameWorld.Systems().Get[*MapSystem]().CreateReturnDoor(dc.owner, cast.key, skillID)
}

func (dc *Doors) SpawnField(key DoorKey, skillID constant.SkillID, returnEp, fieldEp DoorEndpoint) *Door {
	cast := dc.casts[skillID]
	if cast == nil || cast.key != key || cast.door != nil {
		return nil
	}
	m := dc.owner.GetMap()
	if m == nil || m.wz == nil || m != fieldEp.Map {
		return nil
	}
	ref, err := m.Reserve()
	if err != nil {
		return nil
	}

	door := NewDoor(key, dc.owner.GetID(), skillID, fieldEp, returnEp, dc.owner.party.ID())
	door.fieldRef = ref
	cast.door = door
	m.AddDoor(door)
	dc.owner.gameWorld.Systems().Get[*PartySystem]().UpdateMemberAsync(nil, dc.owner)
	return door
}

func (dc *Doors) Remove(buff *SkillBuff, animated bool) {
	skillID := constant.SkillID(buff.wz.ID)
	cast := dc.casts[skillID]
	if cast == nil || cast.buff != buff {
		return
	}
	dc.despawn(skillID, cast, animated)
}

func (dc *Doors) despawn(skillID constant.SkillID, cast *doorCast, animated bool) {
	delete(dc.casts, skillID)
	dc.owner.gameWorld.Systems().Get[*MapSystem]().DespawnDoor(cast.town, cast.key, animated, false)
	if cast.door == nil {
		return
	}

	dc.owner.gameWorld.Systems().Get[*MapSystem]().DespawnDoor(cast.door.field.Map, cast.key, animated, false)
	dc.owner.gameWorld.Systems().Get[*PartySystem]().UpdateMemberAsync(nil, dc.owner)
}

func (dc *Doors) forget(door *Door) {
	cast := dc.casts[door.skillID]
	if cast == nil || cast.door != door {
		return
	}
	delete(dc.casts, door.skillID)
}

func (dc *Doors) Find(skillID constant.SkillID) *Door {
	cast := dc.casts[skillID]
	if cast == nil {
		return nil
	}
	return cast.door
}

func (dc *Doors) SetPartyID(partyID *uint32) {
	maps := dc.owner.gameWorld.Systems().Get[*MapSystem]()
	for _, cast := range dc.casts {
		if cast.door != nil {
			maps.SetDoorPartyID(cast.door.field.Map, cast.key, partyID)
		}
		maps.SetDoorPartyID(cast.town, cast.key, partyID)
	}
}
