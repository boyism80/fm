package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	amqp "github.com/rabbitmq/amqp091-go"
)

type allianceMqGuildAdded struct{ gs *GameServer }

func (allianceMqGuildAdded) New(gs *GameServer) *allianceMqGuildAdded {
	return &allianceMqGuildAdded{gs: gs}
}

func (*allianceMqGuildAdded) EventType() string {
	return "guild_added"
}

func (h *allianceMqGuildAdded) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeAllianceEventEnvelope(raw)
	if !ok {
		return nil
	}
	alliancePb, ok := allianceFromMQPayload(raw)
	if !ok {
		log.Printf("alliance consumer: guild_added alliance_id=%d missing alliance_pb", evt.AllianceID)
		return nil
	}
	var payload struct {
		AddedGuildID uint32 `json:"added_guild_id"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil || payload.AddedGuildID == 0 {
		log.Printf("alliance consumer: guild_added alliance_id=%d missing added_guild_id", evt.AllianceID)
		return nil
	}
	addedGuildPb := guildFromAlliancePb(alliancePb, payload.AddedGuildID)
	if addedGuildPb == nil {
		log.Printf("alliance consumer: guild_added alliance_id=%d guild_id=%d not in alliance_pb", evt.AllianceID, payload.AddedGuildID)
		return nil
	}
	gs.alliance.UpdateAsync(ctx, evt).Then(func(interface{}) (interface{}, error) {
		gs.alliance.Update(alliancePb)
		gs.alliance.BroadcastGuildAdded(alliancePb, addedGuildPb)
		return nil, nil
	})
	log.Printf("alliance consumer: applied guild_added alliance_id=%d revision=%d", evt.AllianceID, evt.Revision)
	return nil
}

func guildFromAlliancePb(alliancePb *internal.Alliance, guildID uint32) *internal.Guild {
	if alliancePb == nil || guildID == 0 {
		return nil
	}
	for _, g := range alliancePb.GetGuilds() {
		if g != nil && g.GetGuildId() == guildID {
			return g
		}
	}
	return nil
}
