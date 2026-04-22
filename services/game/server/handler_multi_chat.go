package server

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	lua "github.com/yuin/gopher-lua"
)

type MultiChat struct {
	gs     *GameServer
	opcode byte
}

func (MultiChat) New(gs *GameServer) *MultiChat {
	return &MultiChat{
		gs:     gs,
		opcode: 0x62,
	}
}

func (h *MultiChat) GetOpcode() byte {
	return h.opcode
}

func (h *MultiChat) Handle(ctx *core.ClientContext, req *request.MultiChat) error {
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
	if h.gs.internalClient == nil {
		return fmt.Errorf("internal client not configured")
	}
	if ctx.ActorContext == nil {
		return fmt.Errorf("multi chat: actor context required")
	}
	mode := pconst.MultiChatMode(req.Type)
	if mode > pconst.MultiChatModeAlliance {
		return nil
	}
	memberID := uint32(0)
	if mode == pconst.MultiChatModeParty {
		pid := ch.GetPartyID()
		if pid == nil {
			return nil
		}
		memberID = *pid
	} else if len(req.Recipients) > 0 {
		first := req.Recipients[0]
		if first > 0 {
			memberID = uint32(first)
		}
	}
	if memberID == 0 {
		return nil
	}
	senderID := ch.GetID()
	senderName := ch.GetName()
	worldID := h.gs.config.WorldId
	async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
		func(c context.Context) (*internal.BroadcastMultiChatReply, error) {
			return h.gs.internalClient.BroadcastMultiChat(c, &internal.BroadcastMultiChatRequest{
				WorldId:           worldID,
				MemberId:          memberID,
				SenderCharacterId: senderID,
				ChatMode:          uint32(mode),
				SenderName:        senderName,
				Message:           req.Message,
			})
		},
		func(reply *internal.BroadcastMultiChatReply) error {
			if !reply.GetOk() {
				log.Printf("MultiChat: broadcast failed character=%d code=%v", senderID, reply.GetErrorCode())
			}
			return nil
		},
	).OnError(func(err error) {
		log.Printf("MultiChat async error: %v", err)
	}).Run()
	return nil
}
