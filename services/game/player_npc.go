package game

import (
	"errors"
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
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

	position := at.GetPosition()
	if seat, ok := m.wz.PlayerNpcSeats[npcID]; ok {
		position = seatPosition(seat)
	}
	return newPlayerNpc(source, npcID, m.TemplateID(), position), nil
}

func newPlayerNpc(source *Character, npcID uint32, mapID uint32, position types.Point[int16]) *PlayerNpc {
	return &PlayerNpc{
		NpcID:       npcID,
		CharacterID: source.GetID(),
		Name:        source.name,
		Gender:      source.look.gender,
		Skin:        source.look.skinColor,
		Face:        source.look.face,
		Hair:        source.look.hair,
		Equips:      source.inventory.look(),
		MapID:       mapID,
		Position:    position,
	}
}

// seatPosition is where a seat's npc stands: its x on the ground under it.
func seatPosition(seat wz.NpcSpawn) types.Point[int16] {
	return types.Point[int16]{X: seat.Position.X, Y: seat.CollisionY}
}

// hallOfFame are the job instructors' maps whose seats take a level 200 explorer's player npc, as the original did.
var hallOfFame = []struct {
	class constant.ClassType
	mapID uint32
}{
	{constant.ClassWarrior, 102000003},
	{constant.ClassMagician, 101000003},
	{constant.ClassBowman, 100000201},
	{constant.ClassThief, 103000003},
	{constant.ClassPirate, 120000000},
}

// celebrateMaxLevel tells the world ch reached level 200 and stands ch's player npc on the first free seat of its job
// instructor's map.
func (ch *Character) celebrateMaxLevel() {
	ch.Post(func(actx actor.Context) {
		medal := ""
		if item := ch.inventory.equipped[constant.EquipmentPartsMedal]; item != nil {
			medal = "<" + ch.gameWorld.GetResources().GetItemName(item.GetModel().GetID()) + "> "
		}
		message := fmt.Sprintf(constant.MaxLevelNoticeMessage, medal, ch.GetName())
		ch.gameWorld.Systems().Get[*InternalSystem]().BroadcastNoticeAsync(actx, ch, constant.MsgLightBlueText, message, false).OnError(func(err error) {
			log.Printf("level 200 notice character=%d: %v", ch.GetID(), err)
		})

		for _, hall := range hallOfFame {
			if !ch.class.Of(hall.class) {
				continue
			}
			npcs := ch.gameWorld.Systems().Get[*PlayerNpcSystem]()
			seat, ok := npcs.FreeSeat(hall.mapID)
			if !ok {
				return
			}
			npcs.CreateAsync(actx, newPlayerNpc(ch, seat.ID, hall.mapID, seatPosition(seat))).Bind(ch).Do(func(created bool) error {
				if created {
					ch.listener.OnMessage(ch, constant.MsgPopup, constant.HallOfFameNpcMessage)
					ch.listener.OnMessage(ch, constant.MsgPinkText, constant.HallOfFameNpcMessage)
				}
				return nil
			}).OnError(func(err error) {
				log.Printf("level 200 player npc character=%d: %v", ch.GetID(), err)
			})
			return
		}
	})
}

func (n *PlayerNpc) BuildS2C() *dto.Character {
	base, overlays, cashWeapon := n.Equips.split()
	return &dto.Character{
		Name:      n.Name,
		Gender:    n.Gender,
		SkinColor: n.Skin,
		Face:      n.Face,
		Hair:      n.Hair,
		BaseLooks: base,
		Overlays:  overlays,
		Weapon:    cashWeapon,
	}
}

type equippedLook map[constant.EquipmentPartsType]uint32

// split follows the client: cash items (slot -101 and below) cover the regular item in the same slot, which is
// kept as an overlay. The cash weapon (-111) is sent apart in the weapon sticker field, 0 when there is none.
func (e equippedLook) split() (base map[int8]uint32, overlays map[int8]uint32, cashWeapon uint32) {
	base = make(map[int8]uint32)
	overlays = make(map[int8]uint32)
	for parts, itemID := range e {
		if parts < -127 {
			continue
		}
		absoluteParts := int8(parts * -1)
		_, exists := base[absoluteParts]
		if absoluteParts < 100 && !exists {
			base[absoluteParts] = itemID
		} else if absoluteParts > 100 && absoluteParts != 111 {
			adjustedParts := absoluteParts - 100
			if existing, exists := base[adjustedParts]; exists {
				overlays[adjustedParts] = existing
			}
			base[adjustedParts] = itemID
		} else if exists {
			overlays[absoluteParts] = itemID
		}
	}
	return base, overlays, e[constant.EquipmentPartsWeapon.Cash()]
}
