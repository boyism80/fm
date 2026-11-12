package game

import (
	"context"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

// The guild and alliance changes below leave every permission to internal, which checks the latest rows under its
// lock. Its refusals of a change the member may not make are the expected answer and send the client nothing.

type guildReply interface {
	GetOk() bool
	GetErrorCode() internal.GuildErrorCode
}

type allianceReply interface {
	GetOk() bool
	GetErrorCode() internal.AllianceErrorCode
}

func guildRefusal(code internal.GuildErrorCode) bool {
	switch code {
	case internal.GuildErrorCode_GUILD_ERROR_ALREADY_IN_GUILD,
		internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD,
		internal.GuildErrorCode_GUILD_ERROR_LEADER_CANNOT_LEAVE,
		internal.GuildErrorCode_GUILD_ERROR_NOT_AUTHORIZED,
		internal.GuildErrorCode_GUILD_ERROR_CANNOT_EXPEL_TARGET,
		internal.GuildErrorCode_GUILD_ERROR_TARGET_NOT_IN_GUILD,
		internal.GuildErrorCode_GUILD_ERROR_INVALID_RANK_TITLES,
		internal.GuildErrorCode_GUILD_ERROR_INVALID_MEMBER_RANK,
		internal.GuildErrorCode_GUILD_ERROR_INVALID_NOTICE:
		return true
	}
	return false
}

func allianceRefusal(code internal.AllianceErrorCode) bool {
	switch code {
	case internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_GUILD_MASTER,
		internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_IN_ALLIANCE,
		internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_ALLIANCE_LEADER,
		internal.AllianceErrorCode_ALLIANCE_ERROR_ALLIANCE_NOT_FOUND,
		internal.AllianceErrorCode_ALLIANCE_ERROR_INVALID_RANK_TITLES,
		internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_AUTHORIZED,
		internal.AllianceErrorCode_ALLIANCE_ERROR_INVALID_MEMBER_RANK,
		internal.AllianceErrorCode_ALLIANCE_ERROR_TARGET_NOT_IN_ALLIANCE,
		internal.AllianceErrorCode_ALLIANCE_ERROR_CANNOT_EXPEL_OWN_GUILD,
		internal.AllianceErrorCode_ALLIANCE_ERROR_INVALID_NOTICE,
		internal.AllianceErrorCode_ALLIANCE_ERROR_INVALID_LEADER_CANDIDATE:
		return true
	}
	return false
}

// guildChange sends a guild change and logs only a failure that is not a refusal.
func guildChange[R guildReply](ctx actor.Context, what string, characterID uint32, call func(context.Context) (R, error)) {
	async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(call, func(reply R) error {
		if !reply.GetOk() && !guildRefusal(reply.GetErrorCode()) {
			log.Printf("guild %s: failed character=%d code=%v", what, characterID, reply.GetErrorCode())
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("guild %s character=%d: %v", what, characterID, err)
	})
}

// allianceChange sends an alliance change and logs only a failure that is not a refusal.
func allianceChange[R allianceReply](ctx actor.Context, what string, characterID uint32, call func(context.Context) (R, error)) {
	async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(call, func(reply R) error {
		if !reply.GetOk() && !allianceRefusal(reply.GetErrorCode()) {
			log.Printf("alliance %s: failed character=%d code=%v", what, characterID, reply.GetErrorCode())
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("alliance %s character=%d: %v", what, characterID, err)
	})
}

func (s *GuildSystem) client() internal.InternalClient {
	return s.gw.Systems().Get[*InternalSystem]().client
}

// AcceptInviteAsync joins the guild that invited ch and shows it the guild and its alliance.
func (s *GuildSystem) AcceptInviteAsync(ctx actor.Context, ch *Character, guildID uint32) {
	worldID := s.gw.GetWorldID()
	charID := ch.GetID()
	member := GuildMemberFromCharacter(ch, worldID, s.gw.GetChannelID(), internal.GuildMemberRank_GUILD_MEMBER_RANK_NEW).BuildRPC()
	async.NewTask(ctx, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(func(c context.Context) (*internal.AcceptGuildInviteReply, error) {
		return s.client().AcceptGuildInvite(c, &internal.AcceptGuildInviteRequest{WorldId: worldID, GuildId: guildID, Member: member})
	}, func(reply *internal.AcceptGuildInviteReply) error {
		if !reply.GetOk() {
			if reply.GetErrorCode() == internal.GuildErrorCode_GUILD_ERROR_GUILD_FULL {
				ch.listener.OnMessage(ch, constant.MsgPopup, constant.GuildFullMessage)
			} else if !guildRefusal(reply.GetErrorCode()) {
				log.Printf("guild accept invite: failed character=%d guild=%d code=%v", charID, guildID, reply.GetErrorCode())
			}
			return nil
		}
		guildPb := reply.GetGuild()
		if guildPb == nil {
			log.Printf("guild accept invite: ok but missing guild character=%d", charID)
			return nil
		}
		s.StoreGuild(ch, guildPb)
		if alliancePb := reply.GetAlliance(); alliancePb != nil {
			s.gw.Systems().Get[*AllianceSystem]().Update(alliancePb)
			ch.listener.OnShowAllianceInfo(ch)
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("guild accept invite character=%d: %v", charID, err)
	})
}

func (s *GuildSystem) LeaveAsync(ctx actor.Context, ch *Character) {
	req := &internal.LeaveGuildRequest{WorldId: s.gw.GetWorldID(), CharacterId: ch.GetID()}
	guildChange(ctx, "leave", ch.GetID(), func(c context.Context) (*internal.LeaveGuildReply, error) {
		return s.client().LeaveGuild(c, req)
	})
}

func (s *GuildSystem) ExpelAsync(ctx actor.Context, ch *Character, targetID uint32) {
	req := &internal.ExpelGuildRequest{WorldId: s.gw.GetWorldID(), RequesterCharacterId: ch.GetID(), TargetCharacterId: targetID}
	guildChange(ctx, "expel", ch.GetID(), func(c context.Context) (*internal.ExpelGuildReply, error) {
		return s.client().ExpelGuild(c, req)
	})
}

func (s *GuildSystem) ChangeRankTitlesAsync(ctx actor.Context, ch *Character, titles []string) {
	req := &internal.ChangeGuildRankTitlesRequest{WorldId: s.gw.GetWorldID(), CharacterId: ch.GetID(), RankTitles: titles}
	guildChange(ctx, "change rank titles", ch.GetID(), func(c context.Context) (*internal.ChangeGuildRankTitlesReply, error) {
		return s.client().ChangeGuildRankTitles(c, req)
	})
}

func (s *GuildSystem) ChangeMemberRankAsync(ctx actor.Context, ch *Character, targetID uint32, rank internal.GuildMemberRank) {
	req := &internal.ChangeGuildMemberRankRequest{
		WorldId:              s.gw.GetWorldID(),
		RequesterCharacterId: ch.GetID(),
		TargetCharacterId:    targetID,
		NewRank:              rank,
	}
	guildChange(ctx, "change member rank", ch.GetID(), func(c context.Context) (*internal.ChangeGuildMemberRankReply, error) {
		return s.client().ChangeGuildMemberRank(c, req)
	})
}

func (s *GuildSystem) ChangeNoticeAsync(ctx actor.Context, ch *Character, notice string) {
	req := &internal.ChangeGuildNoticeRequest{WorldId: s.gw.GetWorldID(), CharacterId: ch.GetID(), Notice: notice}
	guildChange(ctx, "change notice", ch.GetID(), func(c context.Context) (*internal.ChangeGuildNoticeReply, error) {
		return s.client().ChangeGuildNotice(c, req)
	})
}

func (s *AllianceSystem) client() internal.InternalClient {
	return s.gw.Systems().Get[*InternalSystem]().client
}

// AcceptInviteAsync joins ch's guild to the alliance that invited it; a failure tells ch the join did not go through.
func (s *AllianceSystem) AcceptInviteAsync(ctx actor.Context, ch *Character, guildID uint32, allianceID uint32) {
	charID := ch.GetID()
	req := &internal.AcceptAllianceInviteRequest{WorldId: s.gw.GetWorldID(), CharacterId: charID, AllianceId: allianceID, GuildId: guildID}
	async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(func(c context.Context) (*internal.AcceptAllianceInviteReply, error) {
		return s.client().AcceptAllianceInvite(c, req)
	}, func(reply *internal.AcceptAllianceInviteReply) error {
		if !reply.GetOk() {
			log.Printf("alliance accept invite: failed character=%d alliance=%d guild=%d code=%v", charID, allianceID, guildID, reply.GetErrorCode())
			s.gw.Systems().Get[*DispatchSystem]().Visit(charID, func(_ actor.Context, ch *Character) {
				if ch != nil {
					ch.listener.OnMessage(ch, constant.MsgPopup, constant.AllianceAcceptFailedMessage)
				}
			})
			return nil
		}
		s.gw.Systems().Get[*GuildSystem]().ClearAllianceInvite(guildID, allianceID)
		if alliancePb := reply.GetAlliance(); alliancePb != nil {
			s.Update(alliancePb)
		} else {
			log.Printf("alliance accept invite: ok but missing alliance character=%d", charID)
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("alliance accept invite character=%d: %v", charID, err)
	})
}

func (s *AllianceSystem) LeaveAsync(ctx actor.Context, ch *Character) {
	req := &internal.LeaveAllianceRequest{WorldId: s.gw.GetWorldID(), CharacterId: ch.GetID()}
	allianceChange(ctx, "leave", ch.GetID(), func(c context.Context) (*internal.LeaveAllianceReply, error) {
		return s.client().LeaveAlliance(c, req)
	})
}

// ExpelAsync expels a guild from the requester's alliance; allianceID is the one the client named, or 0.
func (s *AllianceSystem) ExpelAsync(ctx actor.Context, ch *Character, targetGuildID uint32, allianceID uint32) {
	req := &internal.ExpelAllianceGuildRequest{
		WorldId:       s.gw.GetWorldID(),
		CharacterId:   ch.GetID(),
		TargetGuildId: targetGuildID,
		AllianceId:    allianceID,
	}
	allianceChange(ctx, "expel", ch.GetID(), func(c context.Context) (*internal.ExpelAllianceGuildReply, error) {
		return s.client().ExpelAllianceGuild(c, req)
	})
}

func (s *AllianceSystem) ChangeLeaderAsync(ctx actor.Context, ch *Character, newLeaderID uint32) {
	req := &internal.ChangeAllianceLeaderRequest{WorldId: s.gw.GetWorldID(), CharacterId: ch.GetID(), NewLeaderCharacterId: newLeaderID}
	allianceChange(ctx, "change leader", ch.GetID(), func(c context.Context) (*internal.ChangeAllianceLeaderReply, error) {
		return s.client().ChangeAllianceLeader(c, req)
	})
}

func (s *AllianceSystem) ChangeNoticeAsync(ctx actor.Context, ch *Character, notice string) {
	req := &internal.ChangeAllianceNoticeRequest{WorldId: s.gw.GetWorldID(), CharacterId: ch.GetID(), Notice: notice}
	allianceChange(ctx, "change notice", ch.GetID(), func(c context.Context) (*internal.ChangeAllianceNoticeReply, error) {
		return s.client().ChangeAllianceNotice(c, req)
	})
}

func (s *AllianceSystem) ChangeRankTitlesAsync(ctx actor.Context, ch *Character, titles []string) {
	req := &internal.ChangeAllianceRankTitlesRequest{WorldId: s.gw.GetWorldID(), CharacterId: ch.GetID(), RankTitles: titles}
	allianceChange(ctx, "change rank titles", ch.GetID(), func(c context.Context) (*internal.ChangeAllianceRankTitlesReply, error) {
		return s.client().ChangeAllianceRankTitles(c, req)
	})
}

func (s *AllianceSystem) ChangeMemberRankAsync(ctx actor.Context, ch *Character, targetID uint32, promote bool) {
	req := &internal.ChangeAllianceMemberRankRequest{
		WorldId:              s.gw.GetWorldID(),
		RequesterCharacterId: ch.GetID(),
		TargetCharacterId:    targetID,
		Promote:              promote,
	}
	allianceChange(ctx, "change member rank", ch.GetID(), func(c context.Context) (*internal.ChangeAllianceMemberRankReply, error) {
		return s.client().ChangeAllianceMemberRank(c, req)
	})
}
