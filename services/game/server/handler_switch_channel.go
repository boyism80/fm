package server

import (
	"context"
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/types"
)

type SwitchChannel struct {
	gs *GameServer
}

func (SwitchChannel) New(gs *GameServer) *SwitchChannel {
	return &SwitchChannel{
		gs: gs,
	}
}

func (h *SwitchChannel) Handle(ctx *core.ClientContext, req *request.SwitchChannel) error {
	if h.gs == nil || req == nil {
		return fmt.Errorf("switch channel: invalid state")
	}
	if h.gs.internalClient == nil {
		return fmt.Errorf("switch channel: internal client not configured")
	}
	if ctx.ActorContext == nil {
		return fmt.Errorf("switch channel: actor context required for internal RPC")
	}

	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("switch channel: invalid client type")
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("switch channel: character not loaded")
	}

	targetChannel := uint32(req.Channel)
	worldID := h.gs.config.WorldId
	if targetChannel == h.gs.config.ChannelId {
		_ = ctx.Client.Send(&response.ServerBlocked{
			Reason: constant.ServerBlockedChannelMoveUnavailable,
		}, types.SEND_POLICY_ENCRYPT)
		return fmt.Errorf("switch channel: already on channel %d", targetChannel)
	}

	h.gs.handOver(ctx, gameClient, character, constant.ServerBlockedChannelMoveUnavailable, func(c context.Context) (*response.SwitchChannel, error) {
		statusReply, err := h.gs.internalClient.GetGameChannelStatus(c, &internal.GetGameChannelStatusRequest{
			WorldId:   worldID,
			ChannelId: targetChannel,
		})
		if err != nil {
			return nil, err
		}
		if !statusReply.GetFound() || !statusReply.GetAlive() || statusReply.GetChannelFull() {
			return nil, nil
		}
		return &response.SwitchChannel{IP: statusReply.GetHost(), Port: uint16(statusReply.GetPort())}, nil
	})
	return nil
}
