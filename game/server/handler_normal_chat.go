package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/protocol/request"
	lua "github.com/yuin/gopher-lua"
)

// NormalChat handles normal chat packet requests
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
		if ctx.LogicActorPID != nil {
			root := luax.GetRootLuaState(ctx.LogicActorPID.String())
			if root != nil {
				result, thread, err := luax.CallWithPID(root, ctx.LogicActorPID, "script/command.lua", "on_chat", character, req.Message, false)
				if thread != nil {
					thread.Close()
				}
				if err != nil {
					log.Printf("Command Lua error: %v", err)
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
