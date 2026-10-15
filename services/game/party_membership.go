package game

type PartyMembership struct {
	owner  *Character
	id     *uint32
	Search *PartySearchConfig
}

func (p *PartyMembership) ID() *uint32 {
	if p.id == nil {
		return nil
	}
	id := *p.id
	return &id
}

func (p *PartyMembership) Is(partyID *uint32) bool {
	return p.id != nil && partyID != nil && *p.id == *partyID
}

func (p *PartyMembership) SetID(id *uint32) {
	p.id = id
	p.owner.Doors.SetPartyID(id)
}

func (p *PartyMembership) MembersOnMap() ([]*Character, *Party) {
	partyID := p.ID()
	if partyID == nil || p.owner.GameWorld == nil {
		return nil, nil
	}
	party := p.owner.GameWorld.Systems().Get[*PartySystem]().Get(*partyID)
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
	if !p.Is(peer.Party.ID()) {
		return
	}
	p.owner.Listener.OnPartyMemberHP(p.owner, peer)
	peer.Listener.OnPartyMemberHP(peer, p.owner)
}
