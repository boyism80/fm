package entity

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

type GuildMember struct {
	WorldID       uint32
	CharacterID   uint32
	CharacterName string
	Level         uint32
	ClassID       uint32
	Rank          internal.GuildMemberRank
	ChannelIndex  *int32
}

func (m *GuildMember) GetCharacterId() uint32 {
	return m.CharacterID
}

func (m *GuildMember) GetCharacterName() string {
	return m.CharacterName
}

func (m *GuildMember) GetLevel() uint32 {
	return m.Level
}

func (m *GuildMember) GetClassId() uint32 {
	return m.ClassID
}

func (m *GuildMember) GetRank() internal.GuildMemberRank {
	return m.Rank
}

func (m *GuildMember) GetChannelIndex() *int32 {
	return m.ChannelIndex
}

func (m *GuildMember) Clone() *GuildMember {
	if m == nil {
		return nil
	}
	out := &GuildMember{
		WorldID:       m.WorldID,
		CharacterID:   m.CharacterID,
		CharacterName: m.CharacterName,
		Level:         m.Level,
		ClassID:       m.ClassID,
		Rank:          m.Rank,
	}
	if m.ChannelIndex != nil {
		c := *m.ChannelIndex
		out.ChannelIndex = &c
	}
	return out
}

func GuildMemberFromProto(pb *internal.GuildMember) *GuildMember {
	if pb == nil {
		return nil
	}
	m := &GuildMember{
		WorldID:       pb.GetWorldId(),
		CharacterID:   pb.GetCharacterId(),
		CharacterName: pb.GetCharacterName(),
		Level:         pb.GetLevel(),
		ClassID:       pb.GetClassId(),
		Rank:          pb.GetRank(),
	}
	if pb.ChannelIndex != nil {
		c := pb.GetChannelIndex()
		m.ChannelIndex = &c
	}
	return m
}

func (m *GuildMember) ToProto() *internal.GuildMember {
	if m == nil {
		return nil
	}
	pm := &internal.GuildMember{
		WorldId:       m.WorldID,
		CharacterId:   m.CharacterID,
		CharacterName: m.CharacterName,
		Level:         m.Level,
		ClassId:       m.ClassID,
		Rank:          m.Rank,
	}
	if m.ChannelIndex != nil {
		c := *m.ChannelIndex
		pm.ChannelIndex = &c
	}
	return pm
}

func GuildMemberFromCharacter(ch *Character, worldID uint32, channelID int32, rank internal.GuildMemberRank) *GuildMember {
	if ch == nil {
		return nil
	}
	m := &GuildMember{
		WorldID:       worldID,
		CharacterID:   ch.GetID(),
		CharacterName: ch.GetName(),
		Level:         uint32(ch.GetLevel()),
		ClassID:       uint32(ch.Class),
		Rank:          rank,
	}
	if channelID >= 0 {
		ci := channelID
		m.ChannelIndex = &ci
	}
	return m
}
