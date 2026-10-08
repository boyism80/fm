package server

import (
	"context"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/entity"
)

func (s partySystem) UpdateMemberAsync(ctx actor.Context, ch *entity.Character) *async.Task {
	p := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	if s.gs == nil || ch == nil || s.gs.internalClient == nil {
		return p
	}
	if ch.GetPartyID() == nil {
		return p
	}
	mm := ch.ToProtoPartyMember(uint32(s.gs.config.WorldId), int32(s.gs.config.ChannelId), internal.PartyMemberRole_PARTY_MEMBER_ROLE_MEMBER)
	if mm == nil {
		return p
	}
	req := &internal.UpdatePartyMemberRequest{Member: mm}

	pc := s.gs.party
	pc.mu.Lock()
	if _, sending := pc.memberUpdates[ch.GetID()]; sending {
		pc.memberUpdates[ch.GetID()] = req
		pc.mu.Unlock()
		return p
	}
	pc.memberUpdates[ch.GetID()] = nil
	pc.mu.Unlock()
	return pc.sendMemberUpdate(p, req)
}

func (pc *PartyCache) sendMemberUpdate(p *async.Task, req *internal.UpdatePartyMemberRequest) *async.Task {
	cid := req.Member.GetCharacterId()
	p.OnError(func(err error) {
		log.Printf("UpdatePartyMember async char %d: %v", cid, err)
	})
	p.Finally(func() {
		pc.mu.Lock()
		next := pc.memberUpdates[cid]
		if next == nil {
			delete(pc.memberUpdates, cid)
			pc.mu.Unlock()
			return
		}
		pc.memberUpdates[cid] = nil
		pc.mu.Unlock()
		pc.sendMemberUpdate(async.NewTask(nil, core.InternalRPCPerStepTimeout), next)
	})
	p.ThenRPC(func(c context.Context) (*internal.UpdatePartyMemberReply, error) {
		return pc.internalClient.UpdatePartyMember(c, req)
	}, nil)
	return p
}
