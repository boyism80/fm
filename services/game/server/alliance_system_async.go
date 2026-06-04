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

func (s allianceSystem) IncCapacityAsync(ctx actor.Context, ch *entity.Character, result *int) *async.Promise {
	fail := int(constant.AllianceIncreaseCapacityResultFailed)
	if result == nil {
		result = &fail
	}
	if s.gs == nil || s.gs.internalClient == nil || ch == nil {
		*result = fail
		return nil
	}
	*result = fail
	guildID, inGuild := ch.GetGuildID()
	if !inGuild {
		*result = int(constant.AllianceIncreaseCapacityResultNotInAlliance)
		return nil
	}
	g := s.gs.guild.Get(guildID)
	if g == nil {
		*result = int(constant.AllianceIncreaseCapacityResultNotInAlliance)
		return nil
	}
	if _, inAlliance := g.GetAllianceID(); !inAlliance {
		*result = int(constant.AllianceIncreaseCapacityResultNotInAlliance)
		return nil
	}
	charID := ch.GetID()
	if g.LeaderCharacterID != charID {
		*result = int(constant.AllianceIncreaseCapacityResultNotGuildMaster)
		return nil
	}
	m := g.FindMember(charID)
	if m == nil {
		*result = int(constant.AllianceIncreaseCapacityResultNotLeader)
		return nil
	}
	ar, isAllianceLeader := m.GetAllianceRank()
	if !isAllianceLeader || ar != 1 {
		*result = int(constant.AllianceIncreaseCapacityResultNotLeader)
		return nil
	}
	worldID := s.gs.config.WorldId
	promise := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	return async.ThenRPC(promise, func(c context.Context) (*internal.IncreaseAllianceCapacityReply, error) {
		return s.gs.internalClient.IncreaseAllianceCapacity(c, &internal.IncreaseAllianceCapacityRequest{
			WorldId:     worldID,
			CharacterId: charID,
		})
	}, func(reply *internal.IncreaseAllianceCapacityReply) error {
		if reply == nil {
			*result = int(constant.AllianceIncreaseCapacityResultFailed)
			return nil
		}
		if !reply.GetOk() {
			switch reply.GetErrorCode() {
			case internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_IN_ALLIANCE,
				internal.AllianceErrorCode_ALLIANCE_ERROR_ALLIANCE_NOT_FOUND,
				internal.AllianceErrorCode_ALLIANCE_ERROR_GUILD_NOT_FOUND:
				*result = int(constant.AllianceIncreaseCapacityResultNotInAlliance)
			case internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_ALLIANCE_LEADER,
				internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_GUILD_MASTER:
				*result = int(constant.AllianceIncreaseCapacityResultNotLeader)
			case internal.AllianceErrorCode_ALLIANCE_ERROR_CAPACITY_MAX:
				*result = int(constant.AllianceIncreaseCapacityResultCapacityMax)
			default:
				*result = int(constant.AllianceIncreaseCapacityResultFailed)
			}
			return nil
		}
		alliancePb := reply.GetAlliance()
		if alliancePb == nil {
			*result = int(constant.AllianceIncreaseCapacityResultFailed)
			return nil
		}
		s.gs.alliance.Update(alliancePb)
		s.gs.alliance.BroadcastCapacityChanged(alliancePb)
		*result = int(constant.AllianceIncreaseCapacityResultOK)
		return nil
	}).OnError(func(err error) {
		log.Printf("alliance inc_capacity: character=%d: %v", charID, err)
		*result = int(constant.AllianceIncreaseCapacityResultFailed)
	})
}
