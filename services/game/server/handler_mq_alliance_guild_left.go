package server

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type allianceMqGuildLeft struct{ gs *GameServer }

func (allianceMqGuildLeft) New(gs *GameServer) *allianceMqGuildLeft {
	return &allianceMqGuildLeft{gs: gs}
}

func (*allianceMqGuildLeft) EventType() string {
	return "guild_left"
}

func allianceGuildFromMQPayload(raw json.RawMessage) (*internal.Guild, bool) {
	var payload struct {
		RemovedGuildPB string `json:"removed_guild_pb"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil || payload.RemovedGuildPB == "" {
		return nil, false
	}
	wire, err := base64.StdEncoding.DecodeString(payload.RemovedGuildPB)
	if err != nil {
		return nil, false
	}
	var guildPb internal.Guild
	if err := proto.Unmarshal(wire, &guildPb); err != nil {
		return nil, false
	}
	return &guildPb, true
}

func (h *allianceMqGuildLeft) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
		log.Printf("alliance consumer: guild_left alliance_id=%d missing alliance_pb", evt.AllianceID)
		return nil
	}
	removedGuildPb, ok := allianceGuildFromMQPayload(raw)
	if !ok {
		log.Printf("alliance consumer: guild_left alliance_id=%d missing removed_guild_pb", evt.AllianceID)
		return nil
	}
	removedGuildID := removedGuildPb.GetGuildId()
	var extra struct {
		Expelled bool `json:"expelled"`
	}
	_ = json.Unmarshal(raw, &extra)
	gs.alliance.UpdateAsync(ctx, evt).Then(func(interface{}) (interface{}, error) {
		gs.alliance.BroadcastGuildLeft(alliancePb, removedGuildPb, extra.Expelled)
		return nil, nil
	})
	log.Printf("alliance consumer: applied guild_left alliance_id=%d removed_guild_id=%d revision=%d", evt.AllianceID, removedGuildID, evt.Revision)
	return nil
}
