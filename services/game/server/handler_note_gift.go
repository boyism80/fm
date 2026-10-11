package server

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
)

type NoteGift struct {
	gs *GameServer
}

func (NoteGift) New(gs *GameServer) *NoteGift {
	return &NoteGift{
		gs: gs,
	}
}

func (h *NoteGift) Handle(ctx *core.ClientContext, _ *request.NoteGift) error {
	if h.gs.internalClient == nil || ctx.ActorContext == nil {
		return nil
	}

	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return nil
	}

	rpc := &internal.NotesRequest{
		WorldId:     h.gs.config.WorldId,
		CharacterId: ch.GetID(),
		Action:      internal.NotesAction_NOTES_ACTION_REMOVE_GIFTS,
	}
	// TODO(#298): internal deletes the notes before replying, so a lost reply
	// drops their thank-you population.
	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.NotesReply, error) {
			return h.gs.internalClient.Notes(c, rpc)
		},
		func(reply *internal.NotesReply) error {
			h.gs.GetDispatchSystem().CallCharacter(ch.GetID(), func(_ actor.Context, current *entity.Character) {
				current.ReceiveNotes(reply.Notes, reply.Population)
				current.Listener.OnUnlockAction(current)
			})
			return nil
		},
	).OnError(func(error) {
		ch.Listener.OnUnlockAction(ch)
	})
	return nil
}
