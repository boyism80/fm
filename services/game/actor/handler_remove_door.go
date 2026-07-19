package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
)

type RemoveDoorHandler struct{}

func (RemoveDoorHandler) New() *RemoveDoorHandler {
	return &RemoveDoorHandler{}
}

func (h *RemoveDoorHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *RemoveDoor) {
	if msg == nil {
		return
	}
	for _, m := range a.Maps() {
		if m != nil && m.GetMapID() == msg.MapID {
			m.RemoveDoorByOwnerSkill(msg.OwnerID, constant.SkillID(msg.SkillID), true)
			return
		}
	}
}
