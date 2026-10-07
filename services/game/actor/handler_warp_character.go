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
	if msg.Character.Destination() != msg.TargetMap {
		log.Printf("warp character %d to map %d: character is moving elsewhere", msg.Character.GetID(), msg.TargetMap.GetMapID())
		msg.Ticket.Release()
		return
	}

	err := msg.TargetMap.AddPlayer(ctx, msg.Character.GetID(), msg.Character, msg.Portal, false)
	msg.Character.FinishMove(msg.TargetMap)
	msg.Ticket.Release()
	if err != nil {
		log.Printf("warp character %d to map %d: %v", msg.Character.GetID(), msg.TargetMap.GetMapID(), err)
		return
	}
	if msg.Character.LoggedOut() {
		_ = msg.TargetMap.LogoutPlayer(ctx, msg.Character.GetID())
		return
	}
	if msg.TargetMap.Closing() {
		a.GameWorld.GetMapSystem().CloseInstance(msg.TargetMap)
	}
	msg.Character.ResumeTimers(ctx.Self())
	msg.Character.Listener.OnPartyMemberFieldsChanged(msg.Character)
	msg.Character.NotifySpouseMap(ctx, false)
	if msg.OnEnter != nil {
		msg.OnEnter(ctx)
	}
}
