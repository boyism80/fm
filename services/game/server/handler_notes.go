package server

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/rpc"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
)

type Notes struct {
	gs *GameServer
}

func (Notes) New(gs *GameServer) *Notes {
	return &Notes{
		gs: gs,
	}
}

func (h *Notes) Handle(ctx *core.ClientContext, req *request.Notes) error {
	if ctx.ActorContext == nil {
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

	switch req.Action {
	case internal.NotesAction_NOTES_ACTION_SEND:
		if req.Kind != pconst.NoteKindNormal && req.Kind != pconst.NoteKindPopulation {
			return nil
		}
		if !game.ClientText(req.Message).Valid(200) {
			return nil
		}
		if !game.ClientText(req.Target).Valid(13) {
			ch.Listener.OnNoteFailed(ch, internal.NoteResult_NOTE_RESULT_WRONG_NAME)
			return nil
		}
	case internal.NotesAction_NOTES_ACTION_REMOVE:
		noteIDs := make([]uint32, 0, len(req.Removals))
		for _, n := range req.Removals {
			noteIDs = append(noteIDs, n.ID)
		}
		ch.Notes.Remove(ctx.ActorContext, req.Action, noteIDs)
		return nil
	case internal.NotesAction_NOTES_ACTION_LIST:
	default:
		return nil
	}

	internalReq := &internal.NotesRequest{
		WorldId:     h.gs.config.WorldId,
		CharacterId: ch.GetID(),
		Action:      req.Action,
		TargetName:  req.Target,
		Message:     req.Message,
		Gift:        uint32(req.Kind),
	}
	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.NotesReply, error) {
			return h.gs.systems.Get[*rpc.InternalConnection]().Client.Notes(c, internalReq)
		},
		func(reply *internal.NotesReply) error {
			h.gs.systems.Get[*game.DispatchSystem]().Visit(ch.GetID(), func(_ actor.Context, current *game.Character) {
				switch req.Action {
				case internal.NotesAction_NOTES_ACTION_SEND:
					if reply.Status == internal.NoteResult_NOTE_RESULT_OK {
						current.Listener.OnNoteSent(current)
					} else {
						current.Listener.OnNoteFailed(current, reply.Status)
					}
				case internal.NotesAction_NOTES_ACTION_LIST:
					current.Notes.Receive(reply.Notes, 0)
					current.Listener.OnNotes(current, current.Notes.Entries)
				}
			})
			return nil
		},
	).OnError(func(error) {
		if req.Action == internal.NotesAction_NOTES_ACTION_SEND {
			ch.Listener.OnNoteFailed(ch, internal.NoteResult_NOTE_RESULT_INBOX_FULL)
		}
	})
	return nil
}
