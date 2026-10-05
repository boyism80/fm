package actor

import (
	"log"

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
	pid := msg.TargetMap.LogicActorPID()
	if pid == nil {
		msg.Character.Destination = nil
		msg.Ticket.Release()
		return
	}
	// A state machine attached or detached while the message was in flight; the ticket keeps the map open until it arrives.
	if pid.Equal(ctx.Self()) == false {
		ctx.Send(pid, msg)
		return
	}

	err := msg.TargetMap.AddPlayer(ctx, msg.Character.GetID(), msg.Character, msg.Portal, false)
	msg.Character.Destination = nil
	msg.Ticket.Release()
	if err != nil {
		log.Printf("warp character %d to map %d: %v", msg.Character.GetID(), msg.TargetMap.GetMapID(), err)
		return
	}
	if msg.Character.LoggedOut() {
		_ = msg.TargetMap.LogoutPlayer(msg.Character.GetID())
		return
	}
	if msg.TargetMap.Closing() {
		a.GameWorld.GetMapSystem().CloseInstance(msg.TargetMap)
	}
	msg.Character.ResumeTimers(ctx.Self())
	msg.Character.Listener.OnPartyMemberFieldsChanged(msg.Character)
	if msg.OnEnter != nil {
		msg.OnEnter(ctx)
	}
}
