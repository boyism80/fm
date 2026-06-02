package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/client"
	gameconst "github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
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
	thread := ch.GetDialog()
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
	if _, err := luax.Resume(root, thread, lua.LNumber(result)); err != nil {
		log.Printf("GuildOperation(create): failed to resume npc script character=%d: %v", ch.GetID(), err)
	}
}

func (h *GuildOperation) inviterCanInvite(characterID, guildID uint32) bool {
	if characterID == 0 || guildID == 0 {
		return false
	}
	g := h.gs.GetGuildByID(guildID)
	if g == nil {
		return false
	}
	for _, m := range g.GetMembers() {
		if m == nil || m.GetCharacterId() != characterID {
			continue
		}
		rank := m.GetRank()
		if rank == internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER {
			return true
		}
		if rank == internal.GuildMemberRank_GUILD_MEMBER_RANK_JUNIOR {
			return true
		}
		return false
	}
	return false
}

func (h *GuildOperation) isGuildMaster(characterID, guildID uint32) bool {
	if characterID == 0 || guildID == 0 {
		return false
	}
	g := h.gs.GetGuildByID(guildID)
	if g == nil {
		return false
	}
	for _, m := range g.GetMembers() {
		if m == nil || m.GetCharacterId() != characterID {
			continue
		}
		return m.GetRank() == internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER
	}
	return false
}

func (h *GuildOperation) isOnGuildEmblemMap(ch *entity.Character) bool {
	if ch == nil {
		return false
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil || mapInstance.Wz == nil {
		return false
	}
	return uint32(mapInstance.Wz.ID) == gameconst.GuildEmblemMapID
}

type guildEmblemChangePayment struct {
	usedCashItem bool
	mesoSpent    int32
}

func (h *GuildOperation) canPayGuildCreateCost(ch *entity.Character) bool {
	if ch == nil {
		return false
	}
	return ch.Meso >= gameconst.GuildCreateMesoCost
}

func (h *GuildOperation) chargeGuildCreateCost(ch *entity.Character) {
	if ch == nil {
		return
	}
	ch.RemoveMeso(gameconst.GuildCreateMesoCost)
}

func (h *GuildOperation) canPayGuildEmblemChangeCost(ch *entity.Character) bool {
	if ch == nil {
		return false
	}
	if ch.HasItem(gameconst.GuildEmblemChangeCashItemID) {
		return true
	}
	return ch.Meso >= gameconst.GuildEmblemChangeMesoCost
}

func (h *GuildOperation) chargeGuildEmblemChangeCost(ch *entity.Character) (guildEmblemChangePayment, bool) {
	if ch == nil {
		return guildEmblemChangePayment{}, false
	}
	if ch.HasItem(gameconst.GuildEmblemChangeCashItemID) {
		if !ch.RemoveItemByID(gameconst.GuildEmblemChangeCashItemID) {
			return guildEmblemChangePayment{}, false
		}
		return guildEmblemChangePayment{usedCashItem: true}, true
	}
	if ch.Meso < gameconst.GuildEmblemChangeMesoCost {
		return guildEmblemChangePayment{}, false
	}
	ch.RemoveMeso(gameconst.GuildEmblemChangeMesoCost)
	return guildEmblemChangePayment{mesoSpent: gameconst.GuildEmblemChangeMesoCost}, true
}

func (h *GuildOperation) refundGuildEmblemChangeCost(ch *entity.Character, payment guildEmblemChangePayment) {
	if ch == nil {
		return
	}
	if payment.usedCashItem {
		item, err := entity.NewItem(gameconst.GuildEmblemChangeCashItemID, 1, h.gs)
		if err != nil {
			log.Printf("GuildOperation(emblem): refund cash item failed character=%d: %v", ch.GetID(), err)
			return
		}
		if _, err := ch.AddItem(item, false); err != nil {
			log.Printf("GuildOperation(emblem): refund cash item failed character=%d: %v", ch.GetID(), err)
		}
		return
	}
	if payment.mesoSpent > 0 {
		ch.AddMeso(payment.mesoSpent)
	}
}

func (h *GuildOperation) applyLocalGuildMembership(ch *entity.Character, guildPb *internal.Guild) {
	if h == nil || h.gs == nil || ch == nil || guildPb == nil || h.gs.guild == nil {
		return
	}
	guildID := guildPb.GetGuildId()
	if guildID == 0 {
		return
	}
	id := guildID
	ch.SetGuildID(&id)
	h.gs.guild.Update(guildPb)
	if ch.Listener == nil {
		return
	}
	ch.Listener.OnShowGuildInfo(ch)
	ch.Listener.OnBroadcastGuildAppearance(ch)
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
		if _, inGuild := ch.GetGuildID(); inGuild {
			_ = ch.Send(&response.GuildMessage{
				Code: pconst.GuildResponseAlreadyInGuild,
			}, types.SEND_POLICY_ENCRYPT)
			h.resumeGuildCreate(ch, gameconst.GuildCreateResultAlreadyInGuild)
			return nil
		}
		if !h.isOnGuildEmblemMap(ch) {
			h.resumeGuildCreate(ch, gameconst.GuildCreateResultNotAllowed)
			return nil
		}
		if !h.canPayGuildCreateCost(ch) {
			h.resumeGuildCreate(ch, gameconst.GuildCreateResultInsufficientMeso)
			return nil
		}

		promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise = async.ThenRPC(promise, func(c context.Context) (*internal.CreateGuildReply, error) {
			leader := ch.ToProtoGuildMember(worldID, int32(h.gs.config.ChannelId), internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER)
			return h.gs.internalClient.CreateGuild(c, &internal.CreateGuildRequest{
				WorldId:   worldID,
				GuildName: req.GuildName,
				Leader:    leader,
			})
		}, func(reply *internal.CreateGuildReply) error {
			if !reply.GetOk() {
				if reply.GetErrorCode() == internal.GuildErrorCode_GUILD_ERROR_ALREADY_IN_GUILD {
					_ = ch.Send(&response.GuildMessage{
						Code: pconst.GuildResponseAlreadyInGuild,
					}, types.SEND_POLICY_ENCRYPT)
				} else {
					log.Printf("GuildOperation(create): failed character=%d code=%v", charID, reply.GetErrorCode())
				}
				h.resumeGuildCreate(ch, gameconst.GuildCreateResultFailed)
				return nil
			}
			guildPb := reply.GetGuild()
			if guildPb == nil || guildPb.GetGuildId() == 0 {
				log.Printf("GuildOperation(create): ok but missing guild character=%d", charID)
				h.resumeGuildCreate(ch, gameconst.GuildCreateResultFailed)
				return nil
			}
			if h.gs.guild != nil {
				h.applyLocalGuildMembership(ch, guildPb)
			}
			h.chargeGuildCreateCost(ch)
			h.resumeGuildCreate(ch, gameconst.GuildCreateResultOK)
			return nil
		})
		promise.OnError(func(err error) {
			log.Printf("GuildOperation(create) async error: %v", err)
			h.resumeGuildCreate(ch, gameconst.GuildCreateResultFailed)
		}).Run()
		return nil
	case pconst.GuildC2SInvite:
		guildID, inGuild := ch.GetGuildID()
		if !inGuild || !h.inviterCanInvite(charID, guildID) {
			return nil
		}
		targetName := req.TargetName
		if targetName == "" {
			return nil
		}
		if h.gs.characterRuntime == nil {
			return nil
		}
		targetID, ok := h.gs.characterRuntime.GetCharacterIDByName(targetName)
		if !ok || targetID == 0 {
			_ = ch.Send(&response.GuildMessage{
				Code: pconst.GuildResponseNotInChannel,
			}, types.SEND_POLICY_ENCRYPT)
			return nil
		}
		if mapInstance := ch.GetMap(); mapInstance != nil {
			if target := mapInstance.GetPlayer(targetID); target != nil {
				if _, targetInGuild := target.GetGuildID(); targetInGuild {
					_ = ch.Send(&response.GuildMessage{
						Code: pconst.GuildResponseAlreadyInGuild,
					}, types.SEND_POLICY_ENCRYPT)
					return nil
				}
				now := time.Now()
				for id, expiresAt := range target.GuildInvites {
					if !now.Before(expiresAt) {
						delete(target.GuildInvites, id)
					}
				}
				if len(target.GuildInvites) > 0 {
					ch.Listener.OnMessage(ch, gameconst.MsgPinkText, gameconst.GuildInviteTargetBusyMessage)
					return nil
				}
				target.GuildInvites[guildID] = now.Add(gameconst.GuildInviteDuration)
				target.Listener.OnGuildInvite(target, guildID, ch.GetName())
				return nil
			}
		}
		h.gs.EnsureSend(nil, targetID, &g_actor.DeliverGuildInvite{
			CharacterID:        targetID,
			InviterCharacterID: charID,
			GuildID:            guildID,
			InviterName:        ch.GetName(),
		})
		return nil
	case pconst.GuildC2SAcceptInvite:
		if _, inGuild := ch.GetGuildID(); inGuild {
			return nil
		}
		if req.GuildID == 0 || req.CharacterID != charID {
			return nil
		}
		now := time.Now()
		expiresAt, hasInvite := ch.GuildInvites[req.GuildID]
		for id, exp := range ch.GuildInvites {
			if !now.Before(exp) {
				delete(ch.GuildInvites, id)
			}
		}
		if !hasInvite || !now.Before(expiresAt) {
			return nil
		}
		delete(ch.GuildInvites, req.GuildID)

		promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise = async.ThenRPC(promise, func(c context.Context) (*internal.AcceptGuildInviteReply, error) {
			member := ch.ToProtoGuildMember(worldID, int32(h.gs.config.ChannelId), internal.GuildMemberRank_GUILD_MEMBER_RANK_NEW)
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
			if guildPb == nil || guildPb.GetGuildId() == 0 {
				log.Printf("GuildOperation(accept invite): ok but missing guild character=%d", charID)
				return nil
			}
			if h.gs.guild != nil {
				h.applyLocalGuildMembership(ch, guildPb)
			}
			return nil
		})
		promise.OnError(func(err error) {
			log.Printf("GuildOperation(accept invite) async error: %v", err)
		}).Run()
		return nil
	case pconst.GuildC2SLeave:
		if req.CharacterID != charID {
			return nil
		}
		if req.CharacterName != ch.GetName() {
			return nil
		}
		if _, inGuild := ch.GetGuildID(); !inGuild {
			return nil
		}
		promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise = async.ThenRPC(promise, func(c context.Context) (*internal.LeaveGuildReply, error) {
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
		}).Run()
		return nil
	case pconst.GuildC2SExpel:
		targetID := req.CharacterID
		if targetID == 0 {
			return nil
		}
		guildID, inGuild := ch.GetGuildID()
		if !inGuild || !h.inviterCanInvite(charID, guildID) {
			return nil
		}
		promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise = async.ThenRPC(promise, func(c context.Context) (*internal.ExpelGuildReply, error) {
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
		}).Run()
		return nil
	case pconst.GuildC2SChangeRankTitles:
		guildID, inGuild := ch.GetGuildID()
		if !inGuild || !h.isGuildMaster(charID, guildID) {
			return nil
		}
		promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise = async.ThenRPC(promise, func(c context.Context) (*internal.ChangeGuildRankTitlesReply, error) {
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
		}).Run()
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
		guildID, inGuild := ch.GetGuildID()
		if !inGuild || !h.inviterCanInvite(charID, guildID) {
			return nil
		}
		if newRank <= 2 && !h.isGuildMaster(charID, guildID) {
			return nil
		}
		protoRank := internal.GuildMemberRank(newRank)
		promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise = async.ThenRPC(promise, func(c context.Context) (*internal.ChangeGuildMemberRankReply, error) {
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
		}).Run()
		return nil
	case pconst.GuildC2SChangeEmblem:
		guildID, inGuild := ch.GetGuildID()
		if !inGuild || !h.isGuildMaster(charID, guildID) || !h.isOnGuildEmblemMap(ch) {
			return nil
		}
		if !h.canPayGuildEmblemChangeCost(ch) {
			ch.Listener.OnMessage(ch, gameconst.MsgPopup, gameconst.GuildEmblemChangeInsufficientCostMessage)
			return nil
		}
		payment, paid := h.chargeGuildEmblemChangeCost(ch)
		if !paid {
			ch.Listener.OnMessage(ch, gameconst.MsgPopup, gameconst.GuildEmblemChangeInsufficientCostMessage)
			return nil
		}
		promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise = async.ThenRPC(promise, func(c context.Context) (*internal.ChangeGuildEmblemReply, error) {
			return h.gs.internalClient.ChangeGuildEmblem(c, &internal.ChangeGuildEmblemRequest{
				WorldId:     worldID,
				CharacterId: charID,
				Logo: &internal.GuildLogo{
					Logo:        uint32(req.Logo),
					LogoColor:   uint32(req.LogoColor),
					LogoBg:      uint32(req.LogoBG),
					LogoBgColor: uint32(req.LogoBGColor),
				},
			})
		}, func(reply *internal.ChangeGuildEmblemReply) error {
			if !reply.GetOk() {
				h.refundGuildEmblemChangeCost(ch, payment)
				if reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD &&
					reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_AUTHORIZED {
					log.Printf("GuildOperation(change emblem): failed character=%d code=%v", charID, reply.GetErrorCode())
				}
				return nil
			}
			return nil
		})
		promise.OnError(func(err error) {
			h.refundGuildEmblemChangeCost(ch, payment)
			log.Printf("GuildOperation(change emblem) async error: %v", err)
		}).Run()
		return nil
	case pconst.GuildC2SChangeNotice:
		guildID, inGuild := ch.GetGuildID()
		if !inGuild || !h.inviterCanInvite(charID, guildID) {
			return nil
		}
		promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		promise = async.ThenRPC(promise, func(c context.Context) (*internal.ChangeGuildNoticeReply, error) {
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
		}).Run()
		return nil
	default:
		return nil
	}
}
