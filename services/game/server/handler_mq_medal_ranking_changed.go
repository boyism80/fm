package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game"
	"google.golang.org/protobuf/encoding/protojson"
)

type medalMqRankingChanged struct{ gs *GameServer }

func (medalMqRankingChanged) New(gs *GameServer) *medalMqRankingChanged {
	return &medalMqRankingChanged{gs: gs}
}

func (*medalMqRankingChanged) EventType() string {
	return "medal_ranking_changed"
}

func (h *medalMqRankingChanged) Handle(_ actor.Context, raw json.RawMessage) error {
	evt := &internal.MedalRankingChanged{}
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(raw, evt); err != nil {
		log.Printf("medalMqRankingChanged: %v", err)
		return nil
	}
	medals := h.gs.systems.Get[*game.MedalSystem]()
	medals.Set(game.NewMedalRankingFromInternalProto(evt.GetRanking()))
	if evt.GetPreviousHolderId() != 0 {
		medals.Revoke(evt.GetPreviousHolderId(), evt.GetRanking().GetKind())
	}
	return nil
}
