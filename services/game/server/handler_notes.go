package server

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
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

	switch req.Action {
	case internal.NotesAction_NOTES_ACTION_SEND:
		if req.Kind != pconst.NoteKindNormal && req.Kind != pconst.NoteKindPopulation {
			return nil
		}
		if !entity.ClientText(req.Message).Valid(200) {
			return nil
		}
		if !entity.ClientText(req.Target).Valid(13) {
			ch.Listener.OnNoteFailed(ch, internal.NoteResult_NOTE_RESULT_WRONG_NAME)
			return nil
		}
	case internal.NotesAction_NOTES_ACTION_REMOVE, internal.NotesAction_NOTES_ACTION_LIST:
	default:
		return nil
	}

	rpc := &internal.NotesRequest{
		WorldId:     h.gs.config.WorldId,
		CharacterId: ch.GetID(),
		Action:      req.Action,
		TargetName:  req.Target,
		Message:     req.Message,
		Gift:        uint32(req.Kind),
	}
	for _, n := range req.Removals {
		rpc.NoteIds = append(rpc.NoteIds, n.ID)
	}
	// TODO(#298): internal deletes the notes before replying, so a lost reply
	// drops their thank-you population.
	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.NotesReply, error) {
			return h.gs.internalClient.Notes(c, rpc)
		},
		func(reply *internal.NotesReply) error {
			h.gs.GetDispatchSystem().CallCharacter(ch.GetID(), func(_ actor.Context, current *entity.Character) {
				switch req.Action {
				case internal.NotesAction_NOTES_ACTION_SEND:
					if reply.Status == internal.NoteResult_NOTE_RESULT_OK {
						current.Listener.OnNoteSent(current)
					} else {
						current.Listener.OnNoteFailed(current, reply.Status)
					}
				case internal.NotesAction_NOTES_ACTION_REMOVE:
					current.ReceiveNotes(reply.Notes, reply.Population)
				case internal.NotesAction_NOTES_ACTION_LIST:
					current.ReceiveNotes(reply.Notes, 0)
					current.Listener.OnNotes(current, current.Notes)
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
