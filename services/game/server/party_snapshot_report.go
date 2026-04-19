package server

import (
	"context"
	"log"

	"github.com/boyism80/fm/core"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/entity"
)

func (gs *GameServer) UpdatePartyMemberAsync(ch *entity.Character) {
	if gs == nil || ch == nil || gs.internalClient == nil {
		return
	}
	if ch.GetPartyID() == nil {
		return
	}
	mm := ch.ToGrpcPartyMember(uint32(gs.config.WorldId), int32(gs.config.ChannelId), "MEMBER")
	if mm == nil {
		return
	}
	req := &internal.UpdatePartyMemberRequest{Member: mm}
	cid := ch.GetID()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), core.InternalRPCPerStepTimeout)
		defer cancel()
		if _, err := gs.internalClient.UpdatePartyMember(ctx, req); err != nil {
			log.Printf("UpdatePartyMember async char %d: %v", cid, err)
		}
	}()
}
