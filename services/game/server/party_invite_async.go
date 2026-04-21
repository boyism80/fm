package server

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (gs *GameServer) RequestAutoInvitePartyAsync(ctx actor.Context, inviterCharacterID uint32, targetCharacterIDs []uint32) *async.Promise {
	p := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	if gs == nil || gs.internalClient == nil || inviterCharacterID == 0 || len(targetCharacterIDs) == 0 {
		return p
	}
	async.ThenRPC(p,
		func(c context.Context) (*internal.AutoInvitePartyReply, error) {
			return gs.internalClient.AutoInviteParty(c, &internal.AutoInvitePartyRequest{
				WorldId:            gs.config.WorldId,
				InviterCharacterId: inviterCharacterID,
				TargetCharacterIds: targetCharacterIDs,
			})
		},
		func(reply *internal.AutoInvitePartyReply) error {
			_ = reply
			return nil
		},
	)
	return p
}
