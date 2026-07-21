package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type WarpCharacterHandler struct{}

func (WarpCharacterHandler) New() *WarpCharacterHandler {
	return &WarpCharacterHandler{}
}

func (h *WarpCharacterHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *WarpCharacter) {
	if msg == nil || msg.Character == nil || msg.TargetMap == nil {
		return
	}
	if pid := msg.TargetMap.LogicActorPID(); pid == nil || !pid.Equal(ctx.Self()) {
		return
	}
	msg.TargetMap.AddPlayer(ctx, msg.Character.GetID(), msg.Character, msg.Portal, false)
	msg.Character.ResumeTimers(ctx.Self())
	msg.Character.Listener.OnPartyMemberFieldsChanged(msg.Character)
}
