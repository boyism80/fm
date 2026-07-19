package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type AddCharacterHandler struct{}

func (AddCharacterHandler) New() *AddCharacterHandler {
	return &AddCharacterHandler{}
}

func (h *AddCharacterHandler) Handle(ctx actor.Context, a *MapActor, msg *AddCharacter) {
	if a.Map == nil {
		return
	}
	a.Map.AddPlayer(ctx, msg.Character.GetID(), msg.Character, msg.SpawnPoint, msg.Init)
	msg.Character.ResumeTimers(ctx.Self())
	msg.Character.Listener.OnPartyMemberFieldsChanged(msg.Character)
	if msg.Init {
		msg.Character.SendBuddyLoginSync()
	}
}
