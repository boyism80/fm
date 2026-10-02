package actor

import "github.com/asynkron/protoactor-go/actor"

type DespawnDoorHandler struct{}

func (DespawnDoorHandler) New() *DespawnDoorHandler {
	return &DespawnDoorHandler{}
}

func (h *DespawnDoorHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DespawnDoor) {
	if msg == nil || msg.Map == nil {
		return
	}
	for _, m := range a.Maps() {
		if m != msg.Map {
			continue
		}
		m.RemoveDoorByKey(msg.Key, msg.Animated, msg.NotifyCounterpart)
		return
	}
}
