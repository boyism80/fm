package model

import (
	"fmt"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
)

func RegisterCharacterHandlers(m *CharacterActor, h *handler.MessageHandler) {
	RegisterLifeHandlers(&m.Life, h)

	handler.RegisterHandler(m, h, onGainExp)
}

func onGainExp(ch *CharacterActor, ctx protoactor.Context, m *msg.CharacterGainExp) {
	ch.Exp += m.Amount

	fmt.Println("gain exp")
}
