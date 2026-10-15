package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type ReceiveWhisperHandler struct{}

func (ReceiveWhisperHandler) New() *ReceiveWhisperHandler {
	return &ReceiveWhisperHandler{}
}

func (h *ReceiveWhisperHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *game.ReceiveWhisper) {
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.CharacterID)
	if ch == nil || ch.LoggedOut() {
		return
	}
	ch.ReceiveWhisper(ctx, msg.Whisper)
}
