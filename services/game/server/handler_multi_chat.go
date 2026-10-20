package server

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/core/rpc"
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
)

type MultiChat struct {
	gs *GameServer
}

func (MultiChat) New(gs *GameServer) *MultiChat {
	return &MultiChat{
		gs: gs,
	}
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
				luax.CallAsync(ctx.ActorContext, root, thread, "on_chat", ch, req.Message, false).OnError(func(err error) {
					log.Printf("Command Lua error: %v", err)
				})
			}
		}
		return nil
	}
	if len(strings.TrimSpace(req.Message)) == 0 {
		return nil
	}
	if ch.Forbidden("chat") {
		ch.Listener().OnMessage(ch, constant.MsgPinkText, "대화 금지 상태이므로 채팅이 불가능합니다.")
		return nil
	}
	mode := pconst.MultiChatMode(req.Type)
	if mode > pconst.MultiChatModeAlliance {
		return nil
	}
	switch mode {
	case pconst.MultiChatModeParty:
		return h.handlePartyMultiChat(ctx, ch, req)
	case pconst.MultiChatModeBuddy:
		return h.handleBuddyMultiChat(ctx, ch, req)
	case pconst.MultiChatModeGuild:
		return h.handleGuildMultiChat(ctx, ch, req)
	case pconst.MultiChatModeAlliance:
		return h.handleAllianceMultiChat(ctx, ch, req)
	default:
		return nil
	}
}

func (h *MultiChat) handleGuildMultiChat(ctx *core.ClientContext, ch *game.Character, req *request.MultiChat) error {
	if ctx.ActorContext == nil {
		return fmt.Errorf("guild multi chat: actor context required")
	}
	guildID, inGuild := ch.Guild().ID()
	if !inGuild {
		return nil
	}
	senderID := ch.GetID()
	senderName := ch.GetName()
	worldID := h.gs.config.WorldId
	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.BroadcastMultiChatReply, error) {
			return h.gs.systems.Get[*rpc.InternalConnection]().Client.BroadcastMultiChat(c, &internal.BroadcastMultiChatRequest{
				WorldId:           worldID,
				MemberId:          guildID,
				SenderCharacterId: senderID,
				ChatMode:          uint32(pconst.MultiChatModeGuild),
				SenderName:        senderName,
				Message:           req.Message,
			})
		},
		func(reply *internal.BroadcastMultiChatReply) error {
			if !reply.GetOk() {
				log.Printf("GuildMultiChat: broadcast failed character=%d code=%v", senderID, reply.GetErrorCode())
			}
			return nil
		},
	).OnError(func(err error) {
		log.Printf("GuildMultiChat async error: %v", err)
	})
	return nil
}

func (h *MultiChat) handleAllianceMultiChat(ctx *core.ClientContext, ch *game.Character, req *request.MultiChat) error {
	if ctx.ActorContext == nil {
		return fmt.Errorf("alliance multi chat: actor context required")
	}
	guildID, inGuild := ch.Guild().ID()
	if !inGuild {
		return nil
	}
	g := h.gs.systems.Get[*game.GuildSystem]().Get(guildID)
	if g == nil {
		return nil
	}
	allianceID, inAlliance := g.GetAllianceID()
	if !inAlliance {
		return nil
	}
	senderID := ch.GetID()
	senderName := ch.GetName()
	worldID := h.gs.config.WorldId
	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.BroadcastMultiChatReply, error) {
			return h.gs.systems.Get[*rpc.InternalConnection]().Client.BroadcastMultiChat(c, &internal.BroadcastMultiChatRequest{
				WorldId:           worldID,
				MemberId:          allianceID,
				SenderCharacterId: senderID,
				ChatMode:          uint32(pconst.MultiChatModeAlliance),
				SenderName:        senderName,
				Message:           req.Message,
			})
		},
		func(reply *internal.BroadcastMultiChatReply) error {
			if !reply.GetOk() {
				log.Printf("AllianceMultiChat: broadcast failed character=%d code=%v", senderID, reply.GetErrorCode())
			}
			return nil
		},
	).OnError(func(err error) {
		log.Printf("AllianceMultiChat async error: %v", err)
	})
	return nil
}

func (h *MultiChat) handlePartyMultiChat(ctx *core.ClientContext, ch *game.Character, req *request.MultiChat) error {
	if ctx.ActorContext == nil {
		return fmt.Errorf("party multi chat: actor context required")
	}
	pid := ch.Party().ID()
	if pid == nil {
		return nil
	}
	memberID := *pid
	if memberID == 0 {
		return nil
	}
	senderID := ch.GetID()
	senderName := ch.GetName()
	worldID := h.gs.config.WorldId
	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.BroadcastMultiChatReply, error) {
			return h.gs.systems.Get[*rpc.InternalConnection]().Client.BroadcastMultiChat(c, &internal.BroadcastMultiChatRequest{
				WorldId:           worldID,
				MemberId:          memberID,
				SenderCharacterId: senderID,
				ChatMode:          uint32(pconst.MultiChatModeParty),
				SenderName:        senderName,
				Message:           req.Message,
			})
		},
		func(reply *internal.BroadcastMultiChatReply) error {
			if !reply.GetOk() {
				log.Printf("PartyMultiChat: broadcast failed character=%d code=%v", senderID, reply.GetErrorCode())
			}
			return nil
		},
	).OnError(func(err error) {
		log.Printf("PartyMultiChat async error: %v", err)
	})
	return nil
}

func (h *MultiChat) handleBuddyMultiChat(ctx *core.ClientContext, ch *game.Character, req *request.MultiChat) error {
	if ctx.ActorContext == nil {
		return fmt.Errorf("buddy multi chat: actor context required")
	}
	recipientIDs := make([]uint32, 0, len(req.Recipients))
	for _, rid := range req.Recipients {
		if rid > 0 {
			recipientIDs = append(recipientIDs, uint32(rid))
		}
	}
	senderID := ch.GetID()
	senderName := ch.GetName()
	worldID := h.gs.config.WorldId
	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.BroadcastMultiChatReply, error) {
			return h.gs.systems.Get[*rpc.InternalConnection]().Client.BroadcastMultiChat(c, &internal.BroadcastMultiChatRequest{
				WorldId:               worldID,
				SenderCharacterId:     senderID,
				ChatMode:              uint32(pconst.MultiChatModeBuddy),
				SenderName:            senderName,
				Message:               req.Message,
				RecipientCharacterIds: recipientIDs,
			})
		},
		func(reply *internal.BroadcastMultiChatReply) error {
			if !reply.GetOk() {
				log.Printf("BuddyMultiChat: broadcast failed character=%d code=%v", senderID, reply.GetErrorCode())
			}
			return nil
		},
	).OnError(func(err error) {
		log.Printf("BuddyMultiChat async error: %v", err)
	})
	return nil
}
