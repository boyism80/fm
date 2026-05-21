package entity

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

// PartyMember is server-side party member state (decoupled from protobuf).
type PartyMember struct {
	WorldID       uint32
	CharacterID   uint32
	CharacterName string
	Level         uint32
	ClassID       uint32
	Role          internal.PartyMemberRole
	MapID         uint32
	ChannelIndex  *int32
	Door          *PartyDoor
}

func (m *PartyMember) GetCharacterId() uint32            { return m.CharacterID }
func (m *PartyMember) GetCharacterName() string          { return m.CharacterName }
func (m *PartyMember) GetLevel() uint32                  { return m.Level }
func (m *PartyMember) GetClassId() uint32                { return m.ClassID }
func (m *PartyMember) GetRole() internal.PartyMemberRole { return m.Role }
func (m *PartyMember) GetMapId() uint32                  { return m.MapID }
func (m *PartyMember) GetChannelIndex() *int32           { return m.ChannelIndex }
func (m *PartyMember) GetDoor() *PartyDoor               { return m.Door }

// Clone returns a deep copy of the member.
func (m *PartyMember) Clone() *PartyMember {
	if m == nil {
		return nil
	}
	out := &PartyMember{
		WorldID:       m.WorldID,
		CharacterID:   m.CharacterID,
		CharacterName: m.CharacterName,
		Level:         m.Level,
		ClassID:       m.ClassID,
		Role:          m.Role,
		MapID:         m.MapID,
	}
	if m.ChannelIndex != nil {
		c := *m.ChannelIndex
		out.ChannelIndex = &c
	}
	out.Door = m.Door.Clone()
	return out
}
