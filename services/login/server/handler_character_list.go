package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/login/client"
	"github.com/boyism80/fm/types"
)

type CharacterList struct {
	ls *LoginServer
}

func (CharacterList) New(ls *LoginServer) *CharacterList {
	return &CharacterList{
		ls: ls,
	}
}

func (h *CharacterList) Handle(ctx *core.ClientContext, req *request.CharacterList) error {
	log.Printf("Character list packet received from %s - Server: %d, Channel: %d",
		ctx.Client.GetConnection().RemoteAddr(), req.Server, req.Channel)

	loginClient, ok := ctx.Client.(*client.LoginClient)
	if !ok {
		return nil
	}

	worldId := uint32(req.Server)
	loginClient.SetWorldId(worldId)
	loginClient.SetChannelId(req.Channel)
	accountId := loginClient.GetAccountId()

	ic := h.ls.internalClient
	if ic == nil {
		charListResp := &response.CharacterList{Characters: nil, SlotCount: 6}
		return ctx.Client.Send(charListResp, types.SEND_POLICY_ENCRYPT)
	}

	reqMsg := &internal.GetCharacterListRequest{
		AccountId: accountId,
		WorldId:   worldId,
	}

	if ctx.ActorContext == nil {
		log.Printf("CharacterList: no actor context, cannot run internal RPC")
		charListResp := &response.CharacterList{Characters: nil, SlotCount: 6}
		_ = ctx.Client.Send(charListResp, types.SEND_POLICY_ENCRYPT)
		return fmt.Errorf("character list: actor context required for internal RPC")
	}

	sendCharListFallback := true
	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.GetGameChannelStatusReply, error) {
		return ic.GetGameChannelStatus(c, &internal.GetGameChannelStatusRequest{
			WorldId:   worldId,
			ChannelId: uint32(req.Channel),
		})
	}, func(statusReply *internal.GetGameChannelStatusReply) error {
		if statusReply == nil || !statusReply.GetAlive() {
			sendCharListFallback = false
			_ = ctx.Client.Send(&response.LoginFailed{Reason: response.LoginFailedReasonTooManyConnections}, types.SEND_POLICY_ENCRYPT)
			return fmt.Errorf("game channel not alive (world=%d channel=%d)", worldId, req.Channel)
		}
		if statusReply.GetChannelFull() {
			sendCharListFallback = false
			_ = ctx.Client.Send(&response.LoginFailed{Reason: response.LoginFailedReasonTooManyConnections}, types.SEND_POLICY_ENCRYPT)
			return fmt.Errorf("game channel at capacity (world=%d channel=%d)", worldId, req.Channel)
		}
		return nil
	})
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.GetCharacterListReply, error) {
		return ic.GetCharacterList(c, reqMsg)
	}, func(reply *internal.GetCharacterListReply) error {
		return h.sendCharacterList(ctx, reply)
	})
	promise.OnError(func(err error) {
		log.Printf("CharacterList (async): %v", err)
		if !sendCharListFallback {
			return
		}
		charListResp := &response.CharacterList{Characters: nil, SlotCount: 6}
		_ = ctx.Client.Send(charListResp, types.SEND_POLICY_ENCRYPT)
	})
	return nil
}

func (h *CharacterList) sendCharacterList(ctx *core.ClientContext, reply *internal.GetCharacterListReply) error {
	characters := make([]dto.Character, 0, len(reply.Characters))
	for _, ov := range reply.Characters {
		characters = append(characters, overviewToDto(ov))
	}

	charListResp := &response.CharacterList{
		Characters: characters,
		SlotCount:  reply.SlotCount,
	}
	return ctx.Client.Send(charListResp, types.SEND_POLICY_ENCRYPT)
}

func overviewToDto(ov *internal.CharacterOverview) dto.Character {
	baseLooks := make(map[int8]uint32)
	for k, v := range ov.BaseLooks {
		baseLooks[int8(k)] = v
	}
	overlays := make(map[int8]uint32)
	for k, v := range ov.Overlays {
		overlays[int8(k)] = v
	}
	return dto.Character{
		ID:            ov.CharacterId,
		Name:          ov.Name,
		Gender:        uint8(ov.Gender),
		SkinColor:     uint8(ov.SkinColor),
		Face:          ov.Face,
		Hair:          ov.Hair,
		Level:         uint8(ov.Level),
		Class:         uint16(ov.ClassId),
		Map:           ov.MapId,
		SpawnPoint:    uint8(ov.SpawnPoint),
		Rank:          uint32(ov.Rank),
		RankDiff:      ov.RankDiff,
		ClassRank:     uint32(ov.ClassRank),
		ClassRankDiff: ov.ClassRankDiff,
		BaseLooks:     baseLooks,
		Overlays:      overlays,
	}
}
