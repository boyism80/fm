package server

import (
	"encoding/base64"
	"encoding/json"
	"log"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type partyMqPartySnapshot struct{ gs *GameServer }

func (partyMqPartySnapshot) New(gs *GameServer) *partyMqPartySnapshot {
	return &partyMqPartySnapshot{gs: gs}
}
func (*partyMqPartySnapshot) EventType() string { return "party_snapshot" }
func (h *partyMqPartySnapshot) Handle(_ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil || gs.party == nil {
		return nil
	}
	pc := gs.party
	var evt PartyEventEnvelope
	if err := json.Unmarshal(raw, &evt); err != nil {
		log.Printf("party consumer: party_snapshot envelope: %v", err)
		return nil
	}
	var payload struct {
		PartySnapshotPB string `json:"party_snapshot_pb"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("party consumer: party_snapshot JSON: %v", err)
		return nil
	}
	if payload.PartySnapshotPB == "" {
		log.Printf("party consumer: party_snapshot: missing party_snapshot_pb")
		return nil
	}
	wire, err := base64.StdEncoding.DecodeString(payload.PartySnapshotPB)
	if err != nil {
		log.Printf("party consumer: party_snapshot_pb base64: %v", err)
		return nil
	}
	var snap internal.PartySnapshot
	if err := proto.Unmarshal(wire, &snap); err != nil {
		log.Printf("party consumer: party_snapshot protobuf: %v", err)
		return nil
	}
	if snap.GetPartyId() != evt.PartyID {
		log.Printf("party_snapshot: party_id mismatch envelope=%d snapshot=%d", evt.PartyID, snap.GetPartyId())
	}
	if _, err := pc.applyEmbeddedPartySnapshot(evt, &snap, true); err != nil {
		log.Printf("party consumer: applyEmbeddedPartySnapshot: %v", err)
	}
	return nil
}
