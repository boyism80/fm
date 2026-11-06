package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/rpc"
	gameconst "github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type GuildOperation struct {
	gs *GameServer
}

func (GuildOperation) New(gs *GameServer) *GuildOperation {
	return &GuildOperation{
		gs: gs,
	}
}

func (h *GuildOperation) chargeGuildEmblemChangeCost(ch *game.Character) (game.GuildEmblemDetail, bool) {
	if ch.Inventory().HasItem(gameconst.GuildEmblemChangeCashItemID) {
		if !ch.Inventory().RemoveByItemIDCount(gameconst.GuildEmblemChangeCashItemID, 1) {
			return game.GuildEmblemDetail{}, false
		}
		return game.GuildEmblemDetail{CashItem: true}, true
	}
	if !ch.Inventory().RemoveMeso(gameconst.GuildEmblemChangeMesoCost) {
		return game.GuildEmblemDetail{}, false
	}
	return game.GuildEmblemDetail{Meso: gameconst.GuildEmblemChangeMesoCost}, true
}

func (h *GuildOperation) Handle(ctx *core.ClientContext, req *request.GuildOperation) error {
	if ctx.ActorContext == nil {
		return fmt.Errorf("guild operation: actor context required")
	}

	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}

	switch req.Operation {
	case pconst.GuildC2SCreate:
		return h.create(ctx, ch, req)
	case pconst.GuildC2SInvite:
		return h.invite(ctx, ch, req)
	case pconst.GuildC2SAcceptInvite:
		return h.acceptInvite(ctx, ch, req)
	case pconst.GuildC2SLeave:
		return h.leave(ctx, ch, req)
	case pconst.GuildC2SExpel:
		return h.expel(ctx, ch, req)
	case pconst.GuildC2SChangeRankTitles:
		return h.changeRankTitles(ctx, ch, req)
	case pconst.GuildC2SChangeMemberRank:
		return h.changeMemberRank(ctx, ch, req)
	case pconst.GuildC2SChangeEmblem:
		return h.changeEmblem(ctx, ch, req)
	case pconst.GuildC2SChangeNotice:
		return h.changeNotice(ctx, ch, req)
	default:
		return nil
	}
}

func (h *GuildOperation) create(ctx *core.ClientContext, ch *game.Character, req *request.GuildOperation) error {
	charID := ch.GetID()

	if _, inGuild := ch.Guild().ID(); inGuild {
		ch.Listener().OnGuildMessage(ch, pconst.GuildResponseAlreadyInGuild)
		return nil
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil || mapInstance.Wz() == nil || uint32(mapInstance.Wz().ID) != gameconst.GuildEmblemMapID {
		ch.Listener().OnMessage(ch, gameconst.MsgPopup, gameconst.GuildCreateNotAllowedMessage)
		return nil
	}
	if !ch.Inventory().RemoveMeso(gameconst.GuildCreateMesoCost) {
		ch.Listener().OnMessage(ch, gameconst.MsgPopup, gameconst.GuildCreateInsufficientMesoMessage)
		return nil
	}

	detail := game.GuildCreateDetail{Name: req.GuildName, Meso: gameconst.GuildCreateMesoCost}
	id := ch.Unsettled().Add(internal.UnsettledType_UNSETTLED_TYPE_GUILD_CREATE, 0, detail)
	h.gs.systems.Get[*game.GuildSystem]().CreateAsync(ctx.ActorContext, ch, id, detail.Name).Bind(ch).Do(func(reply *internal.CreateGuildReply) error {
		ch.Unsettled().Erase(id)
		if !reply.GetOk() {
			switch reply.GetErrorCode() {
			case internal.GuildErrorCode_GUILD_ERROR_ALREADY_IN_GUILD:
				ch.Listener().OnGuildMessage(ch, pconst.GuildResponseAlreadyInGuild)
			case internal.GuildErrorCode_GUILD_ERROR_GUILD_NAME_INVALID:
				ch.Listener().OnMessage(ch, gameconst.MsgPopup, gameconst.GuildNameInvalidMessage)
			default:
				log.Printf("GuildOperation(create): failed character=%d code=%v", charID, reply.GetErrorCode())
				ch.Listener().OnMessage(ch, gameconst.MsgPopup, gameconst.GuildCreateFailedMessage)
			}
			ch.Inventory().AddMeso(detail.Meso)
			return nil
		}
		guildPb := reply.GetGuild()
		if guildPb == nil {
			log.Printf("GuildOperation(create): ok but missing guild character=%d", charID)
			ch.Listener().OnMessage(ch, gameconst.MsgPopup, gameconst.GuildCreateFailedMessage)
			return nil
		}
		h.gs.systems.Get[*game.GuildSystem]().StoreGuild(ch, guildPb)
		ch.Listener().OnMessage(ch, gameconst.MsgPopup, gameconst.GuildCreatedMessage)
		return nil
	}).OnError(func(err error) {
		ch.Unsettled().Lost(id)
		ch.Listener().OnMessage(ch, gameconst.MsgPopup, gameconst.GuildCreateFailedMessage)

		log.Printf("GuildOperation(create) character=%d id=%d: %v", charID, id, err)
	})
	return nil
}

func (h *GuildOperation) invite(ctx *core.ClientContext, ch *game.Character, req *request.GuildOperation) error {
	charID := ch.GetID()

	guildID, inGuild := ch.Guild().ID()
	g := h.gs.systems.Get[*game.GuildSystem]().Get(guildID)
	if !inGuild || g == nil || !g.CanInvite(charID) {
		return nil
	}
	targetName := req.TargetName
	if targetName == "" {
		return nil
	}
	targetID, ok := h.gs.systems.Get[*game.DispatchSystem]().FindCharacterID(targetName)
	if !ok || targetID == 0 {
		ch.Listener().OnGuildMessage(ch, pconst.GuildResponseNotInChannel)
		return nil
	}
	if mapInstance := ch.GetMap(); mapInstance != nil {
		if target := mapInstance.GetPlayer(targetID); target != nil {
			switch target.Guild().Invite(guildID, ch.GetName()) {
			case game.ErrGuildAlreadyJoined:
				ch.Listener().OnGuildMessage(ch, pconst.GuildResponseAlreadyInGuild)
			case game.ErrGuildInviteBusy:
				ch.Listener().OnMessage(ch, gameconst.MsgPinkText, gameconst.GuildInviteTargetBusyMessage)
			}
			return nil
		}
	}
	h.gs.systems.Get[*game.DispatchSystem]().SendToCharacter(nil, targetID, &game.InviteToGuild{
		CharacterID:        targetID,
		InviterCharacterID: charID,
		GuildID:            guildID,
		InviterName:        ch.GetName(),
	})
	return nil
}

func (h *GuildOperation) acceptInvite(ctx *core.ClientContext, ch *game.Character, req *request.GuildOperation) error {
	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	if req.GuildID == 0 || req.CharacterID != charID {
		return nil
	}
	if ch.Guild().AcceptInvite(req.GuildID) != nil {
		return nil
	}

	member := game.GuildMemberFromCharacter(ch, worldID, int32(h.gs.config.ChannelId), internal.GuildMemberRank_GUILD_MEMBER_RANK_NEW).BuildRPC()
	promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).Bind(ch)
	promise.ThenRPC(func(c context.Context) (*internal.AcceptGuildInviteReply, error) {
		return h.gs.systems.Get[*rpc.InternalConnection]().Client.AcceptGuildInvite(c, &internal.AcceptGuildInviteRequest{
			WorldId: worldID,
			GuildId: req.GuildID,
			Member:  member,
		})
	}, func(reply *internal.AcceptGuildInviteReply) error {
		if !reply.GetOk() {
			if reply.GetErrorCode() == internal.GuildErrorCode_GUILD_ERROR_GUILD_FULL {
				ch.Listener().OnMessage(ch, gameconst.MsgPopup, gameconst.GuildFullMessage)
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
		h.gs.systems.Get[*game.GuildSystem]().StoreGuild(ch, guildPb)
		if alliancePb := reply.GetAlliance(); alliancePb != nil {
			h.gs.systems.Get[*game.AllianceSystem]().Update(alliancePb)
			ch.Listener().OnShowAllianceInfo(ch)
		}
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("GuildOperation(accept invite) async error: %v", err)
	})
	return nil
}

func (h *GuildOperation) leave(ctx *core.ClientContext, ch *game.Character, req *request.GuildOperation) error {
	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	if req.CharacterID != charID {
		return nil
	}
	if req.CharacterName != ch.GetName() {
		return nil
	}
	if _, inGuild := ch.Guild().ID(); !inGuild {
		return nil
	}
	promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise.ThenRPC(func(c context.Context) (*internal.LeaveGuildReply, error) {
		return h.gs.systems.Get[*rpc.InternalConnection]().Client.LeaveGuild(c, &internal.LeaveGuildRequest{
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
}

func (h *GuildOperation) expel(ctx *core.ClientContext, ch *game.Character, req *request.GuildOperation) error {
	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	targetID := req.CharacterID
	if targetID == 0 {
		return nil
	}
	guildID, inGuild := ch.Guild().ID()
	g := h.gs.systems.Get[*game.GuildSystem]().Get(guildID)
	if !inGuild || g == nil || !g.CanInvite(charID) {
		return nil
	}
	promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise.ThenRPC(func(c context.Context) (*internal.ExpelGuildReply, error) {
		return h.gs.systems.Get[*rpc.InternalConnection]().Client.ExpelGuild(c, &internal.ExpelGuildRequest{
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
}

func (h *GuildOperation) changeRankTitles(ctx *core.ClientContext, ch *game.Character, req *request.GuildOperation) error {
	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	guildID, inGuild := ch.Guild().ID()
	g := h.gs.systems.Get[*game.GuildSystem]().Get(guildID)
	if !inGuild || g == nil || !g.IsGuildMaster(charID) || !game.RankTitles(req.RankTitles).Valid() {
		return nil
	}
	promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise.ThenRPC(func(c context.Context) (*internal.ChangeGuildRankTitlesReply, error) {
		return h.gs.systems.Get[*rpc.InternalConnection]().Client.ChangeGuildRankTitles(c, &internal.ChangeGuildRankTitlesRequest{
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
}

func (h *GuildOperation) changeMemberRank(ctx *core.ClientContext, ch *game.Character, req *request.GuildOperation) error {
	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	newRank := req.NewMemberRank
	if newRank <= 1 || newRank > 5 {
		return nil
	}
	targetID := req.CharacterID
	if targetID == 0 || targetID == charID {
		return nil
	}
	guildID, inGuild := ch.Guild().ID()
	g := h.gs.systems.Get[*game.GuildSystem]().Get(guildID)
	if !inGuild || g == nil || !g.CanInvite(charID) {
		return nil
	}
	if newRank <= 2 && !g.IsGuildMaster(charID) {
		return nil
	}
	protoRank := internal.GuildMemberRank(newRank)
	promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise.ThenRPC(func(c context.Context) (*internal.ChangeGuildMemberRankReply, error) {
		return h.gs.systems.Get[*rpc.InternalConnection]().Client.ChangeGuildMemberRank(c, &internal.ChangeGuildMemberRankRequest{
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
}

func (h *GuildOperation) changeEmblem(ctx *core.ClientContext, ch *game.Character, req *request.GuildOperation) error {
	charID := ch.GetID()

	guildID, inGuild := ch.Guild().ID()
	g := h.gs.systems.Get[*game.GuildSystem]().Get(guildID)
	mapInstance := ch.GetMap()
	if !inGuild || g == nil || !g.IsGuildMaster(charID) ||
		mapInstance == nil || mapInstance.Wz() == nil || uint32(mapInstance.Wz().ID) != gameconst.GuildEmblemMapID {
		return nil
	}
	detail, paid := h.chargeGuildEmblemChangeCost(ch)
	if !paid {
		ch.Listener().OnMessage(ch, gameconst.MsgPopup, gameconst.GuildEmblemChangeInsufficientCostMessage)
		return nil
	}
	detail.Logo = game.NewGuildLogo(uint32(req.Logo), uint32(req.LogoColor), uint32(req.LogoBG), uint32(req.LogoBGColor))
	id := ch.Unsettled().Add(internal.UnsettledType_UNSETTLED_TYPE_GUILD_EMBLEM, uint64(guildID), detail)
	h.gs.systems.Get[*game.GuildSystem]().ChangeEmblemAsync(ctx.ActorContext, ch, id, guildID, detail.Logo.BuildRPC()).Bind(ch).Do(func(reply *internal.ChangeGuildEmblemReply) error {
		ch.Unsettled().Erase(id)
		if !reply.GetOk() {
			h.gs.systems.Get[*game.GuildSystem]().RefundEmblem(ch, detail)
			if reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD &&
				reply.GetErrorCode() != internal.GuildErrorCode_GUILD_ERROR_NOT_AUTHORIZED {
				log.Printf("GuildOperation(change emblem): failed character=%d code=%v", charID, reply.GetErrorCode())
			}
		}
		return nil
	}).OnError(func(err error) {
		ch.Unsettled().Lost(id)

		log.Printf("GuildOperation(change emblem) character=%d id=%d: %v", charID, id, err)
	})
	return nil
}

func (h *GuildOperation) changeNotice(ctx *core.ClientContext, ch *game.Character, req *request.GuildOperation) error {
	worldID := h.gs.config.WorldId
	charID := ch.GetID()

	guildID, inGuild := ch.Guild().ID()
	g := h.gs.systems.Get[*game.GuildSystem]().Get(guildID)
	if !inGuild || g == nil || !g.CanInvite(charID) {
		return nil
	}
	promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise.ThenRPC(func(c context.Context) (*internal.ChangeGuildNoticeReply, error) {
		return h.gs.systems.Get[*rpc.InternalConnection]().Client.ChangeGuildNotice(c, &internal.ChangeGuildNoticeRequest{
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
}
