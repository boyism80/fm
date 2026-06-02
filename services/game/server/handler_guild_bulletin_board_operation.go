package server

import (
	"context"
	"log"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
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
	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)

	switch req.Action {
	case pconst.GuildBulletinBoardC2SListThreads:
		promise = async.ThenRPC(promise, func(c context.Context) (*internal.ListGuildBulletinBoardThreadsReply, error) {
			return h.gs.internalClient.ListGuildBulletinBoardThreads(c, &internal.ListGuildBulletinBoardThreadsRequest{
				WorldId:     worldID,
				CharacterId: charID,
				Page:        req.Page,
			})
		}, func(reply *internal.ListGuildBulletinBoardThreadsReply) error {
			if !reply.GetOk() {
				return nil
			}
			_ = ch.Send(threadListFromListReply(reply), types.SEND_POLICY_ENCRYPT)
			return nil
		})
	case pconst.GuildBulletinBoardC2SShowThread:
		promise = async.ThenRPC(promise, func(c context.Context) (*internal.ShowGuildBulletinBoardThreadReply, error) {
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
				_ = ch.Send(showThreadFromDetail(detail), types.SEND_POLICY_ENCRYPT)
			}
			return nil
		})
	case pconst.GuildBulletinBoardC2SWriteThread:
		if req.Edit {
			promise = async.ThenRPC(promise, func(c context.Context) (*internal.UpdateGuildBulletinBoardThreadReply, error) {
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
					_ = ch.Send(showThreadFromDetail(detail), types.SEND_POLICY_ENCRYPT)
				}
				return nil
			})
		} else {
			promise = async.ThenRPC(promise, func(c context.Context) (*internal.CreateGuildBulletinBoardThreadReply, error) {
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
					_ = ch.Send(showThreadFromDetail(detail), types.SEND_POLICY_ENCRYPT)
				}
				if len(reply.GetThreads()) > 0 {
					_ = ch.Send(threadListFromCreateThreadReply(reply), types.SEND_POLICY_ENCRYPT)
				}
				return nil
			})
		}
	case pconst.GuildBulletinBoardC2SDeleteThread:
		promise = async.ThenRPC(promise, func(c context.Context) (*internal.DeleteGuildBulletinBoardThreadReply, error) {
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
		promise = async.ThenRPC(promise, func(c context.Context) (*internal.CreateGuildBulletinBoardReplyReply, error) {
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
				_ = ch.Send(showThreadFromDetail(detail), types.SEND_POLICY_ENCRYPT)
			}
			return nil
		})
	case pconst.GuildBulletinBoardC2SDeleteReply:
		promise = async.ThenRPC(promise, func(c context.Context) (*internal.DeleteGuildBulletinBoardReplyReply, error) {
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
				_ = ch.Send(showThreadFromDetail(detail), types.SEND_POLICY_ENCRYPT)
			}
			return nil
		})
	default:
		return nil
	}

	promise.OnError(func(err error) {
		log.Printf("GuildBulletinBoardOperation async error: %v", err)
	}).Run()
	return nil
}

func (h *GuildBulletinBoardOperation) validateBulletinIcon(ch *entity.Character, icon int32) bool {
	if icon >= pconst.GuildBulletinBoardIconCashMin && icon <= pconst.GuildBulletinBoardIconCashMax {
		itemID := uint32(5290000 + icon - pconst.GuildBulletinBoardIconCashMin)
		return ch.HasItem(itemID)
	}
	return icon >= 0 && icon <= 2
}

func threadListFromListReply(reply *internal.ListGuildBulletinBoardThreadsReply) *response.GuildBulletinBoardThreadList {
	return threadListFromEntries(
		reply.GetThreads(),
		int(reply.GetListStart()),
		int(reply.GetThreadCount()),
		reply.GetNotice(),
	)
}

func threadListFromCreateThreadReply(reply *internal.CreateGuildBulletinBoardThreadReply) *response.GuildBulletinBoardThreadList {
	return threadListFromEntries(
		reply.GetThreads(),
		int(reply.GetListStart()),
		int(reply.GetThreadCount()),
		reply.GetNotice(),
	)
}

func threadListFromEntries(
	threads []*internal.GuildBulletinBoardThreadEntry,
	start int,
	totalCount int,
	notice *internal.GuildBulletinBoardThreadEntry,
) *response.GuildBulletinBoardThreadList {
	entries := make([]response.GuildBulletinBoardThreadEntry, 0, len(threads))
	for _, t := range threads {
		entries = append(entries, bulletinThreadEntryFromProto(t))
	}
	var noticeEntry *response.GuildBulletinBoardThreadEntry
	if notice != nil {
		n := bulletinThreadEntryFromProto(notice)
		noticeEntry = &n
	}
	return &response.GuildBulletinBoardThreadList{
		Start:      start,
		TotalCount: totalCount,
		Notice:     noticeEntry,
		Threads:    entries,
	}
}

func bulletinThreadEntryFromProto(t *internal.GuildBulletinBoardThreadEntry) response.GuildBulletinBoardThreadEntry {
	return response.GuildBulletinBoardThreadEntry{
		LocalThreadID:     t.GetLocalThreadId(),
		PosterCharacterID: t.GetPosterCharacterId(),
		Title:             t.GetTitle(),
		Timestamp:         time.UnixMilli(t.GetTimestampUnixMs()),
		Icon:              t.GetIcon(),
		ReplyCount:        t.GetReplyCount(),
	}
}

func showThreadFromDetail(detail *internal.GuildBulletinBoardThreadDetail) *response.GuildBulletinBoardShowThread {
	replies := make([]response.GuildBulletinBoardReplyEntry, 0, len(detail.GetReplies()))
	for _, r := range detail.GetReplies() {
		replies = append(replies, response.GuildBulletinBoardReplyEntry{
			ReplyID:           r.GetReplyId(),
			PosterCharacterID: r.GetPosterCharacterId(),
			Timestamp:         time.UnixMilli(r.GetTimestampUnixMs()),
			Content:           r.GetContent(),
		})
	}
	return &response.GuildBulletinBoardShowThread{
		LocalThreadID:     detail.GetLocalThreadId(),
		PosterCharacterID: detail.GetPosterCharacterId(),
		Timestamp:         time.UnixMilli(detail.GetTimestampUnixMs()),
		Title:             detail.GetTitle(),
		Body:              detail.GetBody(),
		Icon:              detail.GetIcon(),
		Replies:           replies,
	}
}
