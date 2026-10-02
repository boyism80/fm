package actor

import "github.com/asynkron/protoactor-go/actor"

type DespawnDoorHandler struct{}

func (DespawnDoorHandler) New() *DespawnDoorHandler {
	return &DespawnDoorHandler{}
}

func (h *DespawnDoorHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DespawnDoor) {
	if msg == nil || msg.Map == nil || msg.Door == nil {
		return
	}
	for _, m := range a.Maps() {
		if m != msg.Map {
			continue
		}
		if msg.Door.Map != m || msg.Door.OID == 0 {
			return
		}
		m.RemoveDoor(msg.Door.OID, msg.Animated)
		return
	}
}
