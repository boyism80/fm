package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceGuildLeftHandler struct{}

func (DeliverAllianceGuildLeftHandler) New() *DeliverAllianceGuildLeftHandler {
	return &DeliverAllianceGuildLeftHandler{}
}

func (h *DeliverAllianceGuildLeftHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverAllianceGuildLeft) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceGuildLeft(
		ch,
		msg.Info,
		msg.RemovedGuildID,
		msg.RemovedGuild,
		msg.RemovedMembers,
		msg.Expelled,
		msg.Leaving,
	)
}
