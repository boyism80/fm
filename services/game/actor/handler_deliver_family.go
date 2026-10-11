package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type DeliverFamilyHandler struct{}

func (DeliverFamilyHandler) New() *DeliverFamilyHandler {
	return &DeliverFamilyHandler{}
}

func (h *DeliverFamilyHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverFamily) {
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.CharacterID)
	if ch == nil || ch.LoggedOut() {
		return
	}

	switch msg.Event {
	case internal.FamilyEvent_FAMILY_EVENT_INVITE:
		ch.Listener.OnFamilyInvite(ch, msg.SenderID, msg.Name)
	case internal.FamilyEvent_FAMILY_EVENT_INVITE_RESULT:
		ch.Listener.OnFamilyInviteResult(ch, msg.Accepted, msg.Name)
	case internal.FamilyEvent_FAMILY_EVENT_SUMMON:
		ch.Listener.OnFamilySummon(ch, msg.Name, msg.MapName)
	case internal.FamilyEvent_FAMILY_EVENT_JOINED:
		ch.Listener.OnFamilyJoined(ch, msg.Name)
	case internal.FamilyEvent_FAMILY_EVENT_REFRESH:
		ch.Family.Refresh(ctx)
	}
}
