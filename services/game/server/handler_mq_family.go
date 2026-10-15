package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game"
	"google.golang.org/protobuf/encoding/protojson"
)

type familyMq struct{ gs *GameServer }

func (familyMq) New(gs *GameServer) *familyMq {
	return &familyMq{gs: gs}
}

func (*familyMq) EventType() string {
	return "family"
}

// Every member on this server shares the decoded family, so no one modifies it.
func (h *familyMq) Handle(_ actor.Context, raw json.RawMessage) error {
	var payload struct {
		CharacterIDs []uint32        `json:"character_ids"`
		Family       json.RawMessage `json:"family"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	family := &internal.Family{}
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(payload.Family, family); err != nil {
		log.Printf("familyMq: %v", err)
		return nil
	}

	for _, id := range payload.CharacterIDs {
		h.gs.systems.Get[*game.DispatchSystem]().SendToCharacter(nil, id, &game.SetFamily{
			CharacterID: id,
			Family:      family,
		})
	}
	return nil
}
