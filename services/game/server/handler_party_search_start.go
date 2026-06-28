package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
)

type PartySearchStart struct {
	gs *GameServer
}

func (PartySearchStart) New(gs *GameServer) *PartySearchStart {
	return &PartySearchStart{
		gs: gs,
	}
}

func (h *PartySearchStart) Handle(ctx *core.ClientContext, req *request.PartySearchStart) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}

	if req.MaxLevel <= req.MinLevel ||
		req.MaxLevel-req.MinLevel > 30 ||
		req.MembersNeeded <= 0 ||
		req.MembersNeeded > 6 ||
		req.MinLevel > int32(ch.GetLevel()) ||
		req.MaxLevel < int32(ch.GetLevel()) ||
		req.ClassMask == 0 {
		return nil
	}

	cfg := &entity.PartySearchConfig{
		MinLevel:      req.MinLevel,
		MaxLevel:      req.MaxLevel,
		MembersNeeded: req.MembersNeeded,
		ClassMask:     req.ClassMask,
	}

	if pid := ch.GetPartyID(); pid != nil {
		party := h.gs.GetPartySystem().Get(*pid)
		if party == nil || party.GetLeaderCharacterId() != ch.GetID() {
			return nil
		}
	}

	if ch.GetPartyID() == nil {
		async.ThenRPC(
			async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
			func(c context.Context) (*internal.CreatePartyReply, error) {
				leader := ch.ToProtoPartyMember(h.gs.config.WorldId, int32(h.gs.config.ChannelId), internal.PartyMemberRole_PARTY_MEMBER_ROLE_LEADER)
				return h.gs.internalClient.CreateParty(c, &internal.CreatePartyRequest{
					WorldId: h.gs.config.WorldId,
					Leader:  leader,
				})
			},
			func(reply *internal.CreatePartyReply) error {
				if !reply.GetOk() || reply.PartyId == nil {
					return fmt.Errorf("party create failed")
				}
				ch.Listener.OnPartyCreated(ch, *reply.PartyId)
				ch.SetPartySearchConfig(cfg)
				return nil
			},
		)

	} else {
		ch.SetPartySearchConfig(cfg)
	}

	return nil
}
