package domain

import (
	"context"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (s *PartySystem) UpdateMemberAsync(ctx actor.Context, ch *Character) *async.Task {
	p := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	if ch == nil {
		return p
	}
	if ch.Party.ID() == nil {
		return p
	}
	mm := ch.ToProtoPartyMember(uint32(s.gw.GetWorldID()), s.gw.GetChannelID(), internal.PartyMemberRole_PARTY_MEMBER_ROLE_MEMBER)
	if mm == nil {
		return p
	}
	req := &internal.UpdatePartyMemberRequest{Member: mm}

	s.mu.Lock()
	if _, sending := s.memberUpdates[ch.GetID()]; sending {
		s.memberUpdates[ch.GetID()] = req
		s.mu.Unlock()
		return p
	}
	s.memberUpdates[ch.GetID()] = nil
	s.mu.Unlock()
	return s.sendMemberUpdate(p, req)
}

func (s *PartySystem) sendMemberUpdate(p *async.Task, req *internal.UpdatePartyMemberRequest) *async.Task {
	cid := req.Member.GetCharacterId()
	p.OnError(func(err error) {
		log.Printf("UpdatePartyMember async char %d: %v", cid, err)
	})
	p.Finally(func() {
		s.mu.Lock()
		next := s.memberUpdates[cid]
		if next == nil {
			delete(s.memberUpdates, cid)
			s.mu.Unlock()
			return
		}
		s.memberUpdates[cid] = nil
		s.mu.Unlock()
		s.sendMemberUpdate(async.NewTask(nil, core.InternalRPCPerStepTimeout), next)
	})
	p.ThenRPC(func(c context.Context) (*internal.UpdatePartyMemberReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.UpdatePartyMember(c, req)
	}, nil)
	return p
}
