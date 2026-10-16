package game

import (
	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

// PartyMemberFromProto builds an entity PartyMember from protobuf.
func PartyMemberFromProto(pb *internal.PartyMember) *PartyMember {
	if pb == nil {
		return nil
	}
	m := &PartyMember{
		worldID:       pb.GetWorldId(),
		characterID:   pb.GetCharacterId(),
		characterName: pb.GetCharacterName(),
		level:         pb.GetLevel(),
		classID:       pb.GetClassId(),
		role:          pb.GetRole(),
		mapID:         pb.GetMapId(),
	}
	if pb.ChannelIndex != nil {
		c := pb.GetChannelIndex()
		m.channelIndex = &c
	}
	m.door = PartyDoorFromProto(pb.GetDoor())
	return m
}

// BuildRPC serializes this member to protobuf.
func (m *PartyMember) BuildRPC() *internal.PartyMember {
	if m == nil {
		return nil
	}
	pm := &internal.PartyMember{
		WorldId:       m.worldID,
		CharacterId:   m.characterID,
		CharacterName: m.characterName,
		Level:         m.level,
		ClassId:       m.classID,
		Role:          m.role,
		MapId:         m.mapID,
	}
	if m.channelIndex != nil {
		c := *m.channelIndex
		pm.ChannelIndex = &c
	}
	pm.Door = m.door.BuildRPC()
	return pm
}

// PartyMemberFromCharacter builds a PartyMember snapshot from an online character.
func PartyMemberFromCharacter(ch *Character, worldID uint32, channelID int32, role internal.PartyMemberRole) *PartyMember {
	if ch == nil {
		return nil
	}
	mapID := uint32(0)
	if m := ch.GetMap(); m != nil {
		mapID = m.TemplateID()
	}
	m := &PartyMember{
		worldID:       worldID,
		characterID:   ch.GetID(),
		characterName: ch.GetName(),
		level:         uint32(ch.GetLevel()),
		classID:       uint32(ch.class),
		role:          role,
		mapID:         mapID,
	}
	if channelID >= 0 {
		ci := channelID
		m.channelIndex = &ci
	}
	if door := ch.doors.Find(constant.SkillMysticDoor); door != nil {
		if pb := door.BuildRPC(); pb != nil {
			m.door = PartyDoorFromProto(pb)
		}
	}
	return m
}
