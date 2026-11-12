package server

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/clock"
	gameconst "github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
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
	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}

	switch req.Operation {
	case pconst.AllianceC2SLoadInfo:
		if ch.Listener() != nil {
			ch.Listener().OnAllianceUpdateInfo(ch)
		}
		return nil
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

// The changes below leave who may make them to internal, which checks the latest rows under its lock.
func (h *AllianceOperation) handleExpel(ctx *core.ClientContext, ch *game.Character, req *request.AllianceOperation) error {
	if _, inGuild := ch.Guild().ID(); !inGuild || req.TargetGuildID == 0 {
		return nil
	}
	h.gs.systems.Get[*game.AllianceSystem]().ExpelAsync(ctx.ActorContext, ch, req.TargetGuildID, req.AllianceID)
	return nil
}

func (h *AllianceOperation) handleChangeLeader(ctx *core.ClientContext, ch *game.Character, req *request.AllianceOperation) error {
	if _, inGuild := ch.Guild().ID(); !inGuild || req.NewLeaderID == 0 {
		return nil
	}
	h.gs.systems.Get[*game.AllianceSystem]().ChangeLeaderAsync(ctx.ActorContext, ch, req.NewLeaderID)
	return nil
}

func (h *AllianceOperation) handleChangeNotice(ctx *core.ClientContext, ch *game.Character, req *request.AllianceOperation) error {
	if _, inGuild := ch.Guild().ID(); !inGuild {
		return nil
	}
	h.gs.systems.Get[*game.AllianceSystem]().ChangeNoticeAsync(ctx.ActorContext, ch, req.Notice)
	return nil
}

func (h *AllianceOperation) handleChangeRankTitles(ctx *core.ClientContext, ch *game.Character, req *request.AllianceOperation) error {
	if _, inGuild := ch.Guild().ID(); !inGuild {
		return nil
	}
	h.gs.systems.Get[*game.AllianceSystem]().ChangeRankTitlesAsync(ctx.ActorContext, ch, req.RankTitles[:])
	return nil
}

func (h *AllianceOperation) handleChangeMemberRank(ctx *core.ClientContext, ch *game.Character, req *request.AllianceOperation) error {
	if _, inGuild := ch.Guild().ID(); !inGuild || req.TargetCharacterID == 0 {
		return nil
	}
	h.gs.systems.Get[*game.AllianceSystem]().ChangeMemberRankAsync(ctx.ActorContext, ch, req.TargetCharacterID, req.RankChangePromote)
	return nil
}

func (h *AllianceOperation) handleInvite(ctx *core.ClientContext, ch *game.Character, req *request.AllianceOperation) error {
	if h.gs == nil {
		return nil
	}
	guildID, inGuild := ch.Guild().ID()
	if !inGuild {
		return nil
	}
	g := h.gs.systems.Get[*game.GuildSystem]().Get(guildID)
	if g == nil {
		return nil
	}
	charID := ch.GetID()
	if !g.IsAllianceLeader(charID) {
		return nil
	}
	allianceID, inAlliance := g.GetAllianceID()
	if !inAlliance {
		return nil
	}
	a := h.gs.systems.Get[*game.AllianceSystem]().Get(allianceID)
	if a == nil || !a.CanSendInvite(charID) {
		return nil
	}
	guildName := req.TargetGuildLeaderName
	if guildName == "" {
		return nil
	}
	targetGuildID, nameFound := h.gs.systems.Get[*game.GuildSystem]().NameToGuildID(guildName)
	if !nameFound {
		ch.Listener().OnMessage(ch, gameconst.MsgPopup, gameconst.AllianceInviteGuildNotFoundMessage)
		return nil
	}
	targetG := h.gs.systems.Get[*game.GuildSystem]().Get(targetGuildID)
	if targetG == nil {
		ch.Listener().OnMessage(ch, gameconst.MsgPopup, gameconst.AllianceInviteGuildNotFoundMessage)
		return nil
	}
	if _, targetInAlliance := targetG.GetAllianceID(); targetInAlliance || targetG.LeaderCharacterID() == 0 {
		ch.Listener().OnMessage(ch, gameconst.MsgPopup, gameconst.AllianceInviteGuildNotFoundMessage)
		return nil
	}
	targetLeaderID := targetG.LeaderCharacterID()
	if !h.gs.systems.Get[*game.DispatchSystem]().Exists(targetLeaderID) {
		ch.Listener().OnMessage(ch, gameconst.MsgPopup, gameconst.AllianceInviteTargetNotOnlineMessage)
		return nil
	}
	h.gs.systems.Get[*game.DispatchSystem]().SendToCharacter(nil, targetLeaderID, &game.InviteToAlliance{
		CharacterID:        targetLeaderID,
		TargetGuildID:      targetGuildID,
		AllianceID:         allianceID,
		InviterCharacterID: charID,
		InviterGuildID:     guildID,
		InviterName:        ch.GetName(),
		AllianceName:       a.Name(),
		ExpiresAt:          clock.Now().Add(gameconst.GuildInviteDuration),
	})
	return nil
}

func (h *AllianceOperation) handleAcceptInvite(ctx *core.ClientContext, ch *game.Character) error {
	guildID, inGuild := ch.Guild().ID()
	if !inGuild {
		return nil
	}
	g := h.gs.systems.Get[*game.GuildSystem]().Get(guildID)
	if g == nil {
		return nil
	}
	charID := ch.GetID()
	if g.FindMaster(charID) == nil {
		return nil
	}
	if _, inAlliance := g.GetAllianceID(); inAlliance {
		return nil
	}
	allianceID, expiresAt, hasInvite := h.gs.systems.Get[*game.GuildSystem]().PendingAllianceInvite(guildID)
	if !hasInvite || !clock.Now().Before(expiresAt) {
		return nil
	}
	h.gs.systems.Get[*game.AllianceSystem]().AcceptInviteAsync(ctx.ActorContext, ch, guildID, allianceID)
	return nil
}

func (h *AllianceOperation) handleDenyInvite(ctx *core.ClientContext, ch *game.Character) error {
	if h.gs == nil {
		return nil
	}
	guildID, inGuild := ch.Guild().ID()
	if !inGuild {
		return nil
	}
	g := h.gs.systems.Get[*game.GuildSystem]().Get(guildID)
	if g == nil {
		return nil
	}
	charID := ch.GetID()
	if g.FindMaster(charID) == nil {
		return nil
	}
	if _, inAlliance := g.GetAllianceID(); inAlliance {
		return nil
	}
	allianceID, _, hasInvite := h.gs.systems.Get[*game.GuildSystem]().PendingAllianceInvite(guildID)
	if !hasInvite {
		return nil
	}
	guildName := g.Name()
	h.gs.systems.Get[*game.GuildSystem]().ClearAllianceInvite(guildID, allianceID)
	h.notifyDenied(allianceID, guildName)
	return nil
}

func (h *AllianceOperation) notifyDenied(allianceID uint32, guildName string) {
	if allianceID == 0 || guildName == "" {
		return
	}
	a := h.gs.systems.Get[*game.AllianceSystem]().Get(allianceID)
	if a == nil || a.LeaderCharacterID() == 0 {
		return
	}
	leaderID := a.LeaderCharacterID()
	if !h.gs.systems.Get[*game.DispatchSystem]().Exists(leaderID) {
		return
	}
	h.gs.systems.Get[*game.DispatchSystem]().Visit(leaderID, func(_ actor.Context, ch *game.Character) {
		if ch == nil {
			return
		}
		ch.Listener().OnMessage(ch, gameconst.MsgPinkText, guildName+gameconst.AllianceInviteDeniedMessageSuffix)
	})
}

func (h *AllianceOperation) handleLeave(ctx *core.ClientContext, ch *game.Character) error {
	h.gs.systems.Get[*game.AllianceSystem]().LeaveAsync(ctx.ActorContext, ch)
	return nil
}
