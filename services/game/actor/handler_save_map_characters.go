package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/entity"
)

type SaveMapCharactersHandler struct{}

func (SaveMapCharactersHandler) New() *SaveMapCharactersHandler {
	return &SaveMapCharactersHandler{}
}

func (h *SaveMapCharactersHandler) Handle(ctx actor.Context, a *MapActor, _ *SaveMapCharacters) {
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
	chars := make([]*entity.Character, 0)
	for _, m := range maps {
		for _, obj := range m.GetAllPlayers() {
			if ch, ok := obj.(*entity.Character); ok && ch != nil {
				chars = append(chars, ch)
			}
		}
	}
	ack.Saved = len(chars)
	if len(chars) == 0 {
		ctx.Respond(ack)
		return
	}
	p := a.GameWorld.SaveAsync(ctx, chars)
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
