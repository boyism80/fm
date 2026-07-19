package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type DeliverPartyUpdateJoinHandler struct{}

func (DeliverPartyUpdateJoinHandler) New() *DeliverPartyUpdateJoinHandler {
	return &DeliverPartyUpdateJoinHandler{}
}

func (h *DeliverPartyUpdateJoinHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverPartyUpdateJoin) {
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
	ch.Listener.OnPartyUpdateJoin(ch, msg.ForChannel, msg.PartyID, msg.JoinName, msg.LeaderID, msg.Members)
	if ch.GetID() != msg.JoinCharacterID {
		return
	}
	ch.Listener.OnPartyMemberHPChanged(ch, nil)
	for _, obj := range m.GetObjects(constant.ObjectTypeCharacter) {
		peer, ok := obj.(*entity.Character)
		if !ok || peer == nil || peer.GetID() == ch.GetID() {
			continue
		}
		pPeer := peer.GetPartyID()
		if pPeer == nil || *pPeer != msg.PartyID {
			continue
		}
		peer.Listener.OnPartyMemberHPChanged(peer, ch)
	}
}
