package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

func (gs *GameServer) handOver(ctx *core.ClientContext, gameClient *client.GameClient, character *entity.Character, blocked constant.ServerBlockedReason, findRoute func(context.Context) (*response.SwitchChannel, error)) {
	ic := gs.internalClient
	worldID := gs.config.WorldId
	var route *response.SwitchChannel
	var entry *internal.CharacterSaveEntry
	var debuffs []*internal.Debuff
	block := func() {
		_ = ctx.Client.Send(&response.ServerBlocked{Reason: blocked}, types.SEND_POLICY_ENCRYPT)
	}

	// Packets that arrive while the character is saved and handed over are dropped, so nothing changes after the save.
	gameClient.SetChangingChannel(true)
	promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise.ThenRPC(findRoute, func(found *response.SwitchChannel) error {
		if found == nil || found.IP == "" || found.Port == 0 {
			block()
			return fmt.Errorf("hand over: no route (world=%d)", worldID)
		}
		route = found
		entry = character.ToProto(worldID)
		if entry == nil {
			block()
			return fmt.Errorf("hand over: character %d has no map to save", character.GetID())
		}
		debuffs = character.Debuffs.ToProto()
		return nil
	})
	promise.ThenRPC(func(c context.Context) (*internal.SaveCharactersReply, error) {
		return ic.SaveCharacters(c, &internal.SaveCharactersRequest{Entries: []*internal.CharacterSaveEntry{entry}})
	}, func(saveReply *internal.SaveCharactersReply) error {
		if saveReply == nil || !saveReply.GetOk() {
			block()
			return fmt.Errorf("hand over: save failed (world=%d)", worldID)
		}
		return nil
	})
	accID := character.AccountID
	charID := character.GetID()
	sourceChannel := gs.config.ChannelId
	promise.ThenRPC(func(c context.Context) (*internal.BeginGameTransitionReply, error) {
		if accID == 0 || charID == 0 {
			return nil, fmt.Errorf("hand over: missing account or character id")
		}
		return ic.BeginGameTransition(c, &internal.BeginGameTransitionRequest{
			WorldId:         worldID,
			AccountId:       accID,
			CharacterId:     charID,
			ClientIp:        ctx.Client.GetRemoteIP(),
			Debuffs:         debuffs,
			SourceChannelId: &sourceChannel,
		})
	}, func(transReply *internal.BeginGameTransitionReply) error {
		if transReply == nil || !transReply.GetOk() {
			block()
			return fmt.Errorf("hand over: begin transition failed (world=%d code=%v)", worldID, transReply.GetErrorCode())
		}

		// The character leaves this channel now, not when the client disconnects: a client that stays connected must not keep playing it,
		// and the disconnect must neither save it over the next server nor end the moving session.
		left := gameClient.Logout()
		if left == nil {
			_ = ctx.Client.Send(route, types.SEND_POLICY_ENCRYPT)
			return nil
		}
		logout := left.MarkLoggedOut()
		go func() {
			gs.removeCharacter(left, logout)
			_ = ctx.Client.Send(route, types.SEND_POLICY_ENCRYPT)
		}()
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("hand over (async): %v", err)
		gameClient.SetChangingChannel(false)
	})
}
