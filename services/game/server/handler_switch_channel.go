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
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
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
	log.Printf("SwitchChannel recv: channel_id=%d (0-based) from %s",
		req.Channel, ctx.Client.GetConnection().RemoteAddr())
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

	ic := h.gs.internalClient
	var routeHost string
	var routePort uint16

	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.GetGameChannelStatusReply, error) {
		return ic.GetGameChannelStatus(c, &internal.GetGameChannelStatusRequest{
			WorldId:   worldID,
			ChannelId: targetChannel,
		})
	}, func(statusReply *internal.GetGameChannelStatusReply) error {
		if statusReply == nil || !statusReply.GetFound() || !statusReply.GetAlive() || statusReply.GetChannelFull() {
			_ = ctx.Client.Send(&response.ServerBlocked{
				Reason: constant.ServerBlockedChannelMoveUnavailable,
			}, types.SEND_POLICY_ENCRYPT)
			if statusReply == nil {
				return fmt.Errorf("switch channel: nil status reply (world=%d channel=%d)", worldID, targetChannel)
			}
			return fmt.Errorf("switch channel: unavailable (world=%d channel=%d found=%v alive=%v full=%v)",
				worldID, targetChannel, statusReply.GetFound(), statusReply.GetAlive(), statusReply.GetChannelFull())
		}
		routeHost = statusReply.GetHost()
		routePort = uint16(statusReply.GetPort())
		if routeHost == "" || routePort == 0 {
			_ = ctx.Client.Send(&response.ServerBlocked{
				Reason: constant.ServerBlockedChannelMoveUnavailable,
			}, types.SEND_POLICY_ENCRYPT)
			return fmt.Errorf("switch channel: missing route (world=%d channel=%d)", worldID, targetChannel)
		}
		return nil
	})
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.SaveCharactersReply, error) {
		return h.gs.grpcSaveCharacters(c, []*entity.Character{character})
	}, func(saveReply *internal.SaveCharactersReply) error {
		if saveReply == nil || !saveReply.GetOk() {
			_ = ctx.Client.Send(&response.ServerBlocked{
				Reason: constant.ServerBlockedChannelMoveUnavailable,
			}, types.SEND_POLICY_ENCRYPT)
			return fmt.Errorf("switch channel: save failed (world=%d channel=%d)", worldID, targetChannel)
		}
		return nil
	})
	accID := character.AccountID
	charID := character.GetID()
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.BeginGameTransitionReply, error) {
		if accID == 0 || charID == 0 {
			return nil, fmt.Errorf("switch channel: missing account or character id")
		}
		return ic.BeginGameTransition(c, &internal.BeginGameTransitionRequest{
			WorldId:     worldID,
			AccountId:   accID,
			CharacterId: charID,
		})
	}, func(transReply *internal.BeginGameTransitionReply) error {
		if transReply == nil || !transReply.GetOk() {
			_ = ctx.Client.Send(&response.ServerBlocked{
				Reason: constant.ServerBlockedChannelMoveUnavailable,
			}, types.SEND_POLICY_ENCRYPT)
			code := internal.SessionErrorCode_SESSION_NONE
			if transReply != nil {
				code = transReply.GetErrorCode()
			}
			return fmt.Errorf("switch channel: begin transition failed (world=%d channel=%d code=%v)",
				worldID, targetChannel, code)
		}
		if err := ctx.Client.Send(&response.SwitchChannel{
			IP:   routeHost,
			Port: routePort,
		}, types.SEND_POLICY_ENCRYPT); err != nil {
			return err
		}
		gameClient.SetTransferDisconnect(true)
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("SwitchChannel (async): %v", err)
	})
	return nil
}
