package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type AddCharacterHandler struct{}

func (AddCharacterHandler) New() *AddCharacterHandler {
	return &AddCharacterHandler{}
}

func (h *AddCharacterHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *AddCharacter) {
	if msg == nil || msg.Character == nil || msg.TargetMap == nil {
		return
	}
	msg.TargetMap.AddPlayer(ctx, msg.Character.GetID(), msg.Character, msg.SpawnPoint, msg.Init)
	msg.Character.FinishMove(msg.TargetMap)
	if msg.Character.LoggedOut() {
		_ = msg.TargetMap.LogoutPlayer(ctx, msg.Character.GetID())
		return
	}
	msg.Character.ResumeTimers(ctx.Self())
	msg.Character.Listener.OnPartyMemberFieldsChanged(msg.Character)
	msg.Character.Wedding.NotifySpouseMap(ctx, msg.Init)
	if msg.Init {
		msg.Character.SendBuddyLoginSync()
		msg.Character.Duey.CheckArrivals(ctx)
	}
}
