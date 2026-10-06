package actor

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
)

type HandoffCharacterHandler struct{}

func (HandoffCharacterHandler) New() *HandoffCharacterHandler {
	return &HandoffCharacterHandler{}
}

func (h *HandoffCharacterHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *HandoffCharacter) {
	if msg == nil || msg.Character == nil || msg.TargetMap == nil {
		return
	}
	err := a.GameWorld.GetMapSystem().Warp(ctx, msg.Character, msg.TargetMap, msg.Portal, msg.OnEnter)
	msg.Ticket.Release()
	if err != nil {
		log.Printf("handoff character %d to map %d: %v", msg.Character.GetID(), msg.TargetMap.GetMapID(), err)
	}
}
