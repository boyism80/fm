package entity

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

// Party is server-side party aggregate (decoupled from protobuf).
type Party struct {
	GameWorld         GameWorld
	WorldID           uint32
	PartyID           uint32
	LeaderCharacterID uint32
	Revision          uint64
	State             internal.PartyState
	Members           []*PartyMember
}

func (p *Party) GetWorldId() uint32            { return p.WorldID }
func (p *Party) GetPartyId() uint32            { return p.PartyID }
func (p *Party) GetLeaderCharacterId() uint32  { return p.LeaderCharacterID }
func (p *Party) GetRevision() uint64           { return p.Revision }
func (p *Party) GetState() internal.PartyState { return p.State }
func (p *Party) GetMembers() []*PartyMember    { return p.Members }

// Clone returns a deep copy suitable for handing to Lua or other consumers.
func (p *Party) Clone() *Party {
	if p == nil {
		return nil
	}
	out := &Party{
		GameWorld:         p.GameWorld,
		WorldID:           p.WorldID,
		PartyID:           p.PartyID,
		LeaderCharacterID: p.LeaderCharacterID,
		Revision:          p.Revision,
		State:             p.State,
		Members:           make([]*PartyMember, 0, len(p.Members)),
	}
	for _, m := range p.Members {
		out.Members = append(out.Members, m.Clone())
	}
	return out
}

// PartyMemberCharacterIDs returns non-zero character ids from members (order preserved).
func PartyMemberCharacterIDs(members []*PartyMember) []uint32 {
	out := make([]uint32, 0, len(members))
	for _, m := range members {
		if m == nil || m.CharacterID == 0 {
			continue
		}
		out = append(out, m.CharacterID)
	}
	return out
}
