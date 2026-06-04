package server

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type allianceMqCreated struct{ gs *GameServer }

func (allianceMqCreated) New(gs *GameServer) *allianceMqCreated {
	return &allianceMqCreated{gs: gs}
}

func (*allianceMqCreated) EventType() string {
	return "created"
}

func allianceFromMQPayload(raw json.RawMessage) (*internal.Alliance, bool) {
	if raw == nil {
		return nil, false
	}
	var payload struct {
		AlliancePB string `json:"alliance_pb"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil || payload.AlliancePB == "" {
		return nil, false
	}
	wire, err := base64.StdEncoding.DecodeString(payload.AlliancePB)
	if err != nil {
		return nil, false
	}
	var alliancePb internal.Alliance
	if err := proto.Unmarshal(wire, &alliancePb); err != nil {
		return nil, false
	}
	return &alliancePb, true
}

func (h *allianceMqCreated) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
		log.Printf("alliance consumer: created alliance_id=%d missing alliance_pb", evt.AllianceID)
		return nil
	}
	gs.alliance.Update(alliancePb)
	gs.alliance.BroadcastCreate(alliancePb)
	log.Printf("alliance consumer: applied created alliance_id=%d revision=%d", evt.AllianceID, evt.Revision)
	return nil
}
