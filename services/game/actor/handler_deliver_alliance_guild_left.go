package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceGuildLeftHandler struct{}

func (DeliverAllianceGuildLeftHandler) New() *DeliverAllianceGuildLeftHandler {
	return &DeliverAllianceGuildLeftHandler{}
}

func (h *DeliverAllianceGuildLeftHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverAllianceGuildLeft) {
	if msg == nil {
		return
	}
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.CharacterID)
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
