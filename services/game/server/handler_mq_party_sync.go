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

type partyMqPartySync struct{ gs *GameServer }

func (partyMqPartySync) New(gs *GameServer) *partyMqPartySync {
	return &partyMqPartySync{gs: gs}
}
func (*partyMqPartySync) EventType() string { return "party_sync" }
func (h *partyMqPartySync) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil || gs.party == nil {
		return nil
	}
	pc := gs.party
	var evt PartyEventEnvelope
	if err := json.Unmarshal(raw, &evt); err != nil {
		log.Printf("party consumer: party_sync envelope: %v", err)
		return nil
	}
	var payload struct {
		PartyPB string `json:"party_pb"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("party consumer: party_sync JSON: %v", err)
		return nil
	}
	if payload.PartyPB == "" {
		log.Printf("party consumer: party_sync: missing party_pb")
		return nil
	}
	wire, err := base64.StdEncoding.DecodeString(payload.PartyPB)
	if err != nil {
		log.Printf("party consumer: party_pb base64: %v", err)
		return nil
	}
	var partyPb internal.Party
	if err := proto.Unmarshal(wire, &partyPb); err != nil {
		log.Printf("party consumer: party_sync protobuf: %v", err)
		return nil
	}
	if partyPb.GetPartyId() != evt.PartyID {
		log.Printf("party_sync: party_id mismatch envelope=%d party=%d", evt.PartyID, partyPb.GetPartyId())
	}
	if _, err := pc.applyEmbeddedParty(evt, &partyPb, true); err != nil {
		log.Printf("party consumer: applyEmbeddedParty: %v", err)
	}
	return nil
}
