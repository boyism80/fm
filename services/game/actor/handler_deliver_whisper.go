package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverWhisperHandler struct{}

func (DeliverWhisperHandler) New() *DeliverWhisperHandler {
	return &DeliverWhisperHandler{}
}

func (h *DeliverWhisperHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverWhisper) {
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
