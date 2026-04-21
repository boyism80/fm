package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type PartyOperation struct {
	gs     *GameServer
	opcode byte
}

func (PartyOperation) New(gs *GameServer) *PartyOperation {
	return &PartyOperation{
		gs:     gs,
		opcode: 0x66,
	}
}

func (h *PartyOperation) GetOpcode() byte {
	return h.opcode
}

func (h *PartyOperation) Handle(ctx *core.ClientContext, req *request.PartyOperation) error {
	if h.gs.internalClient == nil {
		return fmt.Errorf("internal client not configured")
	}
	if ctx.ActorContext == nil {
		return fmt.Errorf("party operation: actor context required")
	}

	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return fmt.Errorf("character not found")
	}

	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	switch req.Operation {
	case constant.PartyC2SCreate:
		async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
			func(c context.Context) (*internal.CreatePartyReply, error) {
				leader := ch.ToProtoPartyMember(worldID, int32(h.gs.config.ChannelId), "LEADER")
				req := &internal.CreatePartyRequest{
					WorldId: worldID,
					Leader:  leader,
				}
				return h.gs.internalClient.CreateParty(c, req)
			},
			func(reply *internal.CreatePartyReply) error {
				if !reply.GetOk() {
					ch.Listener.OnPartyStatusMessage(ch, constant.PartyStatusForInternalError(constant.PartyC2SCreate, int32(reply.GetErrorCode())))
					log.Printf("PartyOperation(create): failed character=%d code=%v", charID, reply.GetErrorCode())
					return nil
				}
				if reply.PartyId == nil {
					log.Printf("PartyOperation(create): ok but missing party_id character=%d", charID)
					return nil
				}
				ch.Listener.OnPartyCreated(ch, *reply.PartyId)
				return nil
			},
		).OnError(func(err error) {
			log.Printf("PartyOperation(create) async error: %v", err)
		}).Run()
		return nil

	case constant.PartyC2SLeave:
		async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
			func(c context.Context) (*internal.LeavePartyReply, error) {
				return h.gs.internalClient.LeaveParty(c, &internal.LeavePartyRequest{
					WorldId:     worldID,
					CharacterId: charID,
				})
			},
			func(reply *internal.LeavePartyReply) error {
				if !reply.GetOk() {
					ch.Listener.OnPartyStatusMessage(ch, constant.PartyStatusForInternalError(constant.PartyC2SLeave, int32(reply.GetErrorCode())))
					log.Printf("PartyOperation(leave): failed character=%d code=%v", charID, reply.GetErrorCode())
				}
				return nil
			},
		).OnError(func(err error) {
			log.Printf("PartyOperation(leave) async error: %v", err)
		}).Run()
		return nil

	case constant.PartyC2SAcceptInvite:
		async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
			func(c context.Context) (*internal.JoinPartyReply, error) {
				member := ch.ToProtoPartyMember(worldID, int32(h.gs.config.ChannelId), "MEMBER")
				return h.gs.internalClient.JoinParty(c, &internal.JoinPartyRequest{
					WorldId: worldID,
					PartyId: req.PartyID,
					Member:  member,
				})
			},
			func(reply *internal.JoinPartyReply) error {
				if !reply.GetOk() {
					ch.Listener.OnPartyStatusMessage(ch, constant.PartyStatusForInternalError(constant.PartyC2SAcceptInvite, int32(reply.GetErrorCode())))
					log.Printf("PartyOperation(join): failed character=%d party=%d code=%v", charID, req.PartyID, reply.GetErrorCode())
					return nil
				}
				if reply.PartyId == nil {
					log.Printf("PartyOperation(join): ok but missing party_id character=%d", charID)
					return nil
				}
				SyncPartyMemberHPOnMapEnter(ch.GetMap(), ch, *reply.PartyId)
				return nil
			},
		).OnError(func(err error) {
			log.Printf("PartyOperation(join) async error: %v", err)
		}).Run()
		return nil

	case constant.PartyC2SChangeLeader:
		partyIDPtr := ch.GetPartyID()
		if partyIDPtr == nil || req.TargetCharacterID == 0 {
			return nil
		}
		async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
			func(c context.Context) (*internal.ChangePartyLeaderReply, error) {
				return h.gs.internalClient.ChangePartyLeader(c, &internal.ChangePartyLeaderRequest{
					WorldId:              worldID,
					PartyId:              *partyIDPtr,
					RequesterCharacterId: charID,
					NewLeaderCharacterId: req.TargetCharacterID,
				})
			},
			func(reply *internal.ChangePartyLeaderReply) error {
				if !reply.GetOk() {
					ch.Listener.OnPartyStatusMessage(ch, constant.PartyStatusForInternalError(constant.PartyC2SChangeLeader, int32(reply.GetErrorCode())))
					log.Printf("PartyOperation(change leader): failed character=%d party=%d target=%d code=%v", charID, *partyIDPtr, req.TargetCharacterID, reply.GetErrorCode())
				}
				return nil
			},
		).OnError(func(err error) {
			log.Printf("PartyOperation(change leader) async error: %v", err)
		}).Run()
		return nil

	case constant.PartyC2SInvite:
		targetName := req.TargetName
		pid := ch.GetPartyID()
		hasParty := pid != nil
		var inv *internal.InvitePartyReply
		promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		if !hasParty {
			promise = async.ThenRPC(promise,
				func(c context.Context) (*internal.CreatePartyReply, error) {
					leader := ch.ToProtoPartyMember(worldID, int32(h.gs.config.ChannelId), "LEADER")
					req := &internal.CreatePartyRequest{
						WorldId: worldID,
						Leader:  leader,
					}
					return h.gs.internalClient.CreateParty(c, req)
				},
				func(reply *internal.CreatePartyReply) error {
					if reply.GetOk() {
						if reply.PartyId == nil {
							log.Printf("PartyOperation(invite pre-create): ok but missing party_id character=%d", charID)
							return fmt.Errorf("party create before invite failed")
						}
						ch.Listener.OnPartyCreated(ch, *reply.PartyId)
						return nil
					}
					ch.Listener.OnPartyStatusMessage(ch, constant.PartyStatusForInternalError(constant.PartyC2SCreate, int32(reply.GetErrorCode())))
					log.Printf("PartyOperation(invite pre-create): failed character=%d target=%q code=%v", charID, targetName, reply.GetErrorCode())
					return fmt.Errorf("party create before invite failed")
				},
			)
		}
		promise = async.ThenRPC(promise,
			func(c context.Context) (*internal.InvitePartyReply, error) {
				return h.gs.internalClient.InviteParty(c, &internal.InvitePartyRequest{
					WorldId:             worldID,
					InviterCharacterId:  charID,
					TargetCharacterName: targetName,
				})
			},
			func(reply *internal.InvitePartyReply) error {
				inv = reply
				if reply.GetOk() {
					return nil
				}
				code := reply.GetErrorCode()
				ch.Listener.OnPartyStatusMessage(ch, constant.PartyStatusForInternalError(constant.PartyC2SInvite, int32(code)))
				log.Printf("PartyOperation(invite): failed inviter=%d target=%q code=%v", charID, targetName, code)
				return nil
			},
		)
		promise = async.ThenRPC(promise,
			func(c context.Context) (*internal.GetPartyReply, error) {
				if inv == nil || !inv.GetOk() {
					return &internal.GetPartyReply{Found: false}, nil
				}
				if inv.PartyId == nil {
					log.Printf("PartyOperation(invite): ok but missing party_id inviter=%d", charID)
					return &internal.GetPartyReply{Found: false}, nil
				}
				return h.gs.internalClient.GetParty(c, &internal.GetPartyRequest{
					WorldId: worldID,
					PartyId: *inv.PartyId,
				})
			},
			func(gp *internal.GetPartyReply) error {
				if gp.GetFound() && gp.GetParty() != nil && h.gs.party != nil {
					h.gs.party.Update(gp.GetParty())
				}
				return nil
			},
		)
		promise.OnError(func(err error) {
			log.Printf("PartyOperation(invite) async error: %v", err)
		}).Run()
		return nil

	case constant.PartyC2SExpel:
		async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
			func(c context.Context) (*internal.ExpelPartyReply, error) {
				return h.gs.internalClient.ExpelParty(c, &internal.ExpelPartyRequest{
					WorldId:              worldID,
					RequesterCharacterId: charID,
					TargetCharacterId:    req.TargetCharacterID,
				})
			},
			func(reply *internal.ExpelPartyReply) error {
				if !reply.GetOk() {
					ch.Listener.OnPartyStatusMessage(ch, constant.PartyStatusForInternalError(constant.PartyC2SExpel, int32(reply.GetErrorCode())))
					log.Printf("PartyOperation(expel): failed requester=%d target=%d code=%v", charID, req.TargetCharacterID, reply.GetErrorCode())
				}
				return nil
			},
		).OnError(func(err error) {
			log.Printf("PartyOperation(expel) async error: %v", err)
		}).Run()
		return nil

	default:
		return nil
	}
}
