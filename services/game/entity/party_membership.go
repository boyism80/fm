package entity

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

func (p *PartyMembership) SetID(id *uint32) {
	p.id = id
	p.owner.Doors.SetPartyID(id)
}
