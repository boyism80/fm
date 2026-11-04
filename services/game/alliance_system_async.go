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

func (s *AllianceSystem) IncCapacityAsync(ctx actor.Context, ch *Character) (*async.Promise[int], int) {
	guildID, inGuild := ch.guild.ID()
	if !inGuild {
		return nil, int(constant.AllianceIncreaseCapacityResultNotInAlliance)
	}
	g := s.gw.Systems().Get[*GuildSystem]().Get(guildID)
	if g == nil {
		return nil, int(constant.AllianceIncreaseCapacityResultNotInAlliance)
	}
	if _, inAlliance := g.GetAllianceID(); !inAlliance {
		return nil, int(constant.AllianceIncreaseCapacityResultNotInAlliance)
	}
	charID := ch.GetID()
	if g.leaderCharacterID != charID {
		return nil, int(constant.AllianceIncreaseCapacityResultNotGuildMaster)
	}
	if !g.IsAllianceLeader(charID) {
		return nil, int(constant.AllianceIncreaseCapacityResultNotLeader)
	}
	if !ch.inventory.RemoveMeso(constant.AllianceIncreaseCapacityMesoCost) {
		return nil, int(constant.AllianceIncreaseCapacityResultInsufficientMeso)
	}
	allianceID, _ := g.GetAllianceID()
	detail := allianceCapacityDetail{Meso: constant.AllianceIncreaseCapacityMesoCost}
	id := ch.unsettled.Add(internal.UnsettledType_UNSETTLED_TYPE_ALLIANCE_CAPACITY, uint64(allianceID), detail)
	return s.incCapacityAsync(ctx, ch, id, allianceID).Bind(ch).Then(func(reply *internal.IncreaseAllianceCapacityReply) (int, error) {
		ch.unsettled.Erase(id)
		if !reply.GetOk() {
			ch.inventory.AddMeso(detail.Meso)
		}

		if !reply.GetOk() {
			switch reply.GetErrorCode() {
			case internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_IN_ALLIANCE,
				internal.AllianceErrorCode_ALLIANCE_ERROR_ALLIANCE_NOT_FOUND,
				internal.AllianceErrorCode_ALLIANCE_ERROR_GUILD_NOT_FOUND:
				return int(constant.AllianceIncreaseCapacityResultNotInAlliance), nil
			case internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_ALLIANCE_LEADER,
				internal.AllianceErrorCode_ALLIANCE_ERROR_NOT_GUILD_MASTER:
				return int(constant.AllianceIncreaseCapacityResultNotLeader), nil
			case internal.AllianceErrorCode_ALLIANCE_ERROR_CAPACITY_MAX:
				return int(constant.AllianceIncreaseCapacityResultCapacityMax), nil
			default:
				return int(constant.AllianceIncreaseCapacityResultFailed), nil
			}
		}
		if alliancePb := reply.GetAlliance(); alliancePb != nil {
			s.Update(alliancePb)
			s.BroadcastInfoUpdate(alliancePb)
		}
		return int(constant.AllianceIncreaseCapacityResultOK), nil
	}).OnError(func(err error) {
		ch.unsettled.Find(id).inFlight = false

		log.Printf("alliance inc_capacity: character=%d id=%d: %v", charID, id, err)
	}), int(constant.AllianceIncreaseCapacityResultFailed)
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
