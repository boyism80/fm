package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type NormalChat struct {
	gs *GameServer
}

func (NormalChat) New(gs *GameServer) *NormalChat {
	return &NormalChat{
		gs: gs,
	}
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
		if mapInstance := character.GetMap(); mapInstance != nil {
			root := mapInstance.GetLuaRoot()
			if root != nil {
				thread, err := luax.NewThread(root, "script/command.lua")
				if err != nil {
					log.Printf("Command Lua thread error: %v", err)
					return nil
				}
				luax.SetConfiguration(thread, luax.Configuration{
					ActorContext: ctx.ActorContext,
				})
				luax.CallAsync(root, thread, "on_chat", character, req.Message, false).OnError(func(err error) {
					log.Printf("Command Lua error: %v", err)
				})
			}
		}
		return nil
	}

	mapInstance := character.GetMap()
	if mapInstance == nil {
		log.Printf("Map not found for character chat")
		return fmt.Errorf("map not found")
	}

	character.Listener.OnChat(character, req.Message, false, req.DontRecordHistory)

	return nil
}
