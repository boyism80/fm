package actor

import "github.com/asynkron/protoactor-go/actor"

type HandoffCharacterHandler struct{}

func (HandoffCharacterHandler) New() *HandoffCharacterHandler {
	return &HandoffCharacterHandler{}
}

func (h *HandoffCharacterHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *HandoffCharacter) {
	if msg == nil || msg.Character == nil || msg.TargetMap == nil {
		return
	}
	targetPID := msg.TargetMap.LogicActorPID()
	if targetPID == nil {
		msg.Ticket.Release()
		return
	}
	source := a.GetCharacter(msg.Character.GetID())
	if source == nil {
		msg.Ticket.Release()
		return
	}
	if pid := source.LogicActorPID(); pid == nil || !pid.Equal(ctx.Self()) {
		msg.Ticket.Release()
		return
	}
	msg.Character.Destination = msg.TargetMap
	if err := source.RemovePlayer(msg.Character.GetID()); err != nil {
		msg.Character.Destination = nil
		msg.Ticket.Release()
		return
	}
	ctx.Send(targetPID, &WarpCharacter{
		Character: msg.Character,
		TargetMap: msg.TargetMap,
		Portal:    msg.Portal,
		OnEnter:   msg.OnEnter,
		Ticket:    msg.Ticket,
	})
}
