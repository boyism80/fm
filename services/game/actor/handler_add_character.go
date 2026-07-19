package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type AddCharacterHandler struct{}

func (AddCharacterHandler) New() *AddCharacterHandler {
	return &AddCharacterHandler{}
}

func (h *AddCharacterHandler) Handle(ctx actor.Context, a *MapActor, msg *AddCharacter) {
	if msg == nil || msg.Character == nil || msg.TargetMap == nil {
		return
	}
	pid := msg.TargetMap.GetActorPID()
	if pid == nil || !pid.Equal(ctx.Self()) {
		return
	}
	msg.TargetMap.AddPlayer(ctx, msg.Character.GetID(), msg.Character, msg.SpawnPoint, msg.Init)
	msg.Character.ResumeTimers(ctx.Self())
	msg.Character.Listener.OnPartyMemberFieldsChanged(msg.Character)
	if msg.Init {
		msg.Character.SendBuddyLoginSync()
	}
}
