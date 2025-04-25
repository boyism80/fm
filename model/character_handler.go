package model

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
)

func RegisterCharacterHandlers(m *CharacterActor, h *handler.MessageHandler) {
	RegisterLifeHandlers(&m.Life, h)

	handler.RegisterHandler(m, h, onGainExp)
}

func onGainExp(ch *CharacterActor, ctx actor.Context, m *msg.CharacterGainExp) {
	ch.Exp += m.Amount

	fmt.Println("gain exp")
}
