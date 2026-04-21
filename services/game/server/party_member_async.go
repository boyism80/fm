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

// UpdatePartyMemberAsync builds a Promise for UpdatePartyMember. Caller must Run(). Pass ctx nil when unavailable.
func (gs *GameServer) UpdatePartyMemberAsync(ctx actor.Context, ch *entity.Character) *async.Promise {
	p := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	if gs == nil || ch == nil || gs.internalClient == nil {
		return p
	}
	if ch.GetPartyID() == nil {
		return p
	}
	mm := ch.ToProtoPartyMember(uint32(gs.config.WorldId), int32(gs.config.ChannelId), "MEMBER")
	if mm == nil {
		return p
	}
	req := &internal.UpdatePartyMemberRequest{Member: mm}
	cid := ch.GetID()
	p.OnError(func(err error) {
		log.Printf("UpdatePartyMember async char %d: %v", cid, err)
	})
	async.ThenRPC(p, func(c context.Context) (*internal.UpdatePartyMemberReply, error) {
		return gs.internalClient.UpdatePartyMember(c, req)
	}, func(*internal.UpdatePartyMemberReply) error {
		return nil
	})
	return p
}
