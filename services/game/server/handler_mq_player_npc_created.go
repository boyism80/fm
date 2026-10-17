package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game"
	"google.golang.org/protobuf/encoding/protojson"
)

type playerNpcMqCreated struct{ gs *GameServer }

func (playerNpcMqCreated) New(gs *GameServer) *playerNpcMqCreated {
	return &playerNpcMqCreated{gs: gs}
}

func (*playerNpcMqCreated) EventType() string {
	return "player_npc_created"
}

func (h *playerNpcMqCreated) Handle(_ actor.Context, raw json.RawMessage) error {
	npc := &internal.PlayerNpc{}
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(raw, npc); err != nil {
		log.Printf("playerNpcMqCreated: %v", err)
		return nil
	}
	h.gs.systems.Get[*game.PlayerNpcSystem]().Place(game.NewPlayerNpcFromInternalProto(npc))
	return nil
}
