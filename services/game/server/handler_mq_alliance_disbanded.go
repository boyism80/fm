package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/domain"
)

type allianceMqDisbanded struct{ gs *GameServer }

func (allianceMqDisbanded) New(gs *GameServer) *allianceMqDisbanded {
	return &allianceMqDisbanded{gs: gs}
}

func (*allianceMqDisbanded) EventType() string {
	return "disbanded"
}

func (h *allianceMqDisbanded) Handle(ctx actor.Context, raw json.RawMessage) error {
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
		guildIDs = gs.systems.Get[*domain.AllianceSystem]().GuildIDs(evt.AllianceID)
	}
	gs.systems.Get[*domain.AllianceSystem]().UpdateAsync(ctx, evt).Do(func() error {
		gs.systems.Get[*domain.AllianceSystem]().DisbandAsync(ctx, evt.AllianceID, guildIDs)
		return nil
	})
	return nil
}
