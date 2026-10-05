package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/entity"
)

type SaveMapCharactersHandler struct{}

func (SaveMapCharactersHandler) New() *SaveMapCharactersHandler {
	return &SaveMapCharactersHandler{}
}

func (h *SaveMapCharactersHandler) Handle(ctx actor.Context, a *GameLogicActor, _ *SaveMapCharacters) {
	sender := ctx.Sender()
	if sender == nil {
		return
	}
	ack := &SaveMapCharactersAck{}
	maps := a.Maps()
	if len(maps) == 0 || a.GameWorld == nil {
		ctx.Respond(ack)
		return
	}
	if len(maps) == 1 && maps[0].Wz != nil {
		ack.MapID = uint32(maps[0].Wz.ID)
	}
	worldID := a.GameWorld.GetWorldID()
	entries := make([]*internal.CharacterSaveEntry, 0)
	for _, m := range maps {
		for _, obj := range m.GetAllPlayers() {
			ch, ok := obj.(*entity.Character)
			if !ok || ch == nil || ch.LoggedOut() {
				continue
			}
			entry := ch.ToProto(worldID)
			if entry == nil {
				continue
			}
			entries = append(entries, entry)
		}
	}
	ack.Saved = len(entries)
	if len(entries) == 0 {
		ctx.Respond(ack)
		return
	}
	p := a.GameWorld.SaveAsync(ctx, entries)
	if p == nil {
		ack.Err = "nil save promise"
		ctx.Respond(ack)
		return
	}
	var saveErr error
	p.OnError(func(err error) {
		saveErr = err
	}).Finally(func() {
		if saveErr != nil {
			ack.Err = saveErr.Error()
		}
		// Finally runs outside the actor receive turn, so ctx.Respond can lose the original sender/future context.
		// Send ack explicitly to the captured sender PID.
		system := ctx.ActorSystem()
		if system == nil || system.Root == nil {
			return
		}
		system.Root.Send(sender, ack)
	})
}
