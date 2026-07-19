package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildNewMemberHandler struct{}

func (DeliverGuildNewMemberHandler) New() *DeliverGuildNewMemberHandler {
	return &DeliverGuildNewMemberHandler{}
}

func (h *DeliverGuildNewMemberHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildNewMember) {
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
	ch.Listener.OnGuildNewMember(ch, msg.GuildID, msg.Member)
}
