package server

import (
	"github.com/boyism80/fm/core"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
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

	ch.Notes.Remove(ctx.ActorContext, internal.NotesAction_NOTES_ACTION_REMOVE_GIFTS, nil).Finally(func() {
		ch.Listener.OnUnlockAction(ch)
	})
	return nil
}
