package handler

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/model"
	"github.com/boyism80/fm/msg"
)

func handleGainExp(ch *model.Character, ctx actor.Context, m interface{}) {
	gain := m.(*msg.CharacterGainExp)
	ch.Exp += gain.Amount

	fmt.Println("gain exp")
}

func RegisterCharacterHandlers(m *model.Character, handler *MessageHandler) {
	RegisterLifeHandlers(&m.Life, handler)

	RegisterHandler(m, &msg.CharacterGainExp{}, handler, handleGainExp)
}
