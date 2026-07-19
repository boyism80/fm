package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceCreateHandler struct{}

func (DeliverAllianceCreateHandler) New() *DeliverAllianceCreateHandler {
	return &DeliverAllianceCreateHandler{}
}

func (h *DeliverAllianceCreateHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverAllianceCreate) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceCreate(ch, msg.Info, msg.Guilds, msg.MembershipGuilds)
}
