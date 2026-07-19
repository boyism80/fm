package actor

import "github.com/asynkron/protoactor-go/actor"

type HandoffCharacterHandler struct{}

func (HandoffCharacterHandler) New() *HandoffCharacterHandler {
	return &HandoffCharacterHandler{}
}

func (h *HandoffCharacterHandler) Handle(ctx actor.Context, a *MapActor, msg *HandoffCharacter) {
	if msg == nil || msg.Character == nil || msg.TargetMap == nil {
		return
	}
	targetPID := msg.TargetMap.GetActorPID()
	if targetPID == nil {
		return
	}
	source := a.MapForCharacter(msg.Character.GetID())
	if source == nil {
		return
	}
	if pid := source.GetActorPID(); pid == nil || !pid.Equal(ctx.Self()) {
		return
	}
	if err := source.RemovePlayer(msg.Character.GetID()); err != nil {
		return
	}
	ctx.Send(targetPID, &WarpCharacter{
		Character: msg.Character,
		TargetMap: msg.TargetMap,
		Portal:    msg.Portal,
	})
}
