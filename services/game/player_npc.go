package game

import (
	"errors"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/types"
)

var (
	ErrPlayerNpcUnknownTemplate = errors.New("unknown npc template")
	ErrPlayerNpcTaken           = errors.New("npc template already used by a player npc")
	ErrPlayerNpcPlacedOnMap     = errors.New("npc template already placed on this map")
)

type PlayerNpc struct {
	NpcID       uint32
	CharacterID uint32
	Name        string
	Gender      uint8
	Skin        uint8
	Face        uint32
	Hair        uint32
	Equips      equippedLook
	MapID       uint32
	Position    types.Point[int16]
}

// NewPlayerNpc copies source's look onto the npc template and stands it where at is.
func NewPlayerNpc(source *Character, npcID uint32, at *Character) (*PlayerNpc, error) {
	if source.gameWorld.GetResources().Strings.NpcStrings[npcID] == nil {
		return nil, ErrPlayerNpcUnknownTemplate
	}
	if source.gameWorld.Systems().Get[*PlayerNpcSystem]().Find(npcID) != nil {
		return nil, ErrPlayerNpcTaken
	}
	m := at.GetMap()
	if m.NpcByTemplate(npcID) != nil {
		return nil, ErrPlayerNpcPlacedOnMap
	}

	return &PlayerNpc{
		NpcID:       npcID,
		CharacterID: source.GetID(),
		Name:        source.name,
		Gender:      source.look.gender,
		Skin:        source.look.skinColor,
		Face:        source.look.face,
		Hair:        source.look.hair,
		Equips:      source.inventory.look(),
		MapID:       m.TemplateID(),
		Position:    at.GetPosition(),
	}, nil
}

func (n *PlayerNpc) BuildS2C() *dto.Character {
	base, overlays := n.Equips.split()
	return &dto.Character{
		Name:      n.Name,
		Gender:    n.Gender,
		SkinColor: n.Skin,
		Face:      n.Face,
		Hair:      n.Hair,
		BaseLooks: base,
		Overlays:  overlays,
		Weapon:    n.Equips[constant.EquipmentPartsWeapon],
	}
}

type equippedLook map[constant.EquipmentPartsType]uint32

// split follows the client: cash items (slot -101 and below) cover the regular item in the same slot, which is
// kept as an overlay.
func (e equippedLook) split() (base map[int8]uint32, overlays map[int8]uint32) {
	base = make(map[int8]uint32)
	overlays = make(map[int8]uint32)
	for parts, itemID := range e {
		if parts < -127 {
			continue
		}
		absoluteParts := int8(parts * -1)
		if absoluteParts < 100 {
			if _, exists := base[absoluteParts]; !exists {
				base[absoluteParts] = itemID
			}
		} else if absoluteParts > 100 && absoluteParts != 111 {
			adjustedParts := absoluteParts - 100
			if existing, exists := base[adjustedParts]; exists {
				overlays[adjustedParts] = existing
			}
			base[adjustedParts] = itemID
		} else if _, exists := base[absoluteParts]; exists {
			overlays[absoluteParts] = itemID
		}
	}
	return base, overlays
}
