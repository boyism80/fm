package game

import (
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type Party struct {
	gameWorld         GameWorld
	worldID           uint32
	partyID           uint32
	leaderCharacterID uint32
	updatedAt         time.Time
	state             internal.PartyState
	members           []*PartyMember
}

func (p *Party) GetWorldId() uint32            { return p.worldID }
func (p *Party) GetPartyId() uint32            { return p.partyID }
func (p *Party) GetLeaderCharacterId() uint32  { return p.leaderCharacterID }
func (p *Party) GetState() internal.PartyState { return p.state }
func (p *Party) GetMembers() []*PartyMember    { return p.members }

func (p *Party) Clone() *Party {
	if p == nil {
		return nil
	}
	out := &Party{
		gameWorld:         p.gameWorld,
		worldID:           p.worldID,
		partyID:           p.partyID,
		leaderCharacterID: p.leaderCharacterID,
		updatedAt:         p.updatedAt,
		state:             p.state,
		members:           make([]*PartyMember, 0, len(p.members)),
	}
	for _, m := range p.members {
		out.members = append(out.members, m.Clone())
	}
	return out
}

func (p *Party) MemberCharacterIDs() []uint32 {
	out := make([]uint32, 0, len(p.members))
	for _, m := range p.members {
		if m == nil || m.characterID == 0 {
			continue
		}
		out = append(out, m.characterID)
	}
	return out
}
