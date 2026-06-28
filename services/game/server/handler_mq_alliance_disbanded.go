package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type allianceMqDisbanded struct{ gs *GameServer }

func (allianceMqDisbanded) New(gs *GameServer) *allianceMqDisbanded {
	return &allianceMqDisbanded{gs: gs}
}

func (*allianceMqDisbanded) EventType() string {
	return "disbanded"
}

func (h *allianceMqDisbanded) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeAllianceEventEnvelope(raw)
	if !ok {
		return nil
	}
	var extra struct {
		GuildIDs           []uint32 `json:"guild_ids"`
		MemberCharacterIDs []uint32 `json:"member_character_ids"`
	}
	if err := json.Unmarshal(raw, &extra); err != nil {
		log.Printf("alliance consumer: disbanded alliance_id=%d invalid payload: %v", evt.AllianceID, err)
		return nil
	}
	guildIDs := extra.GuildIDs
	if len(guildIDs) == 0 {
		guildIDs = gs.alliance.GuildIDs(evt.AllianceID)
	}
	gs.alliance.UpdateAsync(ctx, evt).Then(func(interface{}) (interface{}, error) {
		gs.alliance.DisbandAfterGuildRefreshAsync(ctx, evt.AllianceID, guildIDs)
		return nil, nil
	})
	log.Printf("alliance consumer: applied disbanded alliance_id=%d revision=%d", evt.AllianceID, evt.Revision)
	return nil
}
