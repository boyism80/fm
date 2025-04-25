package actor

import (
	"fmt"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
	"github.com/boyism80/fm/msg"
)

func RegisterCharacterHandlers(m *model.Character, h *handler.MessageHandler) {
	RegisterLifeHandlers(&m.Life, h)

	handler.RegisterHandler(m, h, handleGainExp)
}

func handleGainExp(ch *model.Character, ctx protoactor.Context, m *msg.CharacterGainExp) {
	ch.Exp += m.Amount

	fmt.Println("gain exp")
}
