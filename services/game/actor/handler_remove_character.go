package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type RemoveCharacterHandler struct{}

func (RemoveCharacterHandler) New() *RemoveCharacterHandler {
	return &RemoveCharacterHandler{}
}

func (h *RemoveCharacterHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *RemoveCharacter) {
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		if ctx.Sender() != nil {
			ctx.Respond(struct{}{})
		}
		return
	}
	_ = m.RemovePlayer(msg.CharacterID)
	if ctx.Sender() != nil {
		ctx.Respond(struct{}{})
	}
}
