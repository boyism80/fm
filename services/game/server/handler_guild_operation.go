package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/client"
	gameconst "github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	lua "github.com/yuin/gopher-lua"
)

type GuildOperation struct {
	gs *GameServer
}

func (GuildOperation) New(gs *GameServer) *GuildOperation {
	return &GuildOperation{
		gs: gs,
	}
}

func (h *GuildOperation) resumeGuildCreate(ch *entity.Character, result gameconst.GuildCreateResult) {
	if ch == nil {
		return
	}
	thread := ch.Dialog.Thread()
	if thread == nil {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	if _, _, err := luax.Resume(root, thread, lua.LNumber(result)); err != nil {
		log.Printf("GuildOperation(create): failed to resume npc script character=%d: %v", ch.GetID(), err)
	}
}

func (h *GuildOperation) chargeGuildEmblemChangeCost(ch *entity.Character) (guildEmblemDetail, bool) {
	if ch.Inventory.HasItem(gameconst.GuildEmblemChangeCashItemID) {
		if !ch.Inventory.RemoveByItemIDCount(gameconst.GuildEmblemChangeCashItemID, 1) {
			return guildEmblemDetail{}, false
		}
		return guildEmblemDetail{CashItem: true}, true
	}
	if !ch.Inventory.RemoveMeso(gameconst.GuildEmblemChangeMesoCost) {
		return guildEmblemDetail{}, false
	}
	return guildEmblemDetail{Meso: gameconst.GuildEmblemChangeMesoCost}, true
}

func (h *GuildOperation) Handle(ctx *core.ClientContext, req *request.GuildOperation) error {
	if h.gs.internalClient == nil {
		return fmt.Errorf("internal client not configured")
	}
	if ctx.ActorContext == nil {
		return fmt.Errorf("guild operation: actor context required")
	}

	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return fmt.Errorf("character not found")
	}

	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	switch req.Operation {
	case pconst.GuildC2SCreate:
		if _, inGuild := ch.Guild.ID(); inGuild {
			ch.Listener.OnGuildMessage(ch, pconst.GuildResponseAlreadyInGuild)
			h.resumeGuildCreate(ch, gameconst.GuildCreateResultAlreadyInGuild)
			return nil
		}
		mapInstance := ch.GetMap()
		if mapInstance == nil || mapInstance.Wz == nil || uint32(mapInstance.Wz.ID) != gameconst.GuildEmblemMapID {
			h.resumeGuildCreate(ch, gameconst.GuildCreateResultNotAllowed)
			return nil
		}
		if !ch.Inventory.RemoveMeso(gameconst.GuildCreateMesoCost) {
			h.resumeGuildCreate(ch, gameconst.GuildCreateResultInsufficientMeso)
			return nil
		}

		detail := guildCreateDetail{Name: req.GuildName, Meso: gameconst.GuildCreateMesoCost}
		id := ch.Unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_GUILD_CREATE, 0, detail)
		h.gs.guildSystem.createAsync(ctx.ActorContext, ch, id, detail.Name).Do(func(reply *internal.CreateGuildReply) error {
			ch.Unsettled.Erase(id)
			if !reply.GetOk() {
				if reply.GetErrorCode() == internal.GuildErrorCode_GUILD_ERROR_ALREADY_IN_GUILD {
					ch.Listener.OnGuildMessage(ch, pconst.GuildResponseAlreadyInGuild)
				} else {
					log.Printf("GuildOperation(create): failed character=%d code=%v", charID, reply.GetErrorCode())
				}
				ch.Inventory.AddMeso(detail.Meso)
				h.resumeGuildCreate(ch, gameconst.GuildCreateResultFailed)
				return nil
			}
			guildPb := reply.GetGuild()
			if guildPb == nil {
				log.Printf("GuildOperation(create): ok but missing guild character=%d", charID)
				h.resumeGuildCreate(ch, gameconst.GuildCreateResultFailed)
				return nil
			}
			h.gs.guildSystem.applyGuildFromProto(ch, guildPb)
			h.resumeGuildCreate(ch, gameconst.GuildCreateResultOK)
			return nil
		}).OnError(func(err error) {
			ch.Unsettled.Find(id).InFlight = false
			log.Printf("GuildOperation(create) character=%d id=%d: %v", charID, id, err)
			h.resumeGuildCreate(ch, gameconst.GuildCreateResultFailed)
		})
		return nil
	case pconst.GuildC2SInvite:
		guildID, inGuild := ch.Guild.ID()
		g := h.gs.GetGuildSystem().Get(guildID)
		if !inGuild || g == nil || !g.CanInvite(charID) {
			return nil
		}
		targetName := req.TargetName
		if targetName == "" {
			return nil
		}
		if h.gs.characterRuntime == nil {
			return nil
		}
		targetID, ok := h.gs.characterRuntime.NameToCharacterID(targetName)
		if !ok || targetID == 0 {
			ch.Listener.OnGuildMessage(ch, pconst.GuildResponseNotInChannel)
			return nil
		}
		if mapInstance := ch.GetMap(); mapInstance != nil {
			if target := mapInstance.GetPlayer(targetID); target != nil {
				switch target.Guild.Invite(guildID, ch.GetName()) {
				case entity.ErrGuildAlreadyJoined:
					ch.Listener.OnGuildMessage(ch, pconst.GuildResponseAlreadyInGuild)
				case entity.ErrGuildInviteBusy:
					ch.Listener.OnMessage(ch, gameconst.MsgPinkText, gameconst.GuildInviteTargetBusyMessage)
				}
				return nil
			}
		}
		h.gs.Deliver(nil, targetID, &g_actor.DeliverGuildInvite{
			CharacterID:        targetID,
			InviterCharacterID: charID,
			GuildID:            guildID,
			InviterName:        ch.GetName(),
		})
		return nil
	case pconst.GuildC2SAcceptInvite:
		if req.GuildID == 0 || req.CharacterID != charID {
			return nil
		}
		if ch.Guild.AcceptInvite(req.GuildID) != nil {
			return nil
		}

		member := ch.ToProtoGuildMember(worldID, int32(h.gs.config.ChannelId), internal.GuildMemberRank_GUILD_MEMBER_RANK_NEW)
		promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise.ThenRPC(func(c context.Context) (*internal.AcceptGuildInviteReply, error) {
			return h.gs.internalClient.AcceptGuildInvite(c, &internal.AcceptGuildInviteRequest{
				WorldId: worldID,
				GuildId: req.GuildID,
				Member:  member,
			})
		}, func(reply *internal.AcceptGuildInviteReply) error {
			if !reply.GetOk() {
				if reply.GetErrorCode() == internal.GuildErrorCode_GUILD_ERROR_GUILD_FULL {
					ch.Listener.OnMessage(ch, gameconst.MsgPopup, gameconst.GuildFullMessage)
				} else if reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_ALREADY_IN_GUILD {
					log.Printf("GuildOperation(accept invite): failed character=%d guild=%d code=%v", charID, req.GuildID, reply.GetErrorCode())
				}
				return nil
			}
			guildPb := reply.GetGuild()
			if guildPb == nil {
				log.Printf("GuildOperation(accept invite): ok but missing guild character=%d", charID)
				return nil
			}
			h.gs.guildSystem.applyGuildFromProto(ch, guildPb)
			return nil
		})
		promise.OnError(func(err error) {
			log.Printf("GuildOperation(accept invite) async error: %v", err)
		})
		return nil
	case pconst.GuildC2SLeave:
		if req.CharacterID != charID {
			return nil
		}
		if req.CharacterName != ch.GetName() {
			return nil
		}
		if _, inGuild := ch.Guild.ID(); !inGuild {
			return nil
		}
		promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise.ThenRPC(func(c context.Context) (*internal.LeaveGuildReply, error) {
			return h.gs.internalClient.LeaveGuild(c, &internal.LeaveGuildRequest{
				WorldId:     worldID,
				CharacterId: charID,
			})
		}, func(reply *internal.LeaveGuildReply) error {
			if !reply.GetOk() {
				if reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD &&
					reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_LEADER_CANNOT_LEAVE {
					log.Printf("GuildOperation(leave): failed character=%d code=%v", charID, reply.GetErrorCode())
				}
				return nil
			}
			return nil
		})
		promise.OnError(func(err error) {
			log.Printf("GuildOperation(leave) async error: %v", err)
		})
		return nil
	case pconst.GuildC2SExpel:
		targetID := req.CharacterID
		if targetID == 0 {
			return nil
		}
		guildID, inGuild := ch.Guild.ID()
		g := h.gs.GetGuildSystem().Get(guildID)
		if !inGuild || g == nil || !g.CanInvite(charID) {
			return nil
		}
		promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise.ThenRPC(func(c context.Context) (*internal.ExpelGuildReply, error) {
			return h.gs.internalClient.ExpelGuild(c, &internal.ExpelGuildRequest{
				WorldId:              worldID,
				RequesterCharacterId: charID,
				TargetCharacterId:    targetID,
			})
		}, func(reply *internal.ExpelGuildReply) error {
			if !reply.GetOk() {
				if reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD &&
					reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_AUTHORIZED &&
					reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_CANNOT_EXPEL_TARGET &&
					reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_TARGET_NOT_IN_GUILD {
					log.Printf("GuildOperation(expel): failed requester=%d target=%d code=%v", charID, targetID, reply.GetErrorCode())
				}
				return nil
			}
			return nil
		})
		promise.OnError(func(err error) {
			log.Printf("GuildOperation(expel) async error: %v", err)
		})
		return nil
	case pconst.GuildC2SChangeRankTitles:
		guildID, inGuild := ch.Guild.ID()
		g := h.gs.GetGuildSystem().Get(guildID)
		if !inGuild || g == nil || !g.IsGuildMaster(charID) {
			return nil
		}
		promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise.ThenRPC(func(c context.Context) (*internal.ChangeGuildRankTitlesReply, error) {
			return h.gs.internalClient.ChangeGuildRankTitles(c, &internal.ChangeGuildRankTitlesRequest{
				WorldId:     worldID,
				CharacterId: charID,
				RankTitles: []string{
					req.RankTitles[0],
					req.RankTitles[1],
					req.RankTitles[2],
					req.RankTitles[3],
					req.RankTitles[4],
				},
			})
		}, func(reply *internal.ChangeGuildRankTitlesReply) error {
			if !reply.GetOk() {
				if reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD &&
					reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_AUTHORIZED &&
					reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_INVALID_RANK_TITLES {
					log.Printf("GuildOperation(change rank titles): failed character=%d code=%v", charID, reply.GetErrorCode())
				}
				return nil
			}
			return nil
		})
		promise.OnError(func(err error) {
			log.Printf("GuildOperation(change rank titles) async error: %v", err)
		})
		return nil
	case pconst.GuildC2SChangeMemberRank:
		newRank := req.NewMemberRank
		if newRank <= 1 || newRank > 5 {
			return nil
		}
		targetID := req.CharacterID
		if targetID == 0 {
			return nil
		}
		guildID, inGuild := ch.Guild.ID()
		g := h.gs.GetGuildSystem().Get(guildID)
		if !inGuild || g == nil || !g.CanInvite(charID) {
			return nil
		}
		if newRank <= 2 && !g.IsGuildMaster(charID) {
			return nil
		}
		protoRank := internal.GuildMemberRank(newRank)
		promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise.ThenRPC(func(c context.Context) (*internal.ChangeGuildMemberRankReply, error) {
			return h.gs.internalClient.ChangeGuildMemberRank(c, &internal.ChangeGuildMemberRankRequest{
				WorldId:              worldID,
				RequesterCharacterId: charID,
				TargetCharacterId:    targetID,
				NewRank:              protoRank,
			})
		}, func(reply *internal.ChangeGuildMemberRankReply) error {
			if !reply.GetOk() {
				if reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD &&
					reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_AUTHORIZED &&
					reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_INVALID_MEMBER_RANK &&
					reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_TARGET_NOT_IN_GUILD {
					log.Printf("GuildOperation(change member rank): failed requester=%d target=%d code=%v", charID, targetID, reply.GetErrorCode())
				}
				return nil
			}
			return nil
		})
		promise.OnError(func(err error) {
			log.Printf("GuildOperation(change member rank) async error: %v", err)
		})
		return nil
	case pconst.GuildC2SChangeEmblem:
		guildID, inGuild := ch.Guild.ID()
		g := h.gs.GetGuildSystem().Get(guildID)
		mapInstance := ch.GetMap()
		if !inGuild || g == nil || !g.IsGuildMaster(charID) ||
			mapInstance == nil || mapInstance.Wz == nil || uint32(mapInstance.Wz.ID) != gameconst.GuildEmblemMapID {
			return nil
		}
		detail, paid := h.chargeGuildEmblemChangeCost(ch)
		if !paid {
			ch.Listener.OnMessage(ch, gameconst.MsgPopup, gameconst.GuildEmblemChangeInsufficientCostMessage)
			return nil
		}
		detail.Logo = &internal.GuildLogo{
			Logo:        uint32(req.Logo),
			LogoColor:   uint32(req.LogoColor),
			LogoBg:      uint32(req.LogoBG),
			LogoBgColor: uint32(req.LogoBGColor),
		}
		id := ch.Unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_GUILD_EMBLEM, uint64(guildID), detail)
		h.gs.guildSystem.changeEmblemAsync(ctx.ActorContext, ch, id, detail.Logo).Do(func(reply *internal.ChangeGuildEmblemReply) error {
			ch.Unsettled.Erase(id)
			if !reply.GetOk() {
				h.gs.guildSystem.refundEmblem(ch, detail)
				if reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD &&
					reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_AUTHORIZED {
					log.Printf("GuildOperation(change emblem): failed character=%d code=%v", charID, reply.GetErrorCode())
				}
			}
			return nil
		}).OnError(func(err error) {
			ch.Unsettled.Find(id).InFlight = false
			log.Printf("GuildOperation(change emblem) character=%d id=%d: %v", charID, id, err)
		})
		return nil
	case pconst.GuildC2SChangeNotice:
		guildID, inGuild := ch.Guild.ID()
		g := h.gs.GetGuildSystem().Get(guildID)
		if !inGuild || g == nil || !g.CanInvite(charID) {
			return nil
		}
		promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise.ThenRPC(func(c context.Context) (*internal.ChangeGuildNoticeReply, error) {
			return h.gs.internalClient.ChangeGuildNotice(c, &internal.ChangeGuildNoticeRequest{
				WorldId:     worldID,
				CharacterId: charID,
				Notice:      req.Notice,
			})
		}, func(reply *internal.ChangeGuildNoticeReply) error {
			if !reply.GetOk() {
				if reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD &&
					reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_AUTHORIZED &&
					reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_INVALID_NOTICE {
					log.Printf("GuildOperation(change notice): failed character=%d code=%v", charID, reply.GetErrorCode())
				}
				return nil
			}
			return nil
		})
		promise.OnError(func(err error) {
			log.Printf("GuildOperation(change notice) async error: %v", err)
		})
		return nil
	default:
		return nil
	}
}
