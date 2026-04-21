package entity

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

// PartyMemberFromProto builds an entity PartyMember from protobuf.
func PartyMemberFromProto(pb *internal.PartyMember) *PartyMember {
	if pb == nil {
		return nil
	}
	m := &PartyMember{
		WorldID:       pb.GetWorldId(),
		CharacterID:   pb.GetCharacterId(),
		CharacterName: pb.GetCharacterName(),
		Level:         pb.GetLevel(),
		ClassID:       pb.GetClassId(),
		Role:          pb.GetRole(),
		MapID:         pb.GetMapId(),
	}
	if pb.ChannelIndex != nil {
		c := pb.GetChannelIndex()
		m.ChannelIndex = &c
	}
	m.Door = PartyDoorFromProto(pb.GetDoor())
	return m
}

// ToProto serializes this member to protobuf.
func (m *PartyMember) ToProto() *internal.PartyMember {
	if m == nil {
		return nil
	}
	pm := &internal.PartyMember{
		WorldId:       m.WorldID,
		CharacterId:   m.CharacterID,
		CharacterName: m.CharacterName,
		Level:         m.Level,
		ClassId:       m.ClassID,
		Role:          m.Role,
		MapId:         m.MapID,
	}
	if m.ChannelIndex != nil {
		c := *m.ChannelIndex
		pm.ChannelIndex = &c
	}
	pm.Door = m.Door.ToProto()
	return pm
}

// PartyMemberFromCharacter builds a PartyMember snapshot from an online character.
func PartyMemberFromCharacter(ch *Character, worldID uint32, channelID int32, role string) *PartyMember {
	if ch == nil {
		return nil
	}
	mapID := uint32(0)
	if m := ch.GetMap(); m != nil {
		mapID = m.GetMapID()
	}
	m := &PartyMember{
		WorldID:       worldID,
		CharacterID:   ch.GetID(),
		CharacterName: ch.GetName(),
		Level:         uint32(ch.GetLevel()),
		ClassID:       uint32(ch.Class),
		Role:          role,
		MapID:         mapID,
	}
	if channelID >= 0 {
		ci := channelID
		m.ChannelIndex = &ci
	}
	if door := ch.GetDoor(0); door != nil {
		if pb := door.ToProto(); pb != nil {
			m.Door = PartyDoorFromProto(pb)
		}
	}
	return m
}
