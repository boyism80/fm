package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceCreateHandler struct{}

func (DeliverAllianceCreateHandler) New() *DeliverAllianceCreateHandler {
	return &DeliverAllianceCreateHandler{}
}

func (h *DeliverAllianceCreateHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverAllianceCreate) {
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
	ch.Listener.OnAllianceCreate(ch, msg.Info, msg.Guilds, msg.MembershipGuilds)
}
