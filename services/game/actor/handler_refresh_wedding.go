package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type RefreshWeddingHandler struct{}

func (RefreshWeddingHandler) New() *RefreshWeddingHandler {
	return &RefreshWeddingHandler{}
}

func (h *RefreshWeddingHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *game.RefreshWedding) {
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Wedding.Refresh(ctx, msg.Event)
}
