package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
)

type Buddy struct {
	gs *GameServer
}

func (Buddy) New(gs *GameServer) *Buddy {
	return &Buddy{
		gs: gs,
	}
}

func (h *Buddy) Handle(ctx *core.ClientContext, req *request.Buddy) error {
	if h.gs.internalClient == nil {
		return fmt.Errorf("internal client not configured")
	}
	if ctx.ActorContext == nil {
		return fmt.Errorf("buddy: actor context required")
	}

	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return fmt.Errorf("character not found")
	}

	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	switch req.Mode {
	case pconst.BuddyAdd:
		targetName := req.Name
		groupName := req.Group
		async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
			func(c context.Context) (*internal.RequestBuddyReply, error) {
				return h.gs.internalClient.RequestBuddy(c, &internal.RequestBuddyRequest{
					WorldId:              worldID,
					RequesterCharacterId: charID,
					TargetCharacterName:  targetName,
					GroupName:            groupName,
				})
			},
			func(reply *internal.RequestBuddyReply) error {
				if !reply.GetOk() {
					ch.Listener.OnBuddyStatusMessage(ch, pconst.BuddyStatusForInternalError(int32(reply.GetErrorCode())))
					log.Printf("Buddy(add): failed requester=%d target=%q code=%v", charID, targetName, reply.GetErrorCode())
					return nil
				}
				if view := reply.GetRequesterView(); view != nil {
					applyBuddyViewToCharacter(ch, view)
					ch.Listener.OnBuddyListUpdate(ch, pconst.BuddyListSyncUpdate, []response.BuddyEntry{buddyEntryResponseFromProto(view)})
				}
				return nil
			},
		).OnError(func(err error) {
			log.Printf("Buddy(add) async error: %v", err)
		})
		return nil

	case pconst.BuddyAccept:
		requesterID := req.CharacterID
		async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
			func(c context.Context) (*internal.AcceptBuddyReply, error) {
				return h.gs.internalClient.AcceptBuddy(c, &internal.AcceptBuddyRequest{
					WorldId:              worldID,
					AccepterCharacterId:  charID,
					RequesterCharacterId: requesterID,
				})
			},
			func(reply *internal.AcceptBuddyReply) error {
				if !reply.GetOk() {
					ch.Listener.OnBuddyStatusMessage(ch, pconst.BuddyStatusForInternalError(int32(reply.GetErrorCode())))
					log.Printf("Buddy(accept): failed accepter=%d requester=%d code=%v", charID, requesterID, reply.GetErrorCode())
					return nil
				}
				if view := reply.GetAccepterView(); view != nil {
					applyBuddyViewToCharacter(ch, view)
					ch.Listener.OnBuddyListUpdate(ch, pconst.BuddyListSyncUpdate, []response.BuddyEntry{buddyEntryResponseFromProto(view)})
				}
				return nil
			},
		).OnError(func(err error) {
			log.Printf("Buddy(accept) async error: %v", err)
		})
		return nil

	case pconst.BuddyDelete:
		buddyID := req.CharacterID
		async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
			func(c context.Context) (*internal.RemoveBuddyReply, error) {
				return h.gs.internalClient.RemoveBuddy(c, &internal.RemoveBuddyRequest{
					WorldId:          worldID,
					CharacterId:      charID,
					BuddyCharacterId: buddyID,
				})
			},
			func(reply *internal.RemoveBuddyReply) error {
				if !reply.GetOk() {
					ch.Listener.OnBuddyStatusMessage(ch, pconst.BuddyStatusForInternalError(int32(reply.GetErrorCode())))
					log.Printf("Buddy(delete): failed character=%d buddy=%d code=%v", charID, buddyID, reply.GetErrorCode())
					return nil
				}
				ch.BuddyList().Remove(buddyID)
				ch.Listener.OnBuddyListUpdate(ch, pconst.BuddyListSyncDelete, ch.BuddyList().SnapshotForClient())
				return nil
			},
		).OnError(func(err error) {
			log.Printf("Buddy(delete) async error: %v", err)
		})
		return nil

	default:
		return nil
	}
}

func applyBuddyViewToCharacter(ch *entity.Character, view *internal.BuddyEntry) {
	if ch == nil || view == nil {
		return
	}
	ch.BuddyList().Upsert(entity.BuddyListEntryFromProto(view))
}

func buddyEntryResponseFromProto(entry *internal.BuddyEntry) response.BuddyEntry {
	return entity.BuddyListEntryFromProto(entry).ToResponse()
}
