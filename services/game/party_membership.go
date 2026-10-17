package game

// PartyMembership keeps the last party and its revision after leaving, so a late snapshot
// of that party cannot bring the character back. A different party always replaces it.
type PartyMembership struct {
	owner    *Character
	id       uint32
	joined   bool
	revision uint64
	search   *PartySearchConfig
}

func (p *PartyMembership) ID() *uint32 {
	if !p.joined {
		return nil
	}
	id := p.id
	return &id
}

func (p *PartyMembership) Is(partyID *uint32) bool {
	return p.joined && partyID != nil && p.id == *partyID
}

func (p *PartyMembership) Join(partyID uint32, revision uint64) {
	if partyID == p.id && revision < p.revision {
		return
	}
	if p.joined && partyID == p.id {
		p.revision = revision
		return
	}
	p.id = partyID
	p.joined = true
	p.revision = revision
	p.owner.doors.SetPartyID(p.ID())
}

func (p *PartyMembership) Leave(partyID uint32, revision uint64) {
	if !p.joined || partyID != p.id || revision < p.revision {
		return
	}
	p.joined = false
	p.revision = revision
	p.owner.doors.SetPartyID(nil)
}

func (p *PartyMembership) MembersOnMap() ([]*Character, *Party) {
	partyID := p.ID()
	if partyID == nil || p.owner.gameWorld == nil {
		return nil, nil
	}
	party := p.owner.gameWorld.Systems().Get[*PartySystem]().Get(*partyID)
	leaderMap := p.owner.GetMap()
	if party == nil || party.GetLeaderCharacterId() != p.owner.GetID() || leaderMap == nil {
		return nil, nil
	}
	members := make([]*Character, 0)
	for _, mem := range party.GetMembers() {
		if mem == nil {
			continue
		}
		member := leaderMap.GetPlayer(mem.GetCharacterId())
		if member == nil {
			return nil, nil
		}
		members = append(members, member)
	}
	if len(members) == 0 {
		return nil, nil
	}
	return members, party
}

func (p *PartyMembership) syncHP(peer *Character) {
	if !p.Is(peer.party.ID()) {
		return
	}
	p.owner.listener.OnPartyMemberHP(p.owner, peer)
	peer.listener.OnPartyMemberHP(peer, p.owner)
}
