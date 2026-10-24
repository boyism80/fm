package game

import (
	"context"
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (s *AllianceSystem) IncCapacityAsync(ctx actor.Context, ch *Character, result *int) *async.Promise[*internal.IncreaseAllianceCapacityReply] {
	fail := int(constant.AllianceIncreaseCapacityResultFailed)
	if result == nil {
		result = &fail
	}
	if ch == nil {
		*result = fail
		return nil
	}
	*result = fail
	guildID, inGuild := ch.guild.ID()
	if !inGuild {
		*result = int(constant.AllianceIncreaseCapacityResultNotInAlliance)
		return nil
	}
	g := s.gw.Systems().Get[*GuildSystem]().Get(guildID)
	if g == nil {
		*result = int(constant.AllianceIncreaseCapacityResultNotInAlliance)
		return nil
	}
	if _, inAlliance := g.GetAllianceID(); !inAlliance {
		*result = int(constant.AllianceIncreaseCapacityResultNotInAlliance)
		return nil
	}
	charID := ch.GetID()
	if g.leaderCharacterID != charID {
		*result = int(constant.AllianceIncreaseCapacityResultNotGuildMaster)
		return nil
	}
	if !g.IsAllianceLeader(charID) {
		*result = int(constant.AllianceIncreaseCapacityResultNotLeader)
		return nil
	}
	if !ch.inventory.RemoveMeso(constant.AllianceIncreaseCapacityMesoCost) {
		*result = int(constant.AllianceIncreaseCapacityResultInsufficientMeso)
		return nil
	}
	allianceID, _ := g.GetAllianceID()
	detail := allianceCapacityDetail{Meso: constant.AllianceIncreaseCapacityMesoCost}
	id := ch.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_ALLIANCE_CAPACITY, uint64(allianceID), detail)
	return s.incCapacityAsync(ctx, ch, id, allianceID).Bind(ch).Do(func(reply *internal.IncreaseAllianceCapacityReply) error {
		ch.unsettled.Erase(id)
		if !reply.GetOk() {
			ch.inventory.AddMeso(detail.Meso)
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
		*result = int(constant.AllianceIncreaseCapacityResultOK)
		if alliancePb := reply.GetAlliance(); alliancePb != nil {
			s.Update(alliancePb)
			s.BroadcastInfoUpdate(alliancePb)
		}
		return nil
	}).OnError(func(err error) {
		ch.unsettled.Find(id).inFlight = false

		log.Printf("alliance inc_capacity: character=%d id=%d: %v", charID, id, err)
		*result = int(constant.AllianceIncreaseCapacityResultFailed)
	})
}

type allianceCapacityDetail struct {
	Meso int32
}

func (s *AllianceSystem) incCapacityAsync(ctx actor.Context, ch *Character, requestID uint64, allianceID uint32) *async.Promise[*internal.IncreaseAllianceCapacityReply] {
	req := &internal.IncreaseAllianceCapacityRequest{
		WorldId:     s.gw.GetWorldID(),
		CharacterId: ch.GetID(),
		RequestId:   requestID,
		AllianceId:  allianceID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).Bind(ch).ThenRPC(func(c context.Context) (*internal.IncreaseAllianceCapacityReply, error) {
		return s.gw.Systems().Get[*InternalSystem]().client.IncreaseAllianceCapacity(c, req)
	}, nil)
}

func (s *AllianceSystem) SettleCapacity(ctx actor.Context, ch *Character, request *UnsettledRequest) {
	var detail allianceCapacityDetail
	if err := json.Unmarshal([]byte(request.detail), &detail); err != nil {
		log.Printf("alliance settle capacity: character=%d id=%d: %v", ch.GetID(), request.id, err)
		return
	}
	s.incCapacityAsync(ctx, ch, request.id, uint32(request.target)).Bind(ch).Do(func(reply *internal.IncreaseAllianceCapacityReply) error {
		ch.unsettled.Erase(request.id)
		if !reply.GetOk() {
			ch.inventory.AddMeso(detail.Meso)
			return nil
		}
		if alliancePb := reply.GetAlliance(); alliancePb != nil {
			s.Update(alliancePb)
			s.BroadcastInfoUpdate(alliancePb)
		}
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("alliance settle capacity: character=%d id=%d: %v", ch.GetID(), request.id, err)
	})
}
