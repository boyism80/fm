package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
)

type RemoveDoorHandler struct{}

func (RemoveDoorHandler) New() *RemoveDoorHandler {
	return &RemoveDoorHandler{}
}

func (h *RemoveDoorHandler) Handle(ctx actor.Context, a *MapActor, msg *RemoveDoor) {
	if a.Map == nil || msg == nil {
		return
	}
	a.Map.RemoveDoorByOwnerSkill(msg.OwnerID, constant.SkillID(msg.SkillID), true)
}
