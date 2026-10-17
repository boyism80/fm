package game

import (
	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/types"
)

func NewPlayerNpcFromInternalProto(pb *internal.PlayerNpc) *PlayerNpc {
	equips := make(equippedLook, len(pb.GetEquips()))
	for _, equip := range pb.GetEquips() {
		equips[constant.EquipmentPartsType(equip.GetSlot())] = equip.GetItemId()
	}
	return &PlayerNpc{
		NpcID:       pb.GetNpcId(),
		CharacterID: pb.GetCharacterId(),
		Name:        pb.GetName(),
		Gender:      uint8(pb.GetGender()),
		Skin:        uint8(pb.GetSkin()),
		Face:        pb.GetFace(),
		Hair:        pb.GetHair(),
		Equips:      equips,
		MapID:       pb.GetMapId(),
		Position:    types.Point[int16]{X: int16(pb.GetX()), Y: int16(pb.GetY())},
	}
}

func (n *PlayerNpc) BuildRPC() *internal.PlayerNpc {
	equips := make([]*internal.PlayerNpcEquip, 0, len(n.Equips))
	for parts, itemID := range n.Equips {
		equips = append(equips, &internal.PlayerNpcEquip{Slot: int32(parts), ItemId: itemID})
	}
	return &internal.PlayerNpc{
		NpcId:       n.NpcID,
		CharacterId: n.CharacterID,
		Name:        n.Name,
		Gender:      uint32(n.Gender),
		Skin:        uint32(n.Skin),
		Face:        n.Face,
		Hair:        n.Hair,
		Equips:      equips,
		MapId:       n.MapID,
		X:           int32(n.Position.X),
		Y:           int32(n.Position.Y),
	}
}
