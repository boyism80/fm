package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/protocol/request"
)

// NpcControl handles NPC control packet requests
type NpcControl struct {
	gs     *GameServer
	opcode byte
}

func (NpcControl) New(gs *GameServer) *NpcControl {
	return &NpcControl{
		gs:     gs,
		opcode: 0x9E,
	}
}

func (h *NpcControl) GetOpcode() byte {
	return h.opcode
}

func (h *NpcControl) Handle(ctx *core.ClientContext, req *request.NpcAction) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	character.Listener.OnNpcAction(character, req.Bytes)

	return nil
}
