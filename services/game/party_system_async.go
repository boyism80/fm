package game

import (
	"context"
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

// CreateAsync creates a party led by ch and puts ch in it. A refusal tells ch why and fails the promise, so steps
// chained after it run only once ch leads the new party.
func (s *PartySystem) CreateAsync(ctx actor.Context, ch *Character) *async.Promise[*internal.CreatePartyReply] {
	req := &internal.CreatePartyRequest{
		WorldId: s.gw.GetWorldID(),
		Leader:  PartyMemberFromCharacter(ch, s.gw.GetWorldID(), s.gw.GetChannelID()).BuildRPC(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(func(c context.Context) (*internal.CreatePartyReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.CreateParty(c, req)
	}, func(reply *internal.CreatePartyReply) error {
		if !reply.GetOk() {
			ch.listener.OnPartyStatusMessage(ch, pconst.PartyStatusForInternalError(pconst.PartyC2SCreate, int32(reply.GetErrorCode())), "")
			return fmt.Errorf("create party: character=%d code=%v", ch.GetID(), reply.GetErrorCode())
		}
		if reply.PartyId == nil || reply.GetParty() == nil {
			return fmt.Errorf("create party: no party in the reply character=%d", ch.GetID())
		}
		s.Update(reply.GetParty())
		ch.party.Join(*reply.PartyId, reply.GetParty().GetRevision())
		ch.listener.OnPartyCreated(ch, *reply.PartyId)
		return nil
	})
}

func (s *PartySystem) UpdateMemberAsync(ctx actor.Context, ch *Character) *async.Task {
	p := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	if ch == nil {
		return p
	}
	if ch.party.ID() == nil {
		return p
	}
	member := PartyMemberFromCharacter(ch, uint32(s.gw.GetWorldID()), s.gw.GetChannelID())
	if member == nil {
		return p
	}

	s.mu.Lock()
	if _, sending := s.memberUpdates[ch.GetID()]; sending {
		s.memberUpdates[ch.GetID()] = member
		s.mu.Unlock()
		return p
	}
	s.memberUpdates[ch.GetID()] = nil
	s.mu.Unlock()
	return s.sendMemberUpdate(p, member)
}

func (s *PartySystem) sendMemberUpdate(p *async.Task, member *PartyMember) *async.Task {
	cid := member.GetCharacterId()
	req := &internal.UpdatePartyMemberRequest{Member: member.BuildRPC()}
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
