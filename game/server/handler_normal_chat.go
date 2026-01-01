package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/protocol/request"
)

// NormalChat handles normal chat packet requests
type NormalChat struct {
	gameServer *GameServer
	opcode     byte
}

func (NormalChat) New(gameServer *GameServer) *NormalChat {
	return &NormalChat{
		gameServer: gameServer,
		opcode:     0x20,
	}
}

func (h *NormalChat) GetOpcode() byte {
	return h.opcode
}

func (h *NormalChat) Handle(ctx *core.ClientContext, req *request.NormalChat) error {
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

	if strings.HasPrefix(req.Message, "/") {
		params := strings.Split(strings.TrimPrefix(req.Message, "/"), " ")
		err := h.gameServer.commandHandler.Handle(client, params...)
		if err != nil {
			log.Printf("Command error: %v", err)
		}
		return nil
	}

	mapInstance := h.gameServer.GetMap(character.GetMap())
	if mapInstance == nil {
		log.Printf("Map %d not found for character chat", character.GetMap())
		return fmt.Errorf("map %d not found", character.GetMap())
	}

	character.Listener.OnChat(req.Message, false, req.DontRecordHistory)

	return nil
}
