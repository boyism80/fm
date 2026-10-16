package game

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

// PartyMember is server-side party member state (decoupled from protobuf).
type PartyMember struct {
	worldID       uint32
	characterID   uint32
	characterName string
	level         uint32
	classID       uint32
	role          internal.PartyMemberRole
	mapID         uint32
	channelIndex  *int32
	door          *PartyDoor
}

func (m *PartyMember) GetCharacterId() uint32            { return m.characterID }
func (m *PartyMember) GetCharacterName() string          { return m.characterName }
func (m *PartyMember) GetLevel() uint32                  { return m.level }
func (m *PartyMember) GetClassId() uint32                { return m.classID }
func (m *PartyMember) GetRole() internal.PartyMemberRole { return m.role }
func (m *PartyMember) GetMapId() uint32                  { return m.mapID }
func (m *PartyMember) GetChannelIndex() *int32           { return m.channelIndex }
func (m *PartyMember) GetDoor() *PartyDoor               { return m.door }

// Clone returns a deep copy of the member.
func (m *PartyMember) Clone() *PartyMember {
	if m == nil {
		return nil
	}
	out := &PartyMember{
		worldID:       m.worldID,
		characterID:   m.characterID,
		characterName: m.characterName,
		level:         m.level,
		classID:       m.classID,
		role:          m.role,
		mapID:         m.mapID,
	}
	if m.channelIndex != nil {
		c := *m.channelIndex
		out.channelIndex = &c
	}
	out.door = m.door.Clone()
	return out
}
