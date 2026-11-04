package server

import (
	"log"
	"strings"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	lua "github.com/yuin/gopher-lua"
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
	character := h.gs.character(ctx)
	if character == nil {
		return nil
	}

	// The client keeps chat starting with "/" from players without GM rights, so "!" also runs a command.
	userCommand := strings.HasPrefix(req.Message, "!")
	if strings.HasPrefix(req.Message, "/") || userCommand {
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
				luax.CallAsync(ctx.ActorContext, root, thread, "on_chat", character, req.Message, false).Do(func(vals []lua.LValue) error {
					if userCommand && (len(vals) == 0 || !lua.LVAsBool(vals[0])) {
						character.Chat(req.Message, req.DontRecordHistory)
					}
					return nil
				}).OnError(func(err error) {
					log.Printf("Command Lua error: %v", err)
				})
				return nil
			}
		}
		if !userCommand {
			return nil
		}
	}

	character.Chat(req.Message, req.DontRecordHistory)
	return nil
}
