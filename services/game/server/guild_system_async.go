package server

import (
	"context"
	"log"

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
