package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/rpc"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (gs *GameServer) handOver(ctx *core.ClientContext, gameClient *client.GameClient, character *game.Character, blocked constant.ServerBlockedReason, findRoute func(context.Context) (*response.SwitchChannel, error)) {
	ic := gs.systems.Get[*rpc.InternalConnection]().Client
	worldID := gs.config.WorldId
	var route *response.SwitchChannel
	var entry *internal.CharacterSaveEntry
	var debuffs []*internal.Debuff
	block := func() {
		_ = ctx.Client.Send(&response.ServerBlocked{Reason: blocked}, types.SEND_POLICY_ENCRYPT)
	}
	// The character leaves this channel now, not when the client disconnects: a client that stays connected must not keep playing it,
	// and the disconnect must neither save it over the next server nor end the moving session.
	leave := func(then func()) {
		left := gameClient.Logout()
		if left == nil {
			then()
			return
		}
		logout := left.MarkLoggedOut()
		go func() {
			gs.removeCharacter(left, logout)
			then()
		}()
	}
	transitionPending := false

	// The save below snapshots escrow as owned, so a room still open after it could hand the same goods over again.
	if character.MiniRoom() != nil {
		block()
		return
	}
	// A script or shop still open after the save could change the character the next server already loaded.
	if character.Dialog().Busy() {
		block()
		return
	}

	character.LieDetector().Escape()

	// Packets that arrive while the character is saved and handed over are dropped, so nothing changes after the save.
	gameClient.SetChangingChannel(true)
	promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise.ThenRPC(findRoute, func(found *response.SwitchChannel) error {
		if found == nil || found.IP == "" || found.Port == 0 {
			block()
			return fmt.Errorf("hand over: no route (world=%d)", worldID)
		}
		route = found
		entry = character.BuildRPC(worldID)
		if entry == nil {
			block()
			return fmt.Errorf("hand over: character %d has no map to save", character.GetID())
		}
		debuffs = character.Debuffs().BuildRPC()
		return nil
	})
	promise.ThenRPC(func(c context.Context) (*internal.SaveCharactersReply, error) {
		return ic.SaveCharacters(c, &internal.SaveCharactersRequest{Entries: []*internal.CharacterSaveEntry{entry}}, grpc.WaitForReady(true))
	}, func(saveReply *internal.SaveCharactersReply) error {
		if saveReply == nil || !saveReply.GetOk() {
			block()
			return fmt.Errorf("hand over: save failed (world=%d)", worldID)
		}
		transitionPending = true
		return nil
	})
	accID := character.AccountID()
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
		}, grpc.WaitForReady(true))
	}, func(transReply *internal.BeginGameTransitionReply) error {
		transitionPending = false
		if transReply == nil || !transReply.GetOk() {
			block()
			return fmt.Errorf("hand over: begin transition failed (world=%d code=%v)", worldID, transReply.GetErrorCode())
		}

		leave(func() {
			_ = ctx.Client.Send(route, types.SEND_POLICY_ENCRYPT)
		})
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("hand over (async): %v", err)

		// A transition internal never answered may have committed; playing on would be dropped unsaved when the session refresh fails.
		// A promise future timeout is not a gRPC status.
		if transitionPending {
			s, ok := status.FromError(err)
			if !ok || s.Code() == codes.DeadlineExceeded || s.Code() == codes.Unavailable || s.Code() == codes.Canceled {
				leave(func() {
					_ = ctx.Client.GetConnection().Close()
				})
				return
			}
		}
		gameClient.SetChangingChannel(false)
	})
}
