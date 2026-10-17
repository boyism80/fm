package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game"
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
		GuildIDs []uint32 `json:"guild_ids"`
		GuildPBs []string `json:"guild_pbs"`
	}
	if err := json.Unmarshal(raw, &extra); err != nil {
		log.Printf("alliance consumer: disbanded alliance_id=%d invalid payload: %v", evt.AllianceID, err)
		return nil
	}
	guildIDs := extra.GuildIDs
	if len(guildIDs) == 0 {
		guildIDs = gs.systems.Get[*game.AllianceSystem]().GuildIDs(evt.AllianceID)
	}
	guildPbs := make([]*internal.Guild, 0, len(extra.GuildPBs))
	for _, encoded := range extra.GuildPBs {
		if guildPb := decodeGuildSnapshot(encoded); guildPb != nil {
			guildPbs = append(guildPbs, guildPb)
		}
	}
	gs.systems.Get[*game.AllianceSystem]().Disband(evt.AllianceID, guildIDs, guildPbs)
	return nil
}
