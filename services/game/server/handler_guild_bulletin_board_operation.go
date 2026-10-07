package server

import (
	"context"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
)

type GuildBulletinBoardOperation struct {
	gs *GameServer
}

func (GuildBulletinBoardOperation) New(gs *GameServer) *GuildBulletinBoardOperation {
	return &GuildBulletinBoardOperation{
		gs: gs,
	}
}

func (h *GuildBulletinBoardOperation) Handle(ctx *core.ClientContext, req *request.GuildBulletinBoardOperation) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return nil
	}
	_, inGuild := ch.GetGuildID()
	if !inGuild {
		return nil
	}
	if h.gs.internalClient == nil {
		return nil
	}

	if req.Action == pconst.GuildBulletinBoardC2SWriteThread {
		if !h.validateBulletinIcon(ch, req.Icon) {
			return nil
		}
	}

	worldID := h.gs.config.WorldId
	charID := ch.GetID()
	promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)

	switch req.Action {
	case pconst.GuildBulletinBoardC2SListThreads:
		promise.ThenRPC(func(c context.Context) (*internal.ListGuildBulletinBoardThreadsReply, error) {
			return h.gs.internalClient.ListGuildBulletinBoardThreads(c, &internal.ListGuildBulletinBoardThreadsRequest{
				WorldId:     worldID,
				CharacterId: charID,
				Page:        req.Page,
			})
		}, func(reply *internal.ListGuildBulletinBoardThreadsReply) error {
			if !reply.GetOk() {
				return nil
			}
			ch.Listener.OnGuildBulletinThreadList(ch, reply.GetThreads(), int(reply.GetListStart()), int(reply.GetThreadCount()), reply.GetNotice())
			return nil
		})
	case pconst.GuildBulletinBoardC2SShowThread:
		promise.ThenRPC(func(c context.Context) (*internal.ShowGuildBulletinBoardThreadReply, error) {
			return h.gs.internalClient.ShowGuildBulletinBoardThread(c, &internal.ShowGuildBulletinBoardThreadRequest{
				WorldId:       worldID,
				CharacterId:   charID,
				LocalThreadId: req.LocalThreadID,
			})
		}, func(reply *internal.ShowGuildBulletinBoardThreadReply) error {
			if !reply.GetOk() {
				return nil
			}
			if detail := reply.GetThread(); detail != nil {
				ch.Listener.OnGuildBulletinThread(ch, detail)
			}
			return nil
		})
	case pconst.GuildBulletinBoardC2SWriteThread:
		if req.Edit {
			promise.ThenRPC(func(c context.Context) (*internal.UpdateGuildBulletinBoardThreadReply, error) {
				return h.gs.internalClient.UpdateGuildBulletinBoardThread(c, &internal.UpdateGuildBulletinBoardThreadRequest{
					WorldId:       worldID,
					CharacterId:   charID,
					LocalThreadId: req.LocalThreadID,
					Title:         req.Title,
					Body:          req.Text,
					Icon:          req.Icon,
				})
			}, func(reply *internal.UpdateGuildBulletinBoardThreadReply) error {
				if !reply.GetOk() {
					return nil
				}
				if detail := reply.GetThread(); detail != nil {
					ch.Listener.OnGuildBulletinThread(ch, detail)
				}
				return nil
			})
		} else {
			promise.ThenRPC(func(c context.Context) (*internal.CreateGuildBulletinBoardThreadReply, error) {
				return h.gs.internalClient.CreateGuildBulletinBoardThread(c, &internal.CreateGuildBulletinBoardThreadRequest{
					WorldId:     worldID,
					CharacterId: charID,
					Notice:      req.Notice,
					Title:       req.Title,
					Body:        req.Text,
					Icon:        req.Icon,
				})
			}, func(reply *internal.CreateGuildBulletinBoardThreadReply) error {
				if !reply.GetOk() {
					return nil
				}
				if detail := reply.GetThread(); detail != nil {
					ch.Listener.OnGuildBulletinThread(ch, detail)
				}
				if len(reply.GetThreads()) > 0 {
					ch.Listener.OnGuildBulletinThreadList(ch, reply.GetThreads(), int(reply.GetListStart()), int(reply.GetThreadCount()), reply.GetNotice())
				}
				return nil
			})
		}
	case pconst.GuildBulletinBoardC2SDeleteThread:
		promise.ThenRPC(func(c context.Context) (*internal.DeleteGuildBulletinBoardThreadReply, error) {
			return h.gs.internalClient.DeleteGuildBulletinBoardThread(c, &internal.DeleteGuildBulletinBoardThreadRequest{
				WorldId:       worldID,
				CharacterId:   charID,
				LocalThreadId: req.LocalThreadID,
			})
		}, func(reply *internal.DeleteGuildBulletinBoardThreadReply) error {
			_ = reply
			return nil
		})
	case pconst.GuildBulletinBoardC2SWriteReply:
		promise.ThenRPC(func(c context.Context) (*internal.CreateGuildBulletinBoardReplyReply, error) {
			return h.gs.internalClient.CreateGuildBulletinBoardReply(c, &internal.CreateGuildBulletinBoardReplyRequest{
				WorldId:       worldID,
				CharacterId:   charID,
				LocalThreadId: req.LocalThreadID,
				Content:       req.Text,
			})
		}, func(reply *internal.CreateGuildBulletinBoardReplyReply) error {
			if !reply.GetOk() {
				return nil
			}
			if detail := reply.GetThread(); detail != nil {
				ch.Listener.OnGuildBulletinThread(ch, detail)
			}
			return nil
		})
	case pconst.GuildBulletinBoardC2SDeleteReply:
		promise.ThenRPC(func(c context.Context) (*internal.DeleteGuildBulletinBoardReplyReply, error) {
			return h.gs.internalClient.DeleteGuildBulletinBoardReply(c, &internal.DeleteGuildBulletinBoardReplyRequest{
				WorldId:       worldID,
				CharacterId:   charID,
				LocalThreadId: req.LocalThreadID,
				ReplyId:       req.ReplyID,
			})
		}, func(reply *internal.DeleteGuildBulletinBoardReplyReply) error {
			if !reply.GetOk() {
				return nil
			}
			if detail := reply.GetThread(); detail != nil {
				ch.Listener.OnGuildBulletinThread(ch, detail)
			}
			return nil
		})
	default:
		return nil
	}

	promise.OnError(func(err error) {
		log.Printf("GuildBulletinBoardOperation async error: %v", err)
	})
	return nil
}

func (h *GuildBulletinBoardOperation) validateBulletinIcon(ch *entity.Character, icon int32) bool {
	if icon >= pconst.GuildBulletinBoardIconCashMin && icon <= pconst.GuildBulletinBoardIconCashMax {
		itemID := uint32(5290000 + icon - pconst.GuildBulletinBoardIconCashMin)
		return ch.Inventory.HasItem(itemID)
	}
	return icon >= 0 && icon <= 2
}
