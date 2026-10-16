package game

import "github.com/boyism80/fm/protocol/response"

func (p *Party) BuildS2C() []response.PartyMemberStatus {
	out := make([]response.PartyMemberStatus, 0, len(p.members))
	for _, m := range p.members {
		if m == nil {
			continue
		}
		ch := int32(-2)
		if m.channelIndex != nil {
			ch = *m.channelIndex
		}
		doorTown := uint32(999999999)
		doorTarget := uint32(999999999)
		doorX := int32(0)
		doorY := int32(0)
		if d := m.door; d != nil {
			doorTown = d.town
			doorTarget = d.target
			doorX = d.x
			doorY = d.y
		}
		out = append(out, response.PartyMemberStatus{
			CharacterID: m.characterID,
			Name:        m.characterName,
			Class:       m.classID,
			Level:       m.level,
			Channel:     ch,
			MapID:       m.mapID,
			DoorTown:    doorTown,
			DoorTarget:  doorTarget,
			DoorX:       doorX,
			DoorY:       doorY,
		})
	}
	return out
}
