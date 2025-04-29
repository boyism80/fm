package actor

import (
	"fmt"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/entity"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
)

func RegisterCharacterHandlers(m *entity.Character, h *handler.MessageHandler) {
	RegisterLifeHandlers(&m.Life, h)

	handler.RegisterHandler(m, h, onGainExp)
}

func onGainExp(ch *entity.Character, ctx protoactor.Context, m *msg.CharacterGainExp) {
	ch.Exp += uint32(m.Amount)

	fmt.Println("gain exp")
}
