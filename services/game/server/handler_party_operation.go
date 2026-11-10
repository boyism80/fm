package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/rpc"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type PartyOperation struct {
	gs *GameServer
}

func (PartyOperation) New(gs *GameServer) *PartyOperation {
	return &PartyOperation{
		gs: gs,
	}
}

func (h *PartyOperation) Handle(ctx *core.ClientContext, req *request.PartyOperation) error {
	if ctx.ActorContext == nil {
		return fmt.Errorf("party operation: actor context required")
	}

	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}

	switch req.Operation {
	case constant.PartyC2SCreate:
		return h.create(ctx, ch, req)
	case constant.PartyC2SLeave:
		return h.leave(ctx, ch, req)
	case constant.PartyC2SAcceptInvite:
		return h.acceptInvite(ctx, ch, req)
	case constant.PartyC2SChangeLeader:
		return h.changeLeader(ctx, ch, req)
	case constant.PartyC2SInvite:
		return h.invite(ctx, ch, req)
	case constant.PartyC2SExpel:
		return h.expel(ctx, ch, req)
	default:
		return nil
	}
}

func (h *PartyOperation) create(ctx *core.ClientContext, ch *game.Character, req *request.PartyOperation) error {
	h.gs.systems.Get[*game.PartySystem]().CreateAsync(ctx.ActorContext, ch).OnError(func(err error) {
		log.Printf("PartyOperation(create): %v", err)
	})
	return nil
}

func (h *PartyOperation) leave(ctx *core.ClientContext, ch *game.Character, req *request.PartyOperation) error {
	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(
		func(c context.Context) (*internal.LeavePartyReply, error) {
			return h.gs.systems.Get[*rpc.InternalConnection]().Client.LeaveParty(c, &internal.LeavePartyRequest{
				WorldId:     worldID,
				CharacterId: charID,
			})
		},
		func(reply *internal.LeavePartyReply) error {
			if !reply.GetOk() {
				ch.Listener().OnPartyStatusMessage(ch, constant.PartyStatusForInternalError(constant.PartyC2SLeave, int32(reply.GetErrorCode())), "")
				log.Printf("PartyOperation(leave): failed character=%d code=%v", charID, reply.GetErrorCode())
				return nil
			}
			if reply.PartyId != nil {
				ch.Party().Leave(*reply.PartyId, reply.GetRevision())
			}
			return nil
		},
	).OnError(func(err error) {
		log.Printf("PartyOperation(leave) async error: %v", err)
	})
	return nil
}

func (h *PartyOperation) acceptInvite(ctx *core.ClientContext, ch *game.Character, req *request.PartyOperation) error {
	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	skipPending := ch.ShouldSkipInvitePendingForPartySearch(req.PartyID)
	member := game.PartyMemberFromCharacter(ch, worldID, int32(h.gs.config.ChannelId)).BuildRPC()
	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(
		func(c context.Context) (*internal.JoinPartyReply, error) {
			return h.gs.systems.Get[*rpc.InternalConnection]().Client.JoinParty(c, &internal.JoinPartyRequest{
				WorldId:                worldID,
				PartyId:                req.PartyID,
				Member:                 member,
				SkipInvitePendingCheck: skipPending,
			})
		},
		func(reply *internal.JoinPartyReply) error {
			if !reply.GetOk() {
				ch.Listener().OnPartyStatusMessage(ch, constant.PartyStatusForInternalError(constant.PartyC2SAcceptInvite, int32(reply.GetErrorCode())), "")
				log.Printf("PartyOperation(join): failed character=%d party=%d code=%v", charID, req.PartyID, reply.GetErrorCode())
				return nil
			}
			if reply.PartyId == nil {
				log.Printf("PartyOperation(join): ok but missing party_id character=%d", charID)
				return nil
			}
			ch.Party().Join(*reply.PartyId, reply.GetRevision())
			return nil
		},
	).OnError(func(err error) {
		log.Printf("PartyOperation(join) async error: %v", err)
	})
	return nil
}

func (h *PartyOperation) changeLeader(ctx *core.ClientContext, ch *game.Character, req *request.PartyOperation) error {
	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	partyIDPtr := ch.Party().ID()
	if partyIDPtr == nil {
		return nil
	}
	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(
		func(c context.Context) (*internal.ChangePartyLeaderReply, error) {
			return h.gs.systems.Get[*rpc.InternalConnection]().Client.ChangePartyLeader(c, &internal.ChangePartyLeaderRequest{
				WorldId:              worldID,
				PartyId:              *partyIDPtr,
				RequesterCharacterId: charID,
				NewLeaderCharacterId: req.TargetCharacterID,
			})
		},
		func(reply *internal.ChangePartyLeaderReply) error {
			if !reply.GetOk() {
				ch.Listener().OnPartyStatusMessage(ch, constant.PartyStatusForInternalError(constant.PartyC2SChangeLeader, int32(reply.GetErrorCode())), "")
				log.Printf("PartyOperation(change leader): failed character=%d party=%d target=%d code=%v", charID, *partyIDPtr, req.TargetCharacterID, reply.GetErrorCode())
			}
			return nil
		},
	).OnError(func(err error) {
		log.Printf("PartyOperation(change leader) async error: %v", err)
	})
	return nil
}

func (h *PartyOperation) invite(ctx *core.ClientContext, ch *game.Character, req *request.PartyOperation) error {
	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	targetName := req.TargetName
	promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).Bind(ch)
	if ch.Party().ID() == nil {
		promise = h.gs.systems.Get[*game.PartySystem]().CreateAsync(ctx.ActorContext, ch).Task()
	}
	promise.ThenRPC(
		func(c context.Context) (*internal.InvitePartyReply, error) {
			return h.gs.systems.Get[*rpc.InternalConnection]().Client.InviteParty(c, &internal.InvitePartyRequest{
				WorldId:             worldID,
				InviterCharacterId:  charID,
				TargetCharacterName: targetName,
			})
		},
		func(reply *internal.InvitePartyReply) error {
			if reply.GetOk() {
				return nil
			}
			code := reply.GetErrorCode()
			ch.Listener().OnPartyStatusMessage(ch, constant.PartyStatusForInternalError(constant.PartyC2SInvite, int32(code)), "")

			log.Printf("PartyOperation(invite): failed inviter=%d target=%q code=%v", charID, targetName, code)
			return nil
		},
	)
	promise.OnError(func(err error) {
		log.Printf("PartyOperation(invite) async error: %v", err)
	})
	return nil
}

func (h *PartyOperation) expel(ctx *core.ClientContext, ch *game.Character, req *request.PartyOperation) error {
	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(
		func(c context.Context) (*internal.ExpelPartyReply, error) {
			return h.gs.systems.Get[*rpc.InternalConnection]().Client.ExpelParty(c, &internal.ExpelPartyRequest{
				WorldId:              worldID,
				RequesterCharacterId: charID,
				TargetCharacterId:    req.TargetCharacterID,
			})
		},
		func(reply *internal.ExpelPartyReply) error {
			if !reply.GetOk() {
				ch.Listener().OnPartyStatusMessage(ch, constant.PartyStatusForInternalError(constant.PartyC2SExpel, int32(reply.GetErrorCode())), "")

				log.Printf("PartyOperation(expel): failed requester=%d target=%d code=%v", charID, req.TargetCharacterID, reply.GetErrorCode())
				return nil
			}
			if reply.PartyId != nil {
				h.gs.systems.Get[*game.DispatchSystem]().SendToCharacter(nil, req.TargetCharacterID, &game.ClearPartyByPartyID{
					CharacterID: req.TargetCharacterID,
					PartyID:     *reply.PartyId,
					Revision:    reply.GetRevision(),
				})
			}
			return nil
		},
	).OnError(func(err error) {
		log.Printf("PartyOperation(expel) async error: %v", err)
	})
	return nil
}
