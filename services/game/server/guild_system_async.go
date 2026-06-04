package server

import (
	"context"
	"log"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

func (s guildSystem) DisbandAsync(ctx actor.Context, ch *entity.Character, result *int) *async.Promise {
	fail := int(constant.GuildDisbandResultFailed)
	if result == nil {
		result = &fail
	}
	if s.gs == nil || s.gs.internalClient == nil || ch == nil {
		*result = fail
		return nil
	}
	*result = fail
	worldID := s.gs.config.WorldId
	charID := ch.GetID()
	promise := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	return async.ThenRPC(promise, func(c context.Context) (*internal.DisbandGuildReply, error) {
		return s.gs.internalClient.DisbandGuild(c, &internal.DisbandGuildRequest{
			WorldId:     worldID,
			CharacterId: charID,
		})
	}, func(reply *internal.DisbandGuildReply) error {
		if reply == nil {
			*result = int(constant.GuildDisbandResultFailed)
			return nil
		}
		if reply.GetOk() {
			*result = int(constant.GuildDisbandResultOK)
			return nil
		}
		switch reply.GetErrorCode() {
		case internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD:
			*result = int(constant.GuildDisbandResultNotInGuild)
		case internal.GuildErrorCode_GUILD_ERROR_NOT_AUTHORIZED:
			*result = int(constant.GuildDisbandResultNotMaster)
		default:
			*result = int(constant.GuildDisbandResultFailed)
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("guild disband: character=%d: %v", charID, err)
		*result = int(constant.GuildDisbandResultFailed)
	})
}

func (s guildSystem) IncCapacityAsync(ctx actor.Context, ch *entity.Character, extendedCap bool, result *int) *async.Promise {
	fail := int(constant.GuildIncreaseCapacityResultFailed)
	if result == nil {
		result = &fail
	}
	if s.gs == nil || s.gs.internalClient == nil || ch == nil {
		*result = fail
		return nil
	}
	if !extendedCap && ch.Meso < constant.GuildCapacityIncreaseMesoCost {
		*result = int(constant.GuildIncreaseCapacityResultInsufficientMeso)
		return nil
	}
	*result = fail
	worldID := s.gs.config.WorldId
	charID := ch.GetID()
	promise := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	return async.ThenRPC(promise, func(c context.Context) (*internal.IncreaseGuildCapacityReply, error) {
		return s.gs.internalClient.IncreaseGuildCapacity(c, &internal.IncreaseGuildCapacityRequest{
			WorldId:     worldID,
			CharacterId: charID,
			ExtendedCap: extendedCap,
		})
	}, func(reply *internal.IncreaseGuildCapacityReply) error {
		if reply == nil {
			*result = int(constant.GuildIncreaseCapacityResultFailed)
			return nil
		}
		if reply.GetOk() {
			if !extendedCap {
				ch.RemoveMeso(constant.GuildCapacityIncreaseMesoCost)
			}
			*result = int(constant.GuildIncreaseCapacityResultOK)
			return nil
		}
		switch reply.GetErrorCode() {
		case internal.GuildErrorCode_GUILD_ERROR_NOT_IN_GUILD:
			*result = int(constant.GuildIncreaseCapacityResultNotInGuild)
		case internal.GuildErrorCode_GUILD_ERROR_NOT_AUTHORIZED:
			*result = int(constant.GuildIncreaseCapacityResultNotMaster)
		case internal.GuildErrorCode_GUILD_ERROR_CAPACITY_REACHED:
			*result = int(constant.GuildIncreaseCapacityResultCapacityReached)
		case internal.GuildErrorCode_GUILD_ERROR_INSUFFICIENT_GUILD_GP:
			*result = int(constant.GuildIncreaseCapacityResultInsufficientGP)
		default:
			*result = int(constant.GuildIncreaseCapacityResultFailed)
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("guild inc_capacity: character=%d: %v", charID, err)
		*result = int(constant.GuildIncreaseCapacityResultFailed)
	})
}

func (s guildSystem) CreateAllianceAsync(ctx actor.Context, ch *entity.Character, allianceName string, result *int) *async.Promise {
	fail := int(constant.AllianceCreateResultFailed)
	if result == nil {
		result = &fail
	}
	if s.gs == nil || s.gs.internalClient == nil || ch == nil {
		*result = fail
		return nil
	}
	partnerID, ok := s.gs.alliance.ValidateCreateAlliance(ch)
	if !ok || partnerID == 0 {
		*result = int(constant.AllianceCreateResultInvalidRequirements)
		return nil
	}
	trimmed := strings.TrimSpace(allianceName)
	if len(trimmed) < 3 || len(trimmed) > 12 {
		*result = int(constant.AllianceCreateResultInvalidName)
		return nil
	}
	if ch.Meso < constant.AllianceCreateMesoCost {
		*result = int(constant.AllianceCreateResultInsufficientMeso)
		return nil
	}
	*result = fail
	worldID := s.gs.config.WorldId
	charID := ch.GetID()
	promise := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	return async.ThenRPC(promise, func(c context.Context) (*internal.CreateAllianceReply, error) {
		return s.gs.internalClient.CreateAlliance(c, &internal.CreateAllianceRequest{
			WorldId:            worldID,
			AllianceName:       trimmed,
			LeaderCharacterId:  charID,
			PartnerCharacterId: partnerID,
		})
	}, func(reply *internal.CreateAllianceReply) error {
		if reply == nil {
			*result = int(constant.AllianceCreateResultFailed)
			return nil
		}
		if !reply.GetOk() {
			switch reply.GetErrorCode() {
			case internal.AllianceErrorCode_ALLIANCE_ERROR_ALLIANCE_NAME_INVALID:
				*result = int(constant.AllianceCreateResultInvalidName)
			case internal.AllianceErrorCode_ALLIANCE_ERROR_ALLIANCE_NAME_TAKEN:
				*result = int(constant.AllianceCreateResultNameTaken)
			default:
				*result = int(constant.AllianceCreateResultFailed)
			}
			return nil
		}
		alliancePb := reply.GetAlliance()
		if alliancePb == nil {
			*result = int(constant.AllianceCreateResultFailed)
			return nil
		}
		ch.RemoveMeso(constant.AllianceCreateMesoCost)
		s.gs.alliance.Update(alliancePb)
		s.gs.alliance.BroadcastCreate(alliancePb)
		*result = int(constant.AllianceCreateResultOK)
		return nil
	}).OnError(func(err error) {
		log.Printf("alliance create: character=%d: %v", charID, err)
		*result = int(constant.AllianceCreateResultFailed)
	})
}

func (s guildSystem) DisbandAllianceAsync(ctx actor.Context, ch *entity.Character, result *int) *async.Promise {
	fail := int(constant.AllianceDisbandResultFailed)
	if result == nil {
		result = &fail
	}
	if s.gs == nil || s.gs.internalClient == nil || ch == nil {
		*result = fail
		return nil
	}
	*result = fail
	guildID, ok := ch.GetGuildID()
	if !ok {
		*result = int(constant.AllianceDisbandResultNotInAlliance)
		return nil
	}
	g := s.gs.guild.Get(guildID)
	allianceID, inAlliance := g.GetAllianceID()
	if g == nil || !inAlliance {
		*result = int(constant.AllianceDisbandResultNotInAlliance)
		return nil
	}
	charID := ch.GetID()
	if g.LeaderCharacterID != charID {
		*result = int(constant.AllianceDisbandResultNotGuildMaster)
		return nil
	}
	m := g.FindMember(charID)
	if m == nil {
		*result = int(constant.AllianceDisbandResultNotLeader)
		return nil
	}
	ar, isAllianceLeader := m.GetAllianceRank()
	if !isAllianceLeader || ar != 1 {
		*result = int(constant.AllianceDisbandResultNotLeader)
		return nil
	}
	worldID := s.gs.config.WorldId
	promise := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	return async.ThenRPC(promise, func(c context.Context) (*internal.DisbandAllianceReply, error) {
		return s.gs.internalClient.DisbandAlliance(c, &internal.DisbandAllianceRequest{
			WorldId:     worldID,
			CharacterId: charID,
		})
	}, func(reply *internal.DisbandAllianceReply) error {
		if reply == nil {
			*result = int(constant.AllianceDisbandResultFailed)
			return nil
		}
		if !reply.GetOk() {
			switch reply.GetErrorCode() {
			case internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_IN_ALLIANCE,
				internal.AllianceErrorCode_ALLIANCE_ERROR_ALLIANCE_NOT_FOUND:
				*result = int(constant.AllianceDisbandResultNotInAlliance)
			case internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_ALLIANCE_LEADER,
				internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_GUILD_MASTER:
				*result = int(constant.AllianceDisbandResultNotLeader)
			default:
				*result = int(constant.AllianceDisbandResultFailed)
			}
			return nil
		}
		guildIDs := s.gs.alliance.GuildIDs(allianceID)
		if len(guildIDs) == 0 {
			guildIDs = []uint32{guildID}
		}
		s.gs.alliance.DisbandAfterGuildRefreshAsync(ctx, allianceID, guildIDs).Run()
		*result = int(constant.AllianceDisbandResultOK)
		return nil
	}).OnError(func(err error) {
		log.Printf("alliance disband: character=%d: %v", charID, err)
		*result = int(constant.AllianceDisbandResultFailed)
	})
}
