package game

import "github.com/boyism80/fm/protocol/response"

func (p *Party) BuildS2C() []response.PartyMemberStatus {
	out := make([]response.PartyMemberStatus, 0, len(p.Members))
	for _, m := range p.Members {
		if m == nil {
			continue
		}
		ch := int32(-2)
		if m.ChannelIndex != nil {
			ch = *m.ChannelIndex
		}
		doorTown := uint32(999999999)
		doorTarget := uint32(999999999)
		doorX := int32(0)
		doorY := int32(0)
		if d := m.Door; d != nil {
			doorTown = d.Town
			doorTarget = d.Target
			doorX = d.X
			doorY = d.Y
		}
		out = append(out, response.PartyMemberStatus{
			CharacterID: m.CharacterID,
			Name:        m.CharacterName,
			Class:       m.ClassID,
			Level:       m.Level,
			Channel:     ch,
			MapID:       m.MapID,
			DoorTown:    doorTown,
			DoorTarget:  doorTarget,
			DoorX:       doorX,
			DoorY:       doorY,
		})
	}
	return out
}
