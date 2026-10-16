package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type ReuseOIDHandler struct{}

func (ReuseOIDHandler) New() *ReuseOIDHandler {
	return &ReuseOIDHandler{}
}

func (h *ReuseOIDHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *game.ReuseOID) {
	if msg == nil {
		return
	}
	for _, m := range a.Maps() {
		if m == nil || m.GetMapID() != msg.MapID {
			continue
		}
		if pid := m.LogicActorPID(); pid == nil || !pid.Equal(ctx.Self()) {
			return
		}
		m.HoldOID(msg.OID)
		return
	}
}
