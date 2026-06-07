package entity

import (
	"log"

	"github.com/boyism80/fm/services/game/constant"
)

func (m *Map) registerFieldDropTimers(fp *FieldPlacement, dropType constant.DropType) {
	if fp == nil || m == nil || m.Wz == nil {
		return
	}
	if m.Wz.Everlast && fp.PlayerDrop {
		return
	}

	fp.RegisterExpire(constant.ItemExpireTime)
	if dropType == constant.DropTypeOwnerOnly || dropType == constant.DropTypeParty {
		fp.RegisterFFA(constant.ItemFFATime)
	}
}

func (m *Map) canLootFieldDrop(fp *FieldPlacement, character *Character) bool {
	if fp == nil || character == nil || m == nil || m.Wz == nil {
		return false
	}
	if fp.Owner != character.GetID() {
		if (!fp.PlayerDrop && fp.DropType == constant.DropTypeOwnerOnly) ||
			(fp.PlayerDrop && m.Wz.Everlast) {
			return false
		}
		if !fp.PlayerDrop && fp.DropType == constant.DropTypeParty && !m.isDropOwnerInParty(character, fp.Owner) {
			return false
		}
	}
	return true
}

func (m *Map) isDropOwnerInParty(character *Character, ownerID uint32) bool {
	if character == nil || ownerID == 0 {
		return false
	}
	if character.GetID() == ownerID {
		return true
	}
	partyID := character.GetPartyID()
	if partyID == nil || m.GameWorld == nil {
		return false
	}
	party := m.GameWorld.GetPartySystem().Get(*partyID)
	if party == nil {
		return false
	}
	for _, member := range party.Members {
		if member != nil && member.CharacterID == ownerID {
			return true
		}
	}
	return false
}

func (m *Map) collectOwnedFieldDrops(character *Character) {
	if m == nil || character == nil || m.Wz == nil || !m.Wz.Everlast {
		return
	}
	if m.objects[constant.ObjectTypeItem] == nil {
		return
	}

	removeIDs := make([]uint32, 0)

	for oid, obj := range m.objects[constant.ObjectTypeItem] {
		switch drop := obj.(type) {
		case Item:
			fp := drop.GetFieldPlacement()
			if fp == nil || fp.Owner != character.GetID() {
				continue
			}
			if _, err := character.AddItem(drop, false); err != nil {
				log.Printf("collectOwnedFieldDropsOnLeave: failed to return item %d: %v", oid, err)
				continue
			}
			removeIDs = append(removeIDs, oid)
		case *Meso:
			fp := drop.GetFieldPlacement()
			if fp == nil || fp.Owner != character.GetID() {
				continue
			}
			character.GainMeso(drop.GetCount32())
			removeIDs = append(removeIDs, oid)
		}
	}

	for _, oid := range removeIDs {
		_ = m.RemoveItem(oid, constant.RemoveItemTypeAnimated, character.GetID())
	}
}
