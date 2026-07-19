package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceMemberOnlineChangeHandler struct{}

func (DeliverAllianceMemberOnlineChangeHandler) New() *DeliverAllianceMemberOnlineChangeHandler {
	return &DeliverAllianceMemberOnlineChangeHandler{}
}

func (h *DeliverAllianceMemberOnlineChangeHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverAllianceMemberOnlineChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceMemberOnlineChange(ch, msg.AllianceID, msg.GuildID, msg.SubjectID, msg.Online)
}
