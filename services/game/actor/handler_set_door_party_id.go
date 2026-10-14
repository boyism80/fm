package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/domain"
)

type SetDoorPartyIDHandler struct{}

func (SetDoorPartyIDHandler) New() *SetDoorPartyIDHandler {
	return &SetDoorPartyIDHandler{}
}

func (h *SetDoorPartyIDHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *domain.SetDoorPartyID) {
	if msg == nil || msg.Map == nil {
		return
	}
	for _, m := range a.Maps() {
		if m != msg.Map {
			continue
		}
		door := m.FindDoorByKey(msg.Key)
		if door == nil {
			return
		}
		door.PartyID = msg.PartyID
		return
	}
}
