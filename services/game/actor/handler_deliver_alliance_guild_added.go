package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/protocol/dto"
)

type DeliverAllianceGuildAddedHandler struct{}

func (DeliverAllianceGuildAddedHandler) New() *DeliverAllianceGuildAddedHandler {
	return &DeliverAllianceGuildAddedHandler{}
}

func (h *DeliverAllianceGuildAddedHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverAllianceGuildAdded) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	var membershipGuild *dto.AllianceMembershipChangeGuild
	if msg.HasMembership {
		membershipGuild = &msg.MembershipGuild
	}
	ch.Listener.OnAllianceGuildAdded(
		ch,
		msg.Info,
		msg.Guilds,
		msg.NewGuildID,
		msg.AddedGuild,
		msg.Members,
		msg.Joining,
		membershipGuild,
	)
}
