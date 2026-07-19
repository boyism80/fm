package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildMemberLeftHandler struct{}

func (DeliverGuildMemberLeftHandler) New() *DeliverGuildMemberLeftHandler {
	return &DeliverGuildMemberLeftHandler{}
}

func (h *DeliverGuildMemberLeftHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildMemberLeft) {
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
	ch.Listener.OnGuildMemberLeft(ch, msg.GuildID, msg.TargetID, msg.TargetName, msg.WasExpelled)
}
