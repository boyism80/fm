package game

type PartyMembership struct {
	owner  *Character
	id     *uint32
	search *PartySearchConfig
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
	p.owner.doors.SetPartyID(id)
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
