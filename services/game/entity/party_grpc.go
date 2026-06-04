package entity

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

// PartyFromProto builds an entity Party from an internal protobuf message.
func PartyFromProto(gw GameWorld, pb *internal.Party) *Party {
	if pb == nil {
		return nil
	}
	members := make([]*PartyMember, 0, len(pb.GetMembers()))
	for _, mm := range pb.GetMembers() {
		if m := PartyMemberFromProto(mm); m != nil {
			members = append(members, m)
		}
	}
	return &Party{
		GameWorld:         gw,
		WorldID:           pb.GetWorldId(),
		PartyID:           pb.GetPartyId(),
		LeaderCharacterID: pb.GetLeaderCharacterId(),
		Revision:          pb.GetRevision(),
		State:             pb.GetState(),
		Members:           members,
	}
}

// ToProto serializes this party to the internal protobuf message.
func (p *Party) ToProto() *internal.Party {
	if p == nil {
		return nil
	}
	members := make([]*internal.PartyMember, 0, len(p.Members))
	for _, m := range p.Members {
		if pm := m.ToProto(); pm != nil {
			members = append(members, pm)
		}
	}
	return &internal.Party{
		WorldId:           p.WorldID,
		PartyId:           p.PartyID,
		LeaderCharacterId: p.LeaderCharacterID,
		Revision:          p.Revision,
		State:             p.State,
		Members:           members,
	}
}
