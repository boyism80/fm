package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type RemoveCharacterHandler struct{}

func (RemoveCharacterHandler) New() *RemoveCharacterHandler {
	return &RemoveCharacterHandler{}
}

func (h *RemoveCharacterHandler) Handle(ctx actor.Context, a *MapActor, msg *RemoveCharacter) {
	if a.Map == nil {
		if ctx.Sender() != nil {
			ctx.Respond(struct{}{})
		}
		return
	}
	_ = a.Map.RemovePlayer(msg.CharacterID)
	if ctx.Sender() != nil {
		ctx.Respond(struct{}{})
	}
}
