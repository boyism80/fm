package entity

import (
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

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
		UpdatedAt:         time.UnixMilli(pb.GetUpdatedAtUnixMs()),
		State:             pb.GetState(),
		Members:           members,
	}
}

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
		UpdatedAtUnixMs:   p.UpdatedAt.UnixMilli(),
		State:             p.State,
		Members:           members,
	}
}
