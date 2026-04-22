package server

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	lua "github.com/yuin/gopher-lua"
)

type PartyChat struct {
	gs     *GameServer
	opcode byte
}

func (PartyChat) New(gs *GameServer) *PartyChat {
	return &PartyChat{
		gs:     gs,
		opcode: 0x62,
	}
}

func (h *PartyChat) GetOpcode() byte {
	return h.opcode
}

func (h *PartyChat) Handle(ctx *core.ClientContext, req *request.PartyChat) error {
	gc, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}
	ch := gc.GetCharacter()
	if ch == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}
	if req.NumRecipients <= 0 {
		return nil
	}
	if strings.HasPrefix(req.Message, "/") {
		if mapInstance := ch.GetMap(); mapInstance != nil {
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
				state, result, err := luax.Execute(root, thread, "on_chat", ch, req.Message, false)
				if err != nil {
					log.Printf("Command Lua error: %v", err)
				} else if state == lua.ResumeYield {
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
	if len(strings.TrimSpace(req.Message)) == 0 {
		return nil
	}
	if req.Type != 1 {
		return nil
	}
	if h.gs.internalClient == nil {
		return fmt.Errorf("internal client not configured")
	}
	if ctx.ActorContext == nil {
		return fmt.Errorf("party chat: actor context required")
	}
	pid := ch.GetPartyID()
	if pid == nil {
		return nil
	}
	senderID := ch.GetID()
	senderName := ch.GetName()
	worldID := h.gs.config.WorldId
	async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
		func(c context.Context) (*internal.BroadcastPartyChatReply, error) {
			return h.gs.internalClient.BroadcastPartyChat(c, &internal.BroadcastPartyChatRequest{
				WorldId:           worldID,
				PartyId:           *pid,
				SenderCharacterId: senderID,
				ChatMode:          1,
				SenderName:        senderName,
				Message:           req.Message,
			})
		},
		func(reply *internal.BroadcastPartyChatReply) error {
			if !reply.GetOk() {
				log.Printf("PartyChat: broadcast failed character=%d code=%v", senderID, reply.GetErrorCode())
			}
			return nil
		},
	).OnError(func(err error) {
		log.Printf("PartyChat async error: %v", err)
	}).Run()
	return nil
}
