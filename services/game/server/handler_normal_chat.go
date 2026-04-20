package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	lua "github.com/yuin/gopher-lua"
)

type NormalChat struct {
	gs     *GameServer
	opcode byte
}

func (NormalChat) New(gs *GameServer) *NormalChat {
	return &NormalChat{
		gs:     gs,
		opcode: 0x20,
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
				state, result, err := luax.Execute(root, thread, "on_chat", character, req.Message, false)
				if err != nil {
					log.Printf("Command Lua error: %v", err)
				} else if state == lua.ResumeYield {
					// async command path (e.g. save/sleep): completion continues via ResumeLua
					return nil
				} else if result != nil && result.Type() == lua.LTBool && lua.LVAsBool(result) {
					return nil
				} else if result != nil && result.Type() == lua.LTBool && !lua.LVAsBool(result) {
					log.Printf("Unknown command: %s", strings.TrimSpace(req.Message))
				}
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
