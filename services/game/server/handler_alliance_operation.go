package server

import (
	"context"
	"fmt"
	"github.com/boyism80/fm/core/clock"
	"log"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/client"
	gameconst "github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type AllianceOperation struct {
	gs *GameServer
}

func (AllianceOperation) New(gs *GameServer) *AllianceOperation {
	return &AllianceOperation{
		gs: gs,
	}
}

func (h *AllianceOperation) Handle(ctx *core.ClientContext, req *request.AllianceOperation) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}
	ch := gameClient.GetCharacter()
	if ch == nil {
		return fmt.Errorf("character not found")
	}

	switch req.Operation {
	case pconst.AllianceC2SLoadInfo:
		return h.handleLoadInfo(ch)
	case pconst.AllianceC2SCreate:
		return h.handleCreate(ctx, ch, req)
	case pconst.AllianceC2SLeave:
		return h.handleLeave(ctx, ch)
	case pconst.AllianceC2SInvite:
		return h.handleInvite(ctx, ch, req)
	case pconst.AllianceC2SAcceptInvite:
		return h.handleAcceptInvite(ctx, ch)
	case pconst.AllianceC2SDenyInvite:
		return h.handleDenyInvite(ctx, ch)
	case pconst.AllianceC2SChangeRankTitles:
		return h.handleChangeRankTitles(ctx, ch, req)
	case pconst.AllianceC2SChangeLeader:
		return h.handleChangeLeader(ctx, ch, req)
	case pconst.AllianceC2SChangeNotice:
		return h.handleChangeNotice(ctx, ch, req)
	case pconst.AllianceC2SChangeMemberRank:
		return h.handleChangeMemberRank(ctx, ch, req)
	case pconst.AllianceC2SExpel:
		return h.handleExpel(ctx, ch, req)
	default:
		return nil
	}
}

func (h *AllianceOperation) handleLoadInfo(ch *entity.Character) error {
	if ch == nil || ch.Listener == nil {
		return nil
	}
	ch.Listener.OnAllianceUpdateInfo(ch)
	return nil
}

func (h *AllianceOperation) handleExpel(ctx *core.ClientContext, ch *entity.Character, req *request.AllianceOperation) error {
	if h.gs == nil || h.gs.internalClient == nil {
		return nil
	}
	targetGuildID := req.TargetGuildID
	if targetGuildID == 0 {
		return nil
	}
	guildID, inGuild := ch.GetGuildID()
	if !inGuild {
		return nil
	}
	if targetGuildID == guildID {
		return nil
	}
	g := h.gs.guild.Get(guildID)
	if g == nil {
		return nil
	}
	charID := ch.GetID()
	if g.LeaderCharacterID != charID {
		return nil
	}
	m := g.FindMember(charID)
	if m == nil || m.GetRank() != internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER {
		return nil
	}
	ar, hasRank := m.GetAllianceRank()
	if !hasRank || ar != 1 {
		return nil
	}
	allianceID, inAlliance := g.GetAllianceID()
	if !inAlliance {
		return nil
	}
	if req.AllianceID != 0 && req.AllianceID != allianceID {
		return nil
	}
	worldID := h.gs.config.WorldId
	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.ExpelAllianceGuildReply, error) {
		return h.gs.internalClient.ExpelAllianceGuild(c, &internal.ExpelAllianceGuildRequest{
			WorldId:       worldID,
			CharacterId:   charID,
			TargetGuildId: targetGuildID,
			AllianceId:    allianceID,
		})
	}, func(reply *internal.ExpelAllianceGuildReply) error {
		if reply == nil || !reply.GetOk() {
			if reply != nil && !reply.GetOk() {
				log.Printf("AllianceOperation(expel): failed character=%d target_guild=%d code=%v",
					charID, targetGuildID, reply.GetErrorCode())
			}
			return nil
		}
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("AllianceOperation(expel) async error: %v", err)
	})
	return nil
}

func (h *AllianceOperation) handleChangeLeader(ctx *core.ClientContext, ch *entity.Character, req *request.AllianceOperation) error {
	if h.gs == nil || h.gs.internalClient == nil {
		return nil
	}
	newLeaderID := req.NewLeaderID
	if newLeaderID == 0 {
		return nil
	}
	guildID, inGuild := ch.GetGuildID()
	if !inGuild {
		return nil
	}
	g := h.gs.guild.Get(guildID)
	if g == nil {
		return nil
	}
	charID := ch.GetID()
	if g.LeaderCharacterID != charID {
		return nil
	}
	m := g.FindMember(charID)
	if m == nil || m.GetRank() != internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER {
		return nil
	}
	ar, hasRank := m.GetAllianceRank()
	if !hasRank || ar != 1 {
		return nil
	}
	allianceID, inAlliance := g.GetAllianceID()
	if !inAlliance {
		return nil
	}
	a := h.cachedAlliance(allianceID)
	if a == nil || a.LeaderCharacterID != charID {
		return nil
	}
	worldID := h.gs.config.WorldId
	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.ChangeAllianceLeaderReply, error) {
		return h.gs.internalClient.ChangeAllianceLeader(c, &internal.ChangeAllianceLeaderRequest{
			WorldId:              worldID,
			CharacterId:          charID,
			NewLeaderCharacterId: newLeaderID,
		})
	}, func(reply *internal.ChangeAllianceLeaderReply) error {
		if reply == nil || !reply.GetOk() {
			if reply != nil && !reply.GetOk() {
				log.Printf("AllianceOperation(change leader): failed character=%d new_leader=%d code=%v",
					charID, newLeaderID, reply.GetErrorCode())
			}
			return nil
		}
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("AllianceOperation(change leader) async error: %v", err)
	})
	return nil
}

func (h *AllianceOperation) handleChangeNotice(ctx *core.ClientContext, ch *entity.Character, req *request.AllianceOperation) error {
	if h.gs == nil || h.gs.internalClient == nil {
		return nil
	}
	guildID, inGuild := ch.GetGuildID()
	if !inGuild {
		return nil
	}
	g := h.gs.guild.Get(guildID)
	if g == nil {
		return nil
	}
	charID := ch.GetID()
	if g.LeaderCharacterID != charID {
		return nil
	}
	m := g.FindMember(charID)
	if m == nil || m.GetRank() != internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER {
		return nil
	}
	ar, hasRank := m.GetAllianceRank()
	if !hasRank || ar > 2 {
		return nil
	}
	if _, inAlliance := g.GetAllianceID(); !inAlliance {
		return nil
	}
	worldID := h.gs.config.WorldId
	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.ChangeAllianceNoticeReply, error) {
		return h.gs.internalClient.ChangeAllianceNotice(c, &internal.ChangeAllianceNoticeRequest{
			WorldId:     worldID,
			CharacterId: charID,
			Notice:      req.Notice,
		})
	}, func(reply *internal.ChangeAllianceNoticeReply) error {
		if reply == nil || !reply.GetOk() {
			if reply != nil && !reply.GetOk() {
				log.Printf("AllianceOperation(change notice): failed character=%d code=%v", charID, reply.GetErrorCode())
			}
			return nil
		}
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("AllianceOperation(change notice) async error: %v", err)
	})
	return nil
}

func (h *AllianceOperation) handleChangeRankTitles(ctx *core.ClientContext, ch *entity.Character, req *request.AllianceOperation) error {
	if h.gs == nil || h.gs.internalClient == nil {
		return nil
	}
	guildID, inGuild := ch.GetGuildID()
	if !inGuild {
		return nil
	}
	g := h.gs.guild.Get(guildID)
	if g == nil {
		return nil
	}
	charID := ch.GetID()
	if g.LeaderCharacterID != charID {
		return nil
	}
	m := g.FindMember(charID)
	if m == nil || m.GetRank() != internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER {
		return nil
	}
	ar, hasRank := m.GetAllianceRank()
	if !hasRank || ar != 1 {
		return nil
	}
	allianceID, inAlliance := g.GetAllianceID()
	if !inAlliance {
		return nil
	}
	a := h.cachedAlliance(allianceID)
	if a == nil || a.LeaderCharacterID != charID {
		return nil
	}
	worldID := h.gs.config.WorldId
	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.ChangeAllianceRankTitlesReply, error) {
		return h.gs.internalClient.ChangeAllianceRankTitles(c, &internal.ChangeAllianceRankTitlesRequest{
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
	}, func(reply *internal.ChangeAllianceRankTitlesReply) error {
		if reply == nil || !reply.GetOk() {
			if reply != nil && !reply.GetOk() {
				log.Printf("AllianceOperation(change rank titles): failed character=%d code=%v", charID, reply.GetErrorCode())
			}
			return nil
		}
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("AllianceOperation(change rank titles) async error: %v", err)
	})
	return nil
}

func (h *AllianceOperation) handleChangeMemberRank(ctx *core.ClientContext, ch *entity.Character, req *request.AllianceOperation) error {
	if h.gs == nil || h.gs.internalClient == nil {
		return nil
	}
	targetID := req.TargetCharacterID
	if targetID == 0 {
		return nil
	}
	guildID, inGuild := ch.GetGuildID()
	if !inGuild {
		return nil
	}
	g := h.gs.guild.Get(guildID)
	if g == nil {
		return nil
	}
	charID := ch.GetID()
	if g.LeaderCharacterID != charID {
		return nil
	}
	m := g.FindMember(charID)
	if m == nil || m.GetRank() != internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER {
		return nil
	}
	ar, hasRank := m.GetAllianceRank()
	if !hasRank || ar > 2 {
		return nil
	}
	if _, inAlliance := g.GetAllianceID(); !inAlliance {
		return nil
	}
	worldID := h.gs.config.WorldId
	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.ChangeAllianceMemberRankReply, error) {
		return h.gs.internalClient.ChangeAllianceMemberRank(c, &internal.ChangeAllianceMemberRankRequest{
			WorldId:              worldID,
			RequesterCharacterId: charID,
			TargetCharacterId:    targetID,
			Promote:              req.RankChangePromote,
		})
	}, func(reply *internal.ChangeAllianceMemberRankReply) error {
		if reply == nil || !reply.GetOk() {
			if reply != nil && !reply.GetOk() {
				log.Printf("AllianceOperation(change member rank): failed requester=%d target=%d code=%v",
					charID, targetID, reply.GetErrorCode())
			}
			return nil
		}
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("AllianceOperation(change member rank) async error: %v", err)
	})
	return nil
}

func (h *AllianceOperation) handleInvite(ctx *core.ClientContext, ch *entity.Character, req *request.AllianceOperation) error {
	if h.gs == nil {
		return nil
	}
	guildID, inGuild := ch.GetGuildID()
	if !inGuild {
		return nil
	}
	g := h.gs.guild.Get(guildID)
	if g == nil {
		return nil
	}
	charID := ch.GetID()
	if g.LeaderCharacterID != charID {
		return nil
	}
	m := g.FindMember(charID)
	if m == nil || m.GetRank() != internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER {
		return nil
	}
	ar, hasRank := m.GetAllianceRank()
	if !hasRank || ar != 1 {
		return nil
	}
	allianceID, inAlliance := g.GetAllianceID()
	if !inAlliance {
		return nil
	}
	a, ok := h.inviterMaySendAllianceInvite(allianceID, charID)
	if !ok {
		return nil
	}
	guildName := req.TargetGuildLeaderName
	if guildName == "" {
		return nil
	}
	targetGuildID, nameFound := h.gs.guild.GuildIDByName(guildName)
	if !nameFound {
		ch.Listener.OnMessage(ch, gameconst.MsgPopup, gameconst.AllianceInviteGuildNotFoundMessage)
		return nil
	}
	targetG := h.gs.guild.Get(targetGuildID)
	if targetG == nil {
		ch.Listener.OnMessage(ch, gameconst.MsgPopup, gameconst.AllianceInviteGuildNotFoundMessage)
		return nil
	}
	if _, targetInAlliance := targetG.GetAllianceID(); targetInAlliance || targetG.LeaderCharacterID == 0 {
		ch.Listener.OnMessage(ch, gameconst.MsgPopup, gameconst.AllianceInviteGuildNotFoundMessage)
		return nil
	}
	targetLeaderID := targetG.LeaderCharacterID
	if h.gs.characterRuntime == nil || !h.gs.characterRuntime.Exists(targetLeaderID) {
		ch.Listener.OnMessage(ch, gameconst.MsgPopup, gameconst.AllianceInviteTargetNotOnlineMessage)
		return nil
	}
	if !h.withStoredGuild(targetGuildID, func(g *entity.Guild) bool {
		if _, inAlliance := g.GetAllianceID(); inAlliance {
			return false
		}
		return !g.HasAllianceInvite()
	}) {
		ch.Listener.OnMessage(ch, gameconst.MsgPinkText, gameconst.GuildInviteTargetBusyMessage)
		return nil
	}
	h.gs.EnsureSend(nil, targetLeaderID, &g_actor.DeliverAllianceInvite{
		CharacterID:        targetLeaderID,
		TargetGuildID:      targetGuildID,
		AllianceID:         allianceID,
		InviterCharacterID: charID,
		InviterGuildID:     guildID,
		InviterName:        ch.GetName(),
		AllianceName:       a.Name,
		ExpiresAt:          clock.Now().Add(gameconst.GuildInviteDuration),
	})
	return nil
}

func (h *AllianceOperation) handleAcceptInvite(ctx *core.ClientContext, ch *entity.Character) error {
	if h.gs == nil || h.gs.internalClient == nil {
		return nil
	}
	guildID, inGuild := ch.GetGuildID()
	if !inGuild {
		return nil
	}
	g := h.gs.guild.Get(guildID)
	if g == nil {
		return nil
	}
	charID := ch.GetID()
	if g.LeaderCharacterID != charID {
		return nil
	}
	m := g.FindMember(charID)
	if m == nil || m.GetRank() != internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER {
		return nil
	}
	if _, inAlliance := g.GetAllianceID(); inAlliance {
		return nil
	}
	var allianceID uint32
	var expiresAt time.Time
	var hasInvite bool
	h.withStoredGuild(guildID, func(sg *entity.Guild) bool {
		allianceID, expiresAt, hasInvite = sg.PendingAllianceInvite()
		return true
	})
	if !hasInvite || !clock.Now().Before(expiresAt) {
		return nil
	}
	worldID := h.gs.config.WorldId
	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.AcceptAllianceInviteReply, error) {
		return h.gs.internalClient.AcceptAllianceInvite(c, &internal.AcceptAllianceInviteRequest{
			WorldId:     worldID,
			CharacterId: charID,
			AllianceId:  allianceID,
			GuildId:     guildID,
		})
	}, func(reply *internal.AcceptAllianceInviteReply) error {
		if reply == nil || !reply.GetOk() {
			if reply != nil && !reply.GetOk() {
				log.Printf("AllianceOperation(accept invite): failed character=%d alliance=%d guild=%d code=%v",
					charID, allianceID, guildID, reply.GetErrorCode())
			}
			return nil
		}
		h.withStoredGuild(guildID, func(g *entity.Guild) bool {
			g.ClearAllianceInvite(allianceID)
			return true
		})
		alliancePb := reply.GetAlliance()
		if alliancePb == nil {
			log.Printf("AllianceOperation(accept invite): ok but missing alliance character=%d", charID)
			return nil
		}
		h.gs.alliance.Update(alliancePb)
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("AllianceOperation(accept invite) async error: %v", err)
	})
	return nil
}

func (h *AllianceOperation) handleDenyInvite(ctx *core.ClientContext, ch *entity.Character) error {
	if h.gs == nil {
		return nil
	}
	guildID, inGuild := ch.GetGuildID()
	if !inGuild {
		return nil
	}
	g := h.gs.guild.Get(guildID)
	if g == nil {
		return nil
	}
	charID := ch.GetID()
	if g.LeaderCharacterID != charID {
		return nil
	}
	m := g.FindMember(charID)
	if m == nil || m.GetRank() != internal.GuildMemberRank_GUILD_MEMBER_RANK_MASTER {
		return nil
	}
	if _, inAlliance := g.GetAllianceID(); inAlliance {
		return nil
	}
	var allianceID uint32
	var hasInvite bool
	h.withStoredGuild(guildID, func(sg *entity.Guild) bool {
		allianceID, _, hasInvite = sg.PendingAllianceInvite()
		return true
	})
	if !hasInvite {
		return nil
	}
	guildName := g.Name
	h.withStoredGuild(guildID, func(sg *entity.Guild) bool {
		sg.ClearAllianceInvite(allianceID)
		return true
	})
	if allianceID != 0 {
		h.notifyAllianceLeaderInviteDenied(allianceID, guildName)
	}
	return nil
}

func (h *AllianceOperation) handleLeave(ctx *core.ClientContext, ch *entity.Character) error {
	if h.gs == nil || h.gs.internalClient == nil {
		return nil
	}
	worldID := h.gs.config.WorldId
	charID := ch.GetID()
	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.LeaveAllianceReply, error) {
		return h.gs.internalClient.LeaveAlliance(c, &internal.LeaveAllianceRequest{
			WorldId:     worldID,
			CharacterId: charID,
		})
	}, func(reply *internal.LeaveAllianceReply) error {
		if reply == nil || !reply.GetOk() {
			if reply != nil && !reply.GetOk() {
				log.Printf("AllianceOperation(leave): failed character=%d code=%v", charID, reply.GetErrorCode())
			}
			return nil
		}
		if reply.GetDisbanded() {
			allianceID := reply.GetAllianceId()
			guildIDs := h.gs.alliance.GuildIDs(allianceID)
			if len(guildIDs) == 0 && reply.GetRemovedGuildId() != 0 {
				guildIDs = []uint32{reply.GetRemovedGuildId()}
			}
			h.gs.alliance.DisbandAfterGuildRefreshAsync(ctx.ActorContext, allianceID, guildIDs)
			return nil
		}
		alliancePb := reply.GetAlliance()
		if alliancePb == nil {
			return nil
		}
		removedGuildID := reply.GetRemovedGuildId()
		actorCtx := ctx.ActorContext
		if removedGuildID == 0 {
			h.gs.alliance.BroadcastGuildLeft(alliancePb, nil, false)
			return nil
		}
		h.gs.guild.RefreshAsync(actorCtx, removedGuildID).Then(func(interface{}) (interface{}, error) {
			var removedGuildPb *internal.Guild
			if g := h.gs.guild.Get(removedGuildID); g != nil {
				removedGuildPb = g.ToProto()
			}
			h.gs.alliance.BroadcastGuildLeft(alliancePb, removedGuildPb, false)
			return nil, nil
		})
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("AllianceOperation(leave) async error: %v", err)
	})
	return nil
}

func (h *AllianceOperation) handleCreate(ctx *core.ClientContext, ch *entity.Character, req *request.AllianceOperation) error {
	if h.gs == nil || h.gs.internalClient == nil {
		return nil
	}
	partnerCharacterID, ok := h.gs.alliance.ValidateCreateAlliance(ch)
	if !ok || partnerCharacterID == 0 {
		return nil
	}

	worldID := h.gs.config.WorldId
	charID := ch.GetID()
	allianceName := req.AllianceName
	if allianceName == "" {
		return nil
	}

	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.CreateAllianceReply, error) {
		return h.gs.internalClient.CreateAlliance(c, &internal.CreateAllianceRequest{
			WorldId:            worldID,
			AllianceName:       allianceName,
			LeaderCharacterId:  charID,
			PartnerCharacterId: partnerCharacterID,
		})
	}, func(reply *internal.CreateAllianceReply) error {
		if !reply.GetOk() {
			log.Printf("AllianceOperation(create): failed character=%d code=%v", charID, reply.GetErrorCode())
			return nil
		}
		alliancePb := reply.GetAlliance()
		if alliancePb == nil {
			log.Printf("AllianceOperation(create): ok but missing alliance character=%d", charID)
			return nil
		}
		h.gs.alliance.Update(alliancePb)
		h.gs.alliance.BroadcastCreate(alliancePb)
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("AllianceOperation(create) async error: %v", err)
	})
	return nil
}

func (h *AllianceOperation) withStoredGuild(guildID uint32, fn func(*entity.Guild) bool) bool {
	if h.gs == nil || h.gs.guild == nil || guildID == 0 || fn == nil {
		return false
	}
	gc := h.gs.guild
	gc.mu.Lock()
	defer gc.mu.Unlock()
	g := gc.guilds[guildID]
	if g == nil {
		return false
	}
	return fn(g)
}

func (h *AllianceOperation) cachedAlliance(allianceID uint32) *entity.Alliance {
	if h.gs == nil || allianceID == 0 {
		return nil
	}
	return h.gs.GetAllianceSystem().Get(allianceID)
}

func (h *AllianceOperation) inviterMaySendAllianceInvite(allianceID uint32, inviterCharacterID uint32) (*entity.Alliance, bool) {
	a := h.cachedAlliance(allianceID)
	if a == nil {
		return nil, false
	}
	if !a.HasInviteCapacity() {
		return nil, false
	}
	if a.LeaderCharacterID != inviterCharacterID {
		return nil, false
	}
	return a, true
}

func (h *AllianceOperation) notifyAllianceLeaderInviteDenied(allianceID uint32, guildName string) {
	if h.gs == nil || allianceID == 0 || guildName == "" {
		return
	}
	a := h.cachedAlliance(allianceID)
	if a == nil || a.LeaderCharacterID == 0 {
		return
	}
	leaderID := a.LeaderCharacterID
	if h.gs.characterRuntime == nil || !h.gs.characterRuntime.Exists(leaderID) {
		return
	}
	h.gs.EnsureSend(nil, leaderID, &g_actor.DeliverMessage{
		CharacterID: leaderID,
		MessageType: gameconst.MsgPinkText,
		Message:     guildName + gameconst.AllianceInviteDeniedMessageSuffix,
	})
}
