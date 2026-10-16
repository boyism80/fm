package game

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

type GuildMember struct {
	worldID       uint32
	characterID   uint32
	characterName string
	level         uint32
	classID       uint32
	rank          internal.GuildMemberRank
	allianceRank  *uint32
	channelIndex  *int32
}

func (m *GuildMember) GetAllianceRank() (uint32, bool) {
	if m == nil || m.allianceRank == nil {
		return 0, false
	}
	return *m.allianceRank, true
}

func (m *GuildMember) GetCharacterId() uint32 {
	return m.characterID
}

func (m *GuildMember) GetCharacterName() string {
	return m.characterName
}

func (m *GuildMember) GetLevel() uint32 {
	return m.level
}

func (m *GuildMember) GetClassId() uint32 {
	return m.classID
}

func (m *GuildMember) GetRank() internal.GuildMemberRank {
	return m.rank
}

func (m *GuildMember) GetChannelIndex() *int32 {
	return m.channelIndex
}

func (m *GuildMember) Clone() *GuildMember {
	if m == nil {
		return nil
	}
	out := &GuildMember{
		worldID:       m.worldID,
		characterID:   m.characterID,
		characterName: m.characterName,
		level:         m.level,
		classID:       m.classID,
		rank:          m.rank,
	}
	if m.allianceRank != nil {
		ar := *m.allianceRank
		out.allianceRank = &ar
	}
	if m.channelIndex != nil {
		c := *m.channelIndex
		out.channelIndex = &c
	}
	return out
}

func GuildMemberFromProto(pb *internal.GuildMember) *GuildMember {
	if pb == nil {
		return nil
	}
	m := &GuildMember{
		worldID:       pb.GetWorldId(),
		characterID:   pb.GetCharacterId(),
		characterName: pb.GetCharacterName(),
		level:         pb.GetLevel(),
		classID:       pb.GetClassId(),
		rank:          pb.GetRank(),
	}
	if pb.ChannelIndex != nil {
		c := pb.GetChannelIndex()
		m.channelIndex = &c
	}
	if pb.AllianceRank != nil {
		ar := pb.GetAllianceRank()
		m.allianceRank = &ar
	}
	return m
}

func (m *GuildMember) BuildRPC() *internal.GuildMember {
	if m == nil {
		return nil
	}
	pm := &internal.GuildMember{
		WorldId:       m.worldID,
		CharacterId:   m.characterID,
		CharacterName: m.characterName,
		Level:         m.level,
		ClassId:       m.classID,
		Rank:          m.rank,
	}
	if m.channelIndex != nil {
		c := *m.channelIndex
		pm.ChannelIndex = &c
	}
	if m.allianceRank != nil {
		ar := *m.allianceRank
		pm.AllianceRank = &ar
	}
	return pm
}

func GuildMemberFromCharacter(ch *Character, worldID uint32, channelID int32, rank internal.GuildMemberRank) *GuildMember {
	if ch == nil {
		return nil
	}
	m := &GuildMember{
		worldID:       worldID,
		characterID:   ch.GetID(),
		characterName: ch.GetName(),
		level:         uint32(ch.GetLevel()),
		classID:       uint32(ch.class),
		rank:          rank,
	}
	if channelID >= 0 {
		ci := channelID
		m.channelIndex = &ci
	}
	return m
}
