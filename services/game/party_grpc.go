package game

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
		gameWorld:         gw,
		worldID:           pb.GetWorldId(),
		partyID:           pb.GetPartyId(),
		leaderCharacterID: pb.GetLeaderCharacterId(),
		updatedAt:         time.UnixMilli(pb.GetUpdatedAtUnixMs()),
		state:             pb.GetState(),
		members:           members,
	}
}

func (p *Party) BuildRPC() *internal.Party {
	if p == nil {
		return nil
	}
	members := make([]*internal.PartyMember, 0, len(p.members))
	for _, m := range p.members {
		if pm := m.BuildRPC(); pm != nil {
			members = append(members, pm)
		}
	}
	return &internal.Party{
		WorldId:           p.worldID,
		PartyId:           p.partyID,
		LeaderCharacterId: p.leaderCharacterID,
		UpdatedAtUnixMs:   p.updatedAt.UnixMilli(),
		State:             p.state,
		Members:           members,
	}
}
