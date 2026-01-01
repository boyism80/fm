package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/protocol/request"
)

type PartySearchStop struct {
	gameServer *GameServer
	opcode     byte
}

func (PartySearchStop) New(gameServer *GameServer) *PartySearchStop {
	return &PartySearchStop{
		gameServer: gameServer,
		opcode:     0xB6,
	}
}

func (h *PartySearchStop) GetOpcode() byte {
	return h.opcode
}

func (h *PartySearchStop) Handle(ctx *core.ClientContext, req *request.PartySearchStop) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		return nil
	}

	return nil
}
