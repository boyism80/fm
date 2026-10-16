package game

import (
	"log"

	"github.com/boyism80/fm/maple/constant"
)

func (m *Map) registerFieldDropTimers(fp *FieldPlacement, dropType constant.DropType) {
	if fp == nil || m == nil || m.wz == nil {
		return
	}
	if m.wz.Everlast && fp.playerDrop {
		return
	}

	fp.RegisterExpire(constant.ItemExpireTime)
	if dropType == constant.DropTypeOwnerOnly || dropType == constant.DropTypeParty {
		fp.RegisterFFA(constant.ItemFFATime)
	}
}

func (m *Map) canLootFieldDrop(fp *FieldPlacement, character *Character) bool {
	if fp == nil || character == nil || m == nil || m.wz == nil {
		return false
	}
	if fp.restrictedPlayerDrop && fp.owner != character.GetID() {
		return false
	}
	if fp.owner != character.GetID() {
		if (!fp.playerDrop && fp.dropType == constant.DropTypeOwnerOnly) ||
			(fp.playerDrop && m.wz.Everlast) {
			return false
		}
		if !fp.playerDrop && fp.dropType == constant.DropTypeParty && !m.isDropOwnerInParty(character, fp.owner) {
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
	partyID := character.party.ID()
	if partyID == nil || m.gameWorld == nil {
		return false
	}
	party := m.gameWorld.Systems().Get[*PartySystem]().Get(*partyID)
	if party == nil {
		return false
	}
	for _, member := range party.members {
		if member != nil && member.characterID == ownerID {
			return true
		}
	}
	return false
}

func (m *Map) collectOwnedFieldDrops(character *Character) {
	if m == nil || character == nil || m.wz == nil || !m.wz.Everlast {
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
			if fp == nil || fp.owner != character.GetID() {
				continue
			}
			if _, err := character.inventory.AddItem(drop, true); err != nil {
				log.Printf("collectOwnedFieldDrops: inventory full, leaving drop %d on map: %v", oid, err)
				continue
			}
			removeIDs = append(removeIDs, oid)
		case *Meso:
			fp := drop.GetFieldPlacement()
			if fp == nil || fp.owner != character.GetID() {
				continue
			}
			before := character.inventory.meso
			character.inventory.GainMeso(drop.GetCount32())
			if character.inventory.meso == before {
				log.Printf("collectOwnedFieldDrops: meso full, leaving drop %d on map", oid)
				continue
			}
			removeIDs = append(removeIDs, oid)
		}
	}

	for _, oid := range removeIDs {
		_ = m.RemoveItem(oid, constant.RemoveItemTypeAnimated, character.GetID())
	}
}
