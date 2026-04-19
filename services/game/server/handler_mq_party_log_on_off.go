package server

import (
	"encoding/base64"
	"encoding/json"
	"log"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type partyMqLogOnOff struct{ gs *GameServer }

func (partyMqLogOnOff) New(gs *GameServer) *partyMqLogOnOff { return &partyMqLogOnOff{gs: gs} }
func (*partyMqLogOnOff) EventType() string                  { return "log_onoff" }
func (h *partyMqLogOnOff) Handle(_ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil || gs.party == nil {
		return nil
	}
	pc := gs.party
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	if raw != nil {
		var payload struct {
			PartySnapshotPB string `json:"party_snapshot_pb"`
			CharacterID     uint32 `json:"character_id"`
		}
		if err := json.Unmarshal(raw, &payload); err == nil && payload.PartySnapshotPB != "" {
			wire, err := base64.StdEncoding.DecodeString(payload.PartySnapshotPB)
			if err == nil {
				var snap internal.PartySnapshot
				if err := proto.Unmarshal(wire, &snap); err == nil {
					applied, err := pc.applyEmbeddedPartySnapshot(evt, &snap, false)
					if err != nil {
						log.Printf("party consumer: log_onoff embedded snapshot: %v", err)
						return nil
					}
					if applied && payload.CharacterID != 0 {
						if snapshot := pc.CachedSnapshot(evt.PartyID); snapshot != nil {
							pc.DeliverPartyLogOnOff(snapshot, payload.CharacterID)
						}
					}
					return nil
				}
			}
		}
	}
	if err := pc.apply(evt); err != nil {
		log.Printf("party consumer: apply type=log_onoff party_id=%d: %v", evt.PartyID, err)
		return nil
	}
	if raw == nil {
		return nil
	}
	var extra struct {
		CharacterID uint32 `json:"character_id"`
	}
	if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 {
		return nil
	}
	snapshot := pc.CachedSnapshot(evt.PartyID)
	if snapshot != nil {
		pc.DeliverPartyLogOnOff(snapshot, extra.CharacterID)
	}
	return nil
}
