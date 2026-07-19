package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceMemberOnlineChangeHandler struct{}

func (DeliverAllianceMemberOnlineChangeHandler) New() *DeliverAllianceMemberOnlineChangeHandler {
	return &DeliverAllianceMemberOnlineChangeHandler{}
}

func (h *DeliverAllianceMemberOnlineChangeHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverAllianceMemberOnlineChange) {
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
	ch.Listener.OnAllianceMemberOnlineChange(ch, msg.AllianceID, msg.GuildID, msg.SubjectID, msg.Online)
}
