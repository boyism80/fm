package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
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
	if req.GuildID == 0 || req.CharacterID != ch.GetID() {
		return nil
	}
	if ch.Guild().AcceptInvite(req.GuildID) != nil {
		return nil
	}
	h.gs.systems.Get[*game.GuildSystem]().AcceptInviteAsync(ctx.ActorContext, ch, req.GuildID)
	return nil
}

func (h *GuildOperation) leave(ctx *core.ClientContext, ch *game.Character, req *request.GuildOperation) error {
	if req.CharacterID != ch.GetID() || req.CharacterName != ch.GetName() {
		return nil
	}
	if _, inGuild := ch.Guild().ID(); !inGuild {
		return nil
	}
	h.gs.systems.Get[*game.GuildSystem]().LeaveAsync(ctx.ActorContext, ch)
	return nil
}

// The changes below leave who may make them to internal, which checks the latest rows under its lock.
func (h *GuildOperation) expel(ctx *core.ClientContext, ch *game.Character, req *request.GuildOperation) error {
	if _, inGuild := ch.Guild().ID(); !inGuild || req.CharacterID == 0 {
		return nil
	}
	h.gs.systems.Get[*game.GuildSystem]().ExpelAsync(ctx.ActorContext, ch, req.CharacterID)
	return nil
}

func (h *GuildOperation) changeRankTitles(ctx *core.ClientContext, ch *game.Character, req *request.GuildOperation) error {
	if _, inGuild := ch.Guild().ID(); !inGuild {
		return nil
	}
	h.gs.systems.Get[*game.GuildSystem]().ChangeRankTitlesAsync(ctx.ActorContext, ch, req.RankTitles[:])
	return nil
}

func (h *GuildOperation) changeMemberRank(ctx *core.ClientContext, ch *game.Character, req *request.GuildOperation) error {
	if _, inGuild := ch.Guild().ID(); !inGuild || req.CharacterID == 0 {
		return nil
	}
	h.gs.systems.Get[*game.GuildSystem]().ChangeMemberRankAsync(ctx.ActorContext, ch, req.CharacterID, internal.GuildMemberRank(req.NewMemberRank))
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
	if _, inGuild := ch.Guild().ID(); !inGuild {
		return nil
	}
	h.gs.systems.Get[*game.GuildSystem]().ChangeNoticeAsync(ctx.ActorContext, ch, req.Notice)
	return nil
}
