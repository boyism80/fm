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
	if a.Map == nil || a.GameWorld == nil {
		ctx.Respond(ack)
		return
	}
	if a.Map.Wz != nil {
		ack.MapID = uint32(a.Map.Wz.ID)
	}
	allPlayers := a.Map.GetAllPlayers()
	if len(allPlayers) == 0 {
		ctx.Respond(ack)
		return
	}
	chars := make([]*entity.Character, 0, len(allPlayers))
	for _, obj := range allPlayers {
		if ch, ok := obj.(*entity.Character); ok && ch != nil {
			chars = append(chars, ch)
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
