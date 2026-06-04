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

func (s partySystem) UpdateMemberAsync(ctx actor.Context, ch *entity.Character) *async.Promise {
	p := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
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
	cid := ch.GetID()
	p.OnError(func(err error) {
		log.Printf("UpdatePartyMember async char %d: %v", cid, err)
	})
	async.ThenRPC(p, func(c context.Context) (*internal.UpdatePartyMemberReply, error) {
		return s.gs.internalClient.UpdatePartyMember(c, req)
	}, func(*internal.UpdatePartyMemberReply) error {
		return nil
	})
	return p
}
